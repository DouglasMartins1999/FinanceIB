package bundlers

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"dotins.eu.org/financeIB/src/interfaces"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

var (
	//go:embed resources/template.xlsx
	template string

	sheetName   = "Demonstrativo"
	sheetRowPad = 3
)

func XLSXInit(settlement interfaces.Statement) {
	var sheet = lo.Must(excelize.OpenReader(strings.NewReader(template)))
	var saveName = "./report_" + time.Now().Format("2006-01-02_15.04.05") + ".xlsx"

	fillValues(sheet, settlement.InFlow, "A", "B", "C")
	fillValues(sheet, settlement.OutFlow, "D", "E", "F")

	sheet.SetCellDefault(sheetName, "I10", settlement.Balance.Started.String())
	sheet.SetCellDefault(sheetName, "I11", settlement.Balance.Current.String())
	sheet.SetCellValue(sheetName, "I12", settlement.Balance.UpdatedAt)

	sheet.UpdateLinkedValue()
	sheet.SaveAs(saveName)
}

func fillValues(sheet *excelize.File, trxs []interfaces.Transaction, dateCol string, valueCol string, descCol string) {
	for index, t := range trxs {
		var position = fmt.Sprint(index + sheetRowPad)

		sheet.SetCellValue(sheetName, dateCol+position, t.CreatedAt)
		sheet.SetCellValue(sheetName, descCol+position, t.Description)
		sheet.SetCellDefault(sheetName, valueCol+position, t.Amount.String())
	}

}
