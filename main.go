package main

import (
	"fmt"
	"os"
	"path"

	"dotins.eu.org/financeIB/src/bundlers"
	"dotins.eu.org/financeIB/src/extractors"
	"github.com/samber/lo"
	"github.com/sqweek/dialog"
)

func main() {
	var inputFileName, _ = lo.Nth(os.Args, 1)
	var outputFileName, _ = lo.Nth(os.Args, 2)
	var extension = path.Ext(inputFileName)

	if inputFileName == "" {
		inputFileName, outputFileName = collect()
		extension = path.Ext(inputFileName)
	}

	fmt.Println(inputFileName, outputFileName, extension)

	if extension == ".ofx" {
		bundlers.XLSXInit(extractors.OFXInit(inputFileName), outputFileName)
		return
	}

	if extension == ".xlsx" || extension == ".xls" {
		bundlers.PDFInit(extractors.XLSXInit(inputFileName, outputFileName))
		return
	}
}

func collect() (input string, output string) {
	var status error
	var extension string
	var saveOutput = false

	dialog.Message("Somente arquivos .OFX ou .XLSX são aceitos").Title("Escolha o arquivo do seu computador").Info()
	for input == "" {
		input, status = dialog.File().Filter("Arquivo .OFX ou .XLSX", "ofx", "xlsx", "xls").Load()
		extension = path.Ext(input)

		if status == dialog.Cancelled {
			exit := dialog.Message("Deseja sair do programa?").YesNo()

			if exit {
				return
			}
		}
	}

	saveOutput = dialog.Message("Quer escolher onde salvar seu arquivo?").YesNo()

	for output == "" && saveOutput {
		output, status = dialog.File().Filter("Salvar Relatório", lo.Ternary(extension == ".ofx", "xlsx", "pdf")).Title("Exportar Relatório").Save()

		if status == dialog.Cancelled {
			exit := dialog.Message("Deseja usar o nome e localização padrão?").YesNo()

			if exit {
				return
			}
		}
	}

	return
}
