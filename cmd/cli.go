package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"syscall"

	"github.com/immanuel-254/myauth/internal"
	"github.com/immanuel-254/myauth/internal/models"
	"golang.org/x/term"
)

func CreateAdminUser() {
	var email string
	var password string

	fmt.Println("Email: ")
	fmt.Scanln(&email)

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	password = string(bytePassword)

	hash, err := internal.HashPassword(password)
	if err != nil {
		panic(err)
	}

	queries := models.New(internal.DB)
	ctx := context.Background()

	_, err = queries.UserCreate(ctx, models.UserCreateParams{
		Email:    email,
		Password: hash,
		Isactive: sql.NullBool{Bool: true, Valid: true},
		Isstaff:  sql.NullBool{Bool: true, Valid: true},
		Isadmin:  sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		panic(err)
	}
}
