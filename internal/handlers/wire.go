package handlers

import (
	"github.com/WildEgor/e-shop-auth/internal/adapters"
	change_email_handler "github.com/WildEgor/e-shop-auth/internal/handlers/change-email"
	change_password_handler "github.com/WildEgor/e-shop-auth/internal/handlers/change-password"
	handlers4 "github.com/WildEgor/e-shop-auth/internal/handlers/change-phone"
	handlers3 "github.com/WildEgor/e-shop-auth/internal/handlers/confirm-email"
	handlers5 "github.com/WildEgor/e-shop-auth/internal/handlers/confirm-phone"
	eh "github.com/WildEgor/e-shop-auth/internal/handlers/errors"
	hch "github.com/WildEgor/e-shop-auth/internal/handlers/health_check"
	handlers "github.com/WildEgor/e-shop-auth/internal/handlers/login"
	logout_handler "github.com/WildEgor/e-shop-auth/internal/handlers/logout"
	me_handler "github.com/WildEgor/e-shop-auth/internal/handlers/me"
	otp_generate_handler "github.com/WildEgor/e-shop-auth/internal/handlers/otp-generate"
	otp_login_handler "github.com/WildEgor/e-shop-auth/internal/handlers/otp-login"
	handlers2 "github.com/WildEgor/e-shop-auth/internal/handlers/ready_check"
	refresh_handler "github.com/WildEgor/e-shop-auth/internal/handlers/refresh"
	reg_handler "github.com/WildEgor/e-shop-auth/internal/handlers/reg"
	update_handler "github.com/WildEgor/e-shop-auth/internal/handlers/update-profile"
	"github.com/WildEgor/e-shop-auth/internal/proto"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	adapters.AdaptersSet,
	services.ServicesSet,
	proto.NewAuthService,
	proto.NewGRPCServer,
	hch.NewHealthCheckHandler,
	eh.NewErrorsHandler,
	me_handler.NewMeHandler,
	refresh_handler.NewRefreshHandler,
	reg_handler.NewRegHandler,
	logout_handler.NewLogoutHandler,
	handlers.NewLoginHandler,
	change_password_handler.NewChangePasswordHandler,
	otp_generate_handler.NewOTPGenHandler,
	otp_login_handler.NewOTPLoginHandler,
	handlers3.NewConfirmEmailHandler,
	handlers4.NewChangePhoneHandler,
	handlers5.NewConfirmPhoneHandler,
	change_email_handler.NewChangeEmailHandler,
	update_handler.NewUpdateProfileHandler,
	handlers2.NewReadyCheckHandler,
)
