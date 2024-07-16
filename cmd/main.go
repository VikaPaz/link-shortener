package main

import (
	"fmt"
	"links-shorter/internal/app"
)

func main() {
	err := app.Run()

	fmt.Println(err)
}
