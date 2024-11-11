package main

import (
	"os"
)

// main - тестовая
func main() {
	os.Exit(0) //want "os.Exit is not allowed in main package"
}
