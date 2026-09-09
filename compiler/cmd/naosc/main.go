package main

import (
	"fmt"

	tsitter "github.com/tree-sitter/go-tree-sitter"
)

func main() {
	parser := tsitter.NewParser()
	defer parser.Close()

	fmt.Println("Tree-sitter OK")
}
