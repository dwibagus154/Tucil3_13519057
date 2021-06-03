package main

import (
	"fmt"

	"example.com/handler"
)

// file ini berfungsi untuk mengintegrasikan handler.go dengan graph.go
func main() {
	fmt.Println("Hello, World!")
	fmt.Println("Selamat datang di aplikasi MAP IKI?")
	handler.Start()
}
