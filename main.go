package main

import (
	"os"
	"strings"

	"dotins.eu.org/financeIB/src/bundlers"
	"dotins.eu.org/financeIB/src/extractors"
)

func main() {
	var fileName = strings.Split(os.Args[1], ".")
	var extension = fileName[len(fileName)-1]

	if extension == "ofx" {
		bundlers.XLSXInit(extractors.OFXInit(os.Args[1]))
		return
	}

	if extension == "xlsx" || extension == "xls" {
		bundlers.PDFInit(extractors.XLSXInit(os.Args[1]))
		return
	}
}
