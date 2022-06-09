package main

import (
	"github.com/obgnail/golangInPython/converter"
	"log"
)

func main() {
	source := "./demo2/source"
	target := "./demo2/target"
	goDir := "/Users/heyingliang/go/go1.16/bin/go"

	err := converter.Convert(source, target, goDir)
	if err != nil {
		log.Fatal(err)
	}
}
