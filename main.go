package main

import (
	"embed"

	"asynclab.club/asynx/backend/cmd"
	"github.com/joho/godotenv"
)

//go:embed frontend/dist/* templates/*
var embedFS embed.FS

func main() {
	_ = godotenv.Load()
	cmd.Main(embedFS)
}
