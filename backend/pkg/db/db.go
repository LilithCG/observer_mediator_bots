package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func Init() (*sql.DB, error) {
	con, err := getEnv()
	if err != nil {
		return nil, err
	}

	db, err := connect(con)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv() (string, error) {
	err := godotenv.Load("C:\\Users\\saity\\PycharmProjects\\observer_bot\\backend\\configs\\variables.env")
	if err != nil {
		return "", err
	}

	connectionString := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=telegram sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("USER"),
		os.Getenv("PASS"),
	)

	return connectionString, err
}

func connect(con string) (*sql.DB, error) {
	db, err := sql.Open("postgres", con)
	if err != nil {
		return nil, err
	}
	return db, nil
}
