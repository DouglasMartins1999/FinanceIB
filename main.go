package main

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"dotins.eu.org/financeIB/src/bundlers"
	"dotins.eu.org/financeIB/src/extractors"
	"dotins.eu.org/financeIB/src/utils"
	"github.com/samber/lo"
	"github.com/sqweek/dialog"
)

func main() {
	var inputFileName, _ = lo.Nth(os.Args, 1)
	var outputFileName, _ = lo.Nth(os.Args, 2)
	var extension = path.Ext(inputFileName)

	defer alert()

	if inputFileName == "" {
		inputFileName, outputFileName = utils.Collect()
		extension = path.Ext(inputFileName)
	}

	if extension == ".ofx" {
		bundlers.XLSXInit(extractors.OFXInit(inputFileName), outputFileName)
		dialog.Message("Extrato coletado com sucesso").Info()
		return
	}

	if extension == ".xlsx" || extension == ".xls" {
		bundlers.PDFInit(extractors.XLSXInit(inputFileName, outputFileName))
		dialog.Message("Relatório gerado com sucesso").Info()
		return
	}
}

func alert() {
	if r := recover(); r != nil {
		if n, _ := lo.Nth(os.Args, 1); n == "" {
			dialog.Message(fmt.Sprint(r) + "\n" + string(debug.Stack())).Title("Ocorreu um erro ao processar").Error()
		} else {
			fmt.Println("---- Ocorreu um erro ao processar ---")
			fmt.Println(fmt.Sprint(r)+"\n", string(debug.Stack())+"\n")
			fmt.Println("---- PRESSIONE QUALQUER TECLA PARA SAIR ---")
			fmt.Scanln()
		}
	}
}
