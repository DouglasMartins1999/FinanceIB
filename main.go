package main

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"dotins.eu.org/financeIB/src/bundlers"
	"dotins.eu.org/financeIB/src/extractors"
	"github.com/samber/lo"
	"github.com/sqweek/dialog"
)

func main() {
	var inputFileName, _ = lo.Nth(os.Args, 1)
	var outputFileName, _ = lo.Nth(os.Args, 2)
	var extension = path.Ext(inputFileName)

	defer alert()

	if inputFileName == "" {
		inputFileName, outputFileName = collect()
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

func collect() (input string, output string) {
	var status error
	var extension string
	var saveOutput = false

	dialog.Message("Escolha o arquivo do seu computador").Info()

	for input == "" {
		input, status = dialog.File().Filter("Arquivo .OFX ou .XLSX", "ofx", "xlsx", "xls").Load()
		extension = path.Ext(input)

		if status == dialog.Cancelled && dialog.Message("Deseja sair do programa?").YesNo() {
			return
		}
	}

	saveOutput = dialog.Message("Quer escolher onde salvar seu arquivo?").YesNo()

	for output == "" && saveOutput {
		ext := lo.Ternary(extension == ".ofx", "xlsx", "pdf")
		title := lo.Ternary(extension == ".ofx", "Salvar Planilha", "Exportar Relatório")
		output, status = dialog.File().Filter(fmt.Sprintf("Arquivo .%s", ext), ext).Title(title).Save()

		if status == dialog.Cancelled && dialog.Message("Deseja usar o nome e localização padrão?").YesNo() {
			return
		}
	}

	return
}
