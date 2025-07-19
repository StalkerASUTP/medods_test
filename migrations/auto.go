package main

import (
	"authentication-app/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func main() {
	conf := config.NewConfig()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DBName) //postgres://<username>:<password>@<host>:<port>/<dbname>

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		panic("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL PRIMARY KEY,             
    refresh_token_hash VARCHAR(255) NOT NULL,        
    refresh_token_expires_at TIMESTAMP NOT NULL,     
    user_agent VARCHAR(500) NOT NULL,                
    ip_address VARCHAR(45) NOT NULL,                 
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW())`)
	if err != nil {
		panic("failed to create table: " + err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		panic("failed to commit transaction: " + err.Error())
	}
	fmt.Println("table created successfully")
}
