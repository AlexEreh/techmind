package document

import (
	"strings"

	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UploadHandler struct {
	documentService service.DocumentService
}

func NewUploadHandler(documentService service.DocumentService) *UploadHandler {
	return &UploadHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Загрузка документа
// @Description  Загружает новый документ в систему с файлом
// @Tags         documents
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        company_id formData string true "ID компании" format:"uuid"
// @Param        name formData string true "Название документа"
// @Param        folder_id formData string false "ID папки" format:"uuid"
// @Param        sender_id formData string false "ID отправителя" format:"uuid"
// @Param        file formData file true "Файл документа"
// @Success      201 {object} DocumentResponse "Документ успешно загружен"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents [post]
func (h *UploadHandler) Handle(c fiber.Ctx) error {
	// Получаем user_id из контекста (установлено JWT middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(handlers.ErrorResponse{
			Error: "unauthorized",
		})
	}

	var req UploadRequest
	if err := c.Bind().Form(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "file is required",
		})
	}

	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: "failed to open file",
		})
	}
	defer fileReader.Close()

	input := service.DocumentUploadInput{
		CompanyID: req.CompanyID,
		FolderID:  req.FolderID,
		Name:      req.Name,
		File:      fileReader,
		FileSize:  file.Size,
		MimeType:  file.Header.Get("Content-Type"),
		SenderID:  req.SenderID,
		UserID:    userID,
	}

	// ...existing code...

	document, err := h.documentService.Upload(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "not supported") {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(handlers.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(DocumentResponse{
		ID:              document.ID,
		CompanyID:       document.CompanyID,
		FolderID:        document.FolderID,
		SenderID:        document.SenderID,
		Name:            document.Name,
		FilePath:        document.FilePath,
		PreviewFilePath: document.PreviewFilePath,
		FileSize:        document.FileSize,
		MimeType:        document.MimeType,
		Checksum:        document.Checksum,
		CreatedBy:       document.CreatedBy,
		UpdatedBy:       document.UpdatedBy,
		CreatedAt:       document.CreatedAt,
		UpdatedAt:       document.UpdatedAt,
	})
}
