package main

import (
	"fmt"
	route "task-5-pbi-fullstack-developer-rulli-damara-putra/router"
)

func main() {

	route := route.SetupRouter()

	port := 8080
	address := fmt.Sprintf("localhost:%d", port)
	fmt.Println("Server running on ", address)

	route.Run(address)
}
