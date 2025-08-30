package main

import (
	"embed"
	"log"

	"asynclab.club/asynx/backend/cmd"
	"github.com/joho/godotenv"
)

//go:embed frontend/dist/* templates/*
var embedFS embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cmd.Main(embedFS)
}
