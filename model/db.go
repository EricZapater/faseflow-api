package model

import (
	"database/sql"
	"faseflow-api/environment"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Setup() {
	env := environment.LoadEnvironment()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		env.DbHost, env.DbPort, env.DbUser, env.DbPass, env.DbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}