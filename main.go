package main

import (
	"os"
	"path"

	"dotins.eu.org/financeIB/src/bundlers"
	"dotins.eu.org/financeIB/src/extractors"
	"github.com/samber/lo"
)

func main() {
	var inputFileName = os.Args[1]
	var outputFileName, _ = lo.Nth(os.Args, 2)
	var extension = path.Ext(inputFileName)

	if extension == ".ofx" {
		bundlers.XLSXInit(extractors.OFXInit(inputFileName), outputFileName)
		return
	}

	if extension == ".xlsx" || extension == ".xls" {
		bundlers.PDFInit(extractors.XLSXInit(inputFileName, outputFileName))
		return
	}
}
