package http

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers/auth"
	"techmind/internal/transport/http/handlers/document"
	"techmind/internal/transport/http/handlers/documenttag"
	"techmind/internal/transport/http/handlers/folder"
	"techmind/pkg/config"

	"github.com/gofiber/fiber/v3"
)

type ServerDeps struct {
	AuthService        service.AuthService
	FolderService      service.FolderService
	DocumentService    service.DocumentService
	DocumentTagService service.DocumentTagService
	Config             *config.Config
}

type Server struct {
	app  *fiber.App
	deps ServerDeps
}

func NewServer(deps ServerDeps) *Server {
	server := &Server{
		app:  fiber.New(),
		deps: deps,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Публичные маршруты
	public := s.app.Group("/api/v1/public")

	// Регистрация маршрутов аутентификации
	authGroup := public.Group("/auth")
	auth.RegisterRoutes(authGroup, s.deps.AuthService)

	// Приватные маршруты
	private := s.app.Group("/api/v1/private")
	private.Use(s.jwtMiddleware)

	// Регистрация маршрутов для папок
	foldersGroup := private.Group("/folders")
	folder.RegisterRoutes(foldersGroup, s.deps.FolderService)

	// Регистрация маршрутов для документов
	documentsGroup := private.Group("/documents")
	document.RegisterRoutes(documentsGroup, s.deps.DocumentService)

	// Регистрация маршрутов для тегов документов
	documentTagsGroup := private.Group("/document-tags")
	documenttag.RegisterRoutes(documentTagsGroup, s.deps.DocumentTagService)
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}
