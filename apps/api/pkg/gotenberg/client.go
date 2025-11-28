package gotenberg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Client представляет клиент для работы с Gotenberg API
type Client struct {
	baseURL    string
	httpClient *http.Client
	trace      string
}

// ClientOption определяет опцию для конфигурации клиента
type ClientOption func(*Client)

// NewClient создает новый экземпляр клиента Gotenberg
func NewClient(baseURL string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithHTTPClient устанавливает пользовательский HTTP клиент
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTrace устанавливает trace ID для всех запросов
func WithTrace(trace string) ClientOption {
	return func(c *Client) {
		c.trace = trace
	}
}

// Response представляет ответ от Gotenberg API
type Response struct {
	Body        []byte
	ContentType string
	Filename    string
	Trace       string
}

// doRequest выполняет HTTP запрос к Gotenberg API
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, contentType string, headers map[string]string) (*Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	if c.trace != "" {
		req.Header.Set("Gotenberg-Trace", c.trace)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	response := &Response{
		Body:        responseBody,
		ContentType: resp.Header.Get("Content-Type"),
		Trace:       resp.Header.Get("Gotenberg-Trace"),
	}

	// Извлечение имени файла из Content-Disposition
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		response.Filename = parseFilename(contentDisposition)
	}

	return response, nil
}

// parseFilename извлекает имя файла из заголовка Content-Disposition
func parseFilename(contentDisposition string) string {
	// Простой парсинг: "attachment; filename=file.pdf"
	const prefix = "filename="
	start := bytes.Index([]byte(contentDisposition), []byte(prefix))
	if start == -1 {
		return ""
	}
	start += len(prefix)
	filename := contentDisposition[start:]
	return filename
}

// File представляет файл для загрузки
type File struct {
	Name    string
	Content []byte
}

// DownloadFrom представляет URL для скачивания файла
type DownloadFrom struct {
	URL              string            `json:"url"`
	ExtraHTTPHeaders map[string]string `json:"extraHttpHeaders,omitempty"`
	Embedded         bool              `json:"embedded,omitempty"`
}

// writeMultipartForm создает multipart/form-data запрос
func writeMultipartForm(writer *multipart.Writer, files []File, fields map[string]string) error {
	// Добавляем файлы
	for _, file := range files {
		part, err := writer.CreateFormFile("files", file.Name)
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := part.Write(file.Content); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}
	}

	// Добавляем поля формы
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	return nil
}

// GetHealth проверяет статус здоровья Gotenberg сервера
func (c *Client) GetHealth(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/health", nil, "", nil)
	if err != nil {
		return nil, err
	}

	var health HealthResponse
	if err := json.Unmarshal(resp.Body, &health); err != nil {
		return nil, fmt.Errorf("failed to unmarshal health response: %w", err)
	}

	return &health, nil
}

// HealthResponse представляет ответ от /health endpoint
type HealthResponse struct {
	Status  string                 `json:"status"`
	Details map[string]HealthCheck `json:"details"`
}

// HealthCheck представляет информацию о проверке здоровья модуля
type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// GetVersion возвращает версию Gotenberg
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/version", nil, "", nil)
	if err != nil {
		return "", err
	}

	return string(resp.Body), nil
}

// marshalJSONField маршалирует значение в JSON строку для использования в форме
func marshalJSONField(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// addFilesToWriter добавляет файлы в multipart writer
func addFilesToWriter(writer *multipart.Writer, files []File, fieldName string) error {
	for _, file := range files {
		part, err := writer.CreateFormFile(fieldName, file.Name)
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(file.Content); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}
	}
	return nil
}

// addFieldsToWriter добавляет поля формы в multipart writer
func addFieldsToWriter(writer *multipart.Writer, fields map[string]string) error {
	for key, value := range fields {
		if value != "" {
			if err := writer.WriteField(key, value); err != nil {
				return fmt.Errorf("failed to write field %s: %w", key, err)
			}
		}
	}
	return nil
}
