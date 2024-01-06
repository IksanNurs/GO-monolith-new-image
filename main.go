package main

import (
	"akuntansi/route"
	"akuntansi/utils"
	_ "embed"
	"os"
)

//go:embed .env
var env string

func main() {
	utils.LoadEnv(env)
	// database.StartDB()
	router := route.GetGinRoute()
	router.Run(":" + os.Getenv("PORT"))
}
