package main

import "fmt"

var ()

func dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		EnvOrDefault("POSTGRES_USER", "postgres"),
		EnvOrDefault("POSTGRES_PASSWORD", "postgres"),
		EnvOrDefault("POSTGRES_ADDRESS", "localhost:5432"),
		EnvOrDefault("POSTGRES_DB", "postgres"),
	)
}
