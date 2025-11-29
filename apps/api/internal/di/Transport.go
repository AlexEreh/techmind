package di

import (
	"context"
	"fmt"

	"techmind/internal/service"
	"techmind/internal/transport/http"
	"techmind/pkg/config"

	"go.uber.org/fx"
)

var Transport = fx.Options(
	fx.Provide(
		func(
			authService service.AuthService,
			folderService service.FolderService,
			documentService service.DocumentService,
			documentTagService service.DocumentTagService,
			senderService service.SenderService,
			companyUserService service.CompanyUserService,
			companyService service.CompanyService,
			cfg *config.Config,
		) *http.Server {
			deps := http.ServerDeps{
				AuthService:        authService,
				FolderService:      folderService,
				DocumentService:    documentService,
				DocumentTagService: documentTagService,
				SenderService:      senderService,
				CompanyUserService: companyUserService,
				CompanyService:     companyService,
				Config:             cfg,
			}
			return http.NewServer(deps)
		},
	),
	fx.Invoke(startHTTPServer),
)

func startHTTPServer(server *http.Server, cfg *config.Config, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%d", cfg.HTTPPort)
				if err := server.Listen(addr); err != nil {
					fmt.Printf("HTTP server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown()
		},
	})
}
