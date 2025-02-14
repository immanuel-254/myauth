package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/immanuel-254/myauth/internal"
)

var (
	Login = internal.View{
		Route:   "/login",
		Handler: http.HandlerFunc(internal.Login),
	}

	Logout = internal.View{
		Route:   "/logout",
		Handler: http.HandlerFunc(internal.Logout),
	}

	Signup = internal.View{
		Route:   "/signup",
		Handler: http.HandlerFunc(internal.Signup),
	}

	ActivateEmail = internal.View{
		Route:   "/activate",
		Handler: http.HandlerFunc(internal.ActivateEmail),
	}

	UserRead = internal.View{
		Route:       "/read",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.UserRead),
	}

	UserList = internal.View{
		Route:       "/list",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.UserList),
	}

	ChangeEmailRequest = internal.View{
		Route:       "/change-email-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangeEmailRequest),
	}

	ChangeEmail = internal.View{
		Route:       "/change-email",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangeEmail),
	}

	ChangePasswordRequest = internal.View{
		Route:       "/change-password-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangePasswordRequest),
	}

	ChangePassword = internal.View{
		Route:       "/change-password",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangePassword),
	}

	ResetPasswordRequest = internal.View{
		Route:       "/reset-password-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ResetPasswordRequest),
	}

	ResetPassword = internal.View{
		Route:       "/reset-password",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ResetPassword),
	}

	DeleteUserRequest = internal.View{
		Route:       "/delete-user-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.DeleteUserRequest),
	}

	DeleteUser = internal.View{
		Route:       "/delete-user",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.DeleteUser),
	}

	IsActiveChange = internal.View{
		Route:       "/isactive",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.IsActiveChange),
	}

	IsStaffChange = internal.View{
		Route:       "/isstaff",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.IsStaffChange),
	}

	SessionList = internal.View{
		Route:       "/session/list",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.SessionList),
	}

	LogList = internal.View{
		Route:       "/log/list",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.LogList),
	}

	DashBoard = internal.View{
		Route:       "/",
		Middlewares: []func(http.Handler) http.Handler{internal.DashRequireAdmin},
		Handler:     http.HandlerFunc(internal.Dashboard),
	}

	DashLogin = internal.View{
		Route:   "/dash-login",
		Handler: http.HandlerFunc(internal.Dashlogin),
	}

	DashLogout = internal.View{
		Route:   "/dash-logout",
		Handler: http.HandlerFunc(internal.Dashlogout),
	}

	style = internal.View{
		Route:   "/static/styles.css",
		Handler: http.HandlerFunc(internal.StyleCss),
	}
	script = internal.View{
		Route:   "/static/script.js",
		Handler: http.HandlerFunc(internal.ScriptJs),
	}
)

func Api() {
	mux := http.NewServeMux()

	allviews := []internal.View{
		Login,
		Logout,
		Signup,
		ActivateEmail,
		UserRead,
		UserList,
		ChangeEmailRequest,
		ChangeEmail,
		ResetPasswordRequest,
		ResetPassword,
		DeleteUserRequest,
		DeleteUser,
		IsActiveChange,
		IsStaffChange,

		SessionList,

		LogList,

		DashBoard,
		DashLogin,
		DashLogout,

		style,
		script,
	}

	internal.Routes(mux, allviews)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("PORT")), // Custom port
		//Handler:      internal.LoggingMiddleware(internal.Cors(internal.New(internal.ConfigDefault)(mux))), // Attach the mux as the handler
		Handler:      internal.LoggingMiddleware(mux),
		ReadTimeout:  10 * time.Second, // Set read timeout
		WriteTimeout: 10 * time.Second, // Set write timeout
		IdleTimeout:  30 * time.Second, // Set idle timeout
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
