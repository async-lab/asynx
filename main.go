package main

import (
	"embed"

	"asynclab.club/asynx/backend/cmd"
	"github.com/joho/godotenv"
)

//go:embed all:frontend/dist/* all:templates/*
var embedFS embed.FS

// @title Asynx API 文档
// @version 1.0
// @description Asynx API 接口文档

// @license.name AGPL-v3.0
// @license.url https://raw.githubusercontent.com/async-lab/asynx/refs/heads/main/LICENSE

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 Bearer Token，格式为 "Bearer <token>"

// @host
// @BasePath /api
func main() {
	_ = godotenv.Load()
	cmd.Main(embedFS)
}
