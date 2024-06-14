package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/kociumba/gabagool-file/gabagool"
)

var g = new(gabagool.GabagoolFile)

func main() {

	// testing the package
	// f := new(g.GabagoolFile)
	// err := f.CreateAndSave("./test", g.Text, []byte("test"))
	// if err != nil {
	// 	panic(err)
	// }

	if os.Args[1] == "create" {
		gabagool, err := g.CreateFile(gabagool.Text, []byte(os.Args[2]))
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		gabagool.SaveWithBitPacking(os.Args[3], gabagool)
		os.Exit(0)
	}

	// Check if a file path was provided as a command line argument
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		if filepath.Ext(filePath) != ".gabagool" {
			log.Error("not a .gabagool file")
			os.Exit(0)
		}
		filePath = filepath.Clean(filePath)
		filePath, _ = strings.CutSuffix(filePath, ".gabagool")
		handleDataTypes(filePath)
	} else {
		log.Error("no file path provided")
		os.Exit(0)
	}

}

func handleDataTypes(filePath string) {
	// Open the file using the Open function
	f, err := g.ParseFile(filePath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info(f)

	switch f.DataType {
	case gabagool.Text:
		displayData(f.Data, f.DataType)
	case gabagool.Image:
		displayData(f.Data, f.DataType)
	case gabagool.Bytes:
		displayData(f.Data, f.DataType)
	}

}
