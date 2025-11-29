package di

import (
	"techmind/internal/service/auth"
	"techmind/internal/service/company"
	"techmind/internal/service/company_user"
	"techmind/internal/service/document"
	"techmind/internal/service/documenttag"
	"techmind/internal/service/folder"
	"techmind/internal/service/sender"

	"go.uber.org/fx"
)

var Service = fx.Options(
	fx.Provide(
		auth.NewService,
		document.NewService,
		documenttag.NewService,
		folder.NewService,
		company_user.NewService,
		company.NewService,
		sender.NewService,
	),
)
