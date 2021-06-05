package main

import (
	"flag"
	"github.com/eyebrow-fish/gogo"
	"log"
	"os"
)

func main() {
	flag.Parse()
	file := flag.Arg(0)
	if file == "" {
		log.Fatal("expected file for first argument")
	}

	fileText, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	tokens := gogo.Parse(string(fileText))
	tree, err := gogo.BuildTrees(tokens)
	if err != nil {
		log.Fatal(err)
	}

	if err = gogo.Go(*tree); err != nil {
		log.Fatal(err)
	}
}
