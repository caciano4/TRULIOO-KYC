package main

import (
	"fmt"
	"trullio-kyc/config"
	"trullio-kyc/routes"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	fmt.Print("teste")
	routes.InitRoutes(db)
}
