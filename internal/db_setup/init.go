package db_setup

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"

	"full-auth/database"
)

//go:embed schema.sql
var ddl string

var ctx = context.Background()

func DbInit() (*database.Queries, error) {
	db, err := sql.Open("sqlite3", "accounts.db")
	if err != nil {
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := database.New(db)

	return queries, nil
}

func GetUser(username string) (string, string) {
	queries, err := DbInit()
	if err != nil {
		panic(err)
	}
	accountName, err := queries.GetAccountUsername(ctx, username)
	if err != nil {
		return "", ""
	}

	accountPassword, err := queries.GetAccountPassword(ctx, username)
	if err != nil {
		return "", ""
	}

	return accountName, accountPassword
}

func CreateAccount(username string, hashedPassword string, email string) string {
	queries, err := DbInit()
	if err != nil {
		panic(err)
	}

	accountName, err := queries.CraeteAccount(ctx, database.CraeteAccountParams{Username: username, Password: hashedPassword, Email: email})
	if err != nil {
		panic(err)
	}

	return accountName
}

func DeleteAccount(username string) {
	queries, err := DbInit()
	if err != nil {
		panic(err)
	}

	queries.DeleteAccount(ctx, username)
}
