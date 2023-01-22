package main

import (
	"os"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/database"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/router"
)
func main() {
	db := database.SetupDB()

	r := router.InitRoutes(db)
	r.Run(":" + os.Getenv("PORT"))
}