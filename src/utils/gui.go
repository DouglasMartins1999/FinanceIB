package utils

import (
	"fmt"
	"path"

	"github.com/samber/lo"
	"github.com/sqweek/dialog"
)

func Collect() (input string, output string) {
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
