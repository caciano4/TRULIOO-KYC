package main

import (
	"trullio-kyc/config"
	"trullio-kyc/routes"
)

func main() {
	config.LoadEnv()
	routes.InitRoutes()
}
