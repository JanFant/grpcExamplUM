package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

var (
	tableUsers = `
	CREATE TABLE users(
	    id SERIAL PRIMARY KEY,
	    name TEXT,
	    age INT
	);`
)

func ConnectPDB() (*pgx.Conn, error) {
	var ctx = context.Background()
	conn, err := pgx.Connect(ctx, GlobalConfig.PSQLConfig.GetPSQLUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, `SELECT * FROM "testDB".public.users limit(1);`)
	if err != nil {
		fmt.Println("Create Users table!")
		_, _ = conn.Exec(ctx, tableUsers)
	}
	return conn, nil
}
