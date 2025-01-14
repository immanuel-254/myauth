package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/immanuel-254/myauth/internal"
)

var (
	Login = internal.View{
		Route:   "/login",
		Handler: http.HandlerFunc(internal.Login),
		Methods: []string{"POST"},
	}

	Logout = internal.View{
		Route:   "/logout",
		Handler: http.HandlerFunc(internal.Logout),
		Methods: []string{"POST"},
	}

	Signup = internal.View{
		Route:   "/signup",
		Handler: http.HandlerFunc(internal.Signup),
		Methods: []string{"POST"},
	}

	ActivateEmail = internal.View{
		Route:   "/activate",
		Handler: http.HandlerFunc(internal.ActivateEmail),
		Methods: []string{"PUT"},
	}

	ChangeEmailRequest = internal.View{
		Route:       "/change-email-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangeEmailRequest),
		Methods:     []string{"POST"},
	}

	ChangeEmail = internal.View{
		Route:       "/change-email",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangeEmail),
		Methods:     []string{"PUT"},
	}

	ChangePasswordRequest = internal.View{
		Route:       "/change-password-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangePasswordRequest),
		Methods:     []string{"POST"},
	}

	ChangePassword = internal.View{
		Route:       "/change-password",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ChangePassword),
		Methods:     []string{"PUT"},
	}

	ResetPasswordRequest = internal.View{
		Route:       "/reset-password-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ResetPasswordRequest),
		Methods:     []string{"POST"},
	}

	ResetPassword = internal.View{
		Route:       "/reset-password",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.ResetPassword),
		Methods:     []string{"PUT"},
	}

	DeleteUserRequest = internal.View{
		Route:       "/delete-user-request",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.DeleteUserRequest),
		Methods:     []string{"POST"},
	}

	DeleteUser = internal.View{
		Route:       "/delete-user",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAuth},
		Handler:     http.HandlerFunc(internal.DeleteUser),
		Methods:     []string{"DELETE"},
	}

	IsActiveChange = internal.View{
		Route:       "/isactive",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.IsActiveChange),
		Methods:     []string{"PUT"},
	}

	IsStaffChange = internal.View{
		Route:       "/isstaff",
		Middlewares: []func(http.Handler) http.Handler{internal.RequireAdmin},
		Handler:     http.HandlerFunc(internal.IsStaffChange),
		Methods:     []string{"PUT"},
	}
)

func Api() {
	mux := http.NewServeMux()

	allviews := []internal.View{
		Login,
		Logout,
		Signup,
		ActivateEmail,
		ChangeEmailRequest,
		ChangeEmail,
		ResetPasswordRequest,
		ResetPassword,
		DeleteUserRequest,
		DeleteUser,
		IsActiveChange,
		IsStaffChange,
	}

	internal.Routes(mux, allviews)

	server := &http.Server{
		Addr:         ":8080",                                                  // Custom port
		Handler:      internal.Cors(internal.New(internal.ConfigDefault)(mux)), // Attach the mux as the handler
		ReadTimeout:  10 * time.Second,                                         // Set read timeout
		WriteTimeout: 10 * time.Second,                                         // Set write timeout
		IdleTimeout:  30 * time.Second,                                         // Set idle timeout
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
