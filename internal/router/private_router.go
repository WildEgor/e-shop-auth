package router

import (
	change_email_handler "github.com/WildEgor/e-shop-auth/internal/handlers/change-email"
	change_password_handler "github.com/WildEgor/e-shop-auth/internal/handlers/change-password"
	handlers "github.com/WildEgor/e-shop-auth/internal/handlers/change-phone"
	handlers3 "github.com/WildEgor/e-shop-auth/internal/handlers/confirm-email"
	handlers2 "github.com/WildEgor/e-shop-auth/internal/handlers/confirm-phone"
	me_handler "github.com/WildEgor/e-shop-auth/internal/handlers/me"
	handlers4 "github.com/WildEgor/e-shop-auth/internal/handlers/refresh"
	handlers5 "github.com/WildEgor/e-shop-auth/internal/handlers/update-profile"
	auth_middleware "github.com/WildEgor/e-shop-auth/internal/middlewares/auth"
	"github.com/WildEgor/e-shop-auth/internal/repositories"
	"github.com/WildEgor/e-shop-auth/internal/services"
	"github.com/gofiber/fiber/v3"
)

type PrivateRouter struct {
	cp  *change_password_handler.ChangePasswordHandler
	ce  *change_email_handler.ChangeEmailHandler
	cpo *handlers.ChangePhoneHandler
	cop *handlers2.ConfirmPhoneHandler
	coe *handlers3.ConfirmEmailHandler
	me  *me_handler.MeHandler
	rt  *handlers4.RefreshHandler
	uph *handlers5.UpdateProfileHandler
	ur  *repositories.UserRepository
	jwt *services.JWTAuthenticator
}

func NewPrivateRouter(
	cp *change_password_handler.ChangePasswordHandler,
	ce *change_email_handler.ChangeEmailHandler,
	cpo *handlers.ChangePhoneHandler,
	cop *handlers2.ConfirmPhoneHandler,
	coe *handlers3.ConfirmEmailHandler,
	me *me_handler.MeHandler,
	rt *handlers4.RefreshHandler,
	uph *handlers5.UpdateProfileHandler,
	ur *repositories.UserRepository,
	jwt *services.JWTAuthenticator,
) *PrivateRouter {
	return &PrivateRouter{
		cp,
		ce,
		cpo,
		cop,
		coe,
		me,
		rt,
		uph,
		ur,
		jwt,
	}
}

func (r *PrivateRouter) Setup(app *fiber.App) {
	v1 := app.Group("/api/v1")
	ac := v1.Group("/auth")
	uc := v1.Group("/user")

	am := auth_middleware.NewAuthMiddleware(auth_middleware.AuthMiddlewareConfig{
		UR:  r.ur,
		JWT: r.jwt,
	})

	ac.Post("change-password", am, r.cp.Handle)
	ac.Post("change-phone", am, r.cpo.Handle)
	ac.Post("change-email", am, r.ce.Handle)
	ac.Post("confirm-phone", am, r.cop.Handle)
	ac.Post("confirm-email", am, r.coe.Handle)
	ac.Post("refresh-token", am, r.rt.Handle)
	uc.Put("profile", am, r.uph.Handle)
	uc.Get("me", am, r.me.Handle)
}
