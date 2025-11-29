package gotenberg

import (
	"context"
	"net/http"
	"techmind/pkg/config"
	"techmind/pkg/gotenberg"
	"time"

	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, cfg *config.Config) *gotenberg.Client {
	// Если Gotenberg отключен, возвращаем nil
	//if !cfg.Gotenberg.Enabled {
	//	return nil
	//}

	// Устанавливаем таймаут (по умолчанию 60 секунд)
	timeout := 60 * time.Second
	if cfg.Gotenberg.Timeout > 0 {
		timeout = time.Duration(cfg.Gotenberg.Timeout) * time.Second
	}

	// Создаем HTTP клиент с настроенным таймаутом
	httpClient := &http.Client{
		Timeout: timeout,
	}

	// Создаем Gotenberg клиент
	client := gotenberg.NewClient(
		cfg.Gotenberg.URL,
		gotenberg.WithHTTPClient(httpClient),
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Gotenberg клиент не требует явного закрытия
			return nil
		},
	})

	return client
}
