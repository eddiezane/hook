package main

import (
	"log"

	"github.com/eddiezane/hook/cmd"
)

func main() {
	if err := cmd.GenMarkdownTree("../../docs"); err != nil {
		log.Fatal(err)
	}
}
