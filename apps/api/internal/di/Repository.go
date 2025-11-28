package di

import (
	"techmind/internal/repo/company"
	"techmind/internal/repo/company_user"
	"techmind/internal/repo/document"
	"techmind/internal/repo/document_tag"
	"techmind/internal/repo/folder"
	"techmind/internal/repo/sender"
	"techmind/internal/repo/tag"
	"techmind/internal/repo/user"

	"go.uber.org/fx"
)

var Repository = fx.Options(
	fx.Provide(
		user.NewRepository,
		company.NewRepository,
		company_user.NewRepository,
		folder.NewRepository,
		sender.NewRepository,
		document.NewRepository,
		tag.NewRepository,
		document_tag.NewRepository,
	),
)
