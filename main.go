package main

import (
	"faseflow-api/model"
	"faseflow-api/routes"
)

func main() {
	model.Setup()
	routes.Setup()
}