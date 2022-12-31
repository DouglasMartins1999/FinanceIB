package bundlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"dotins.eu.org/financeIB/src/interfaces"
	"github.com/goodsign/monday"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/leekchan/accounting"
	"github.com/samber/lo"
)

var (
	entityName = "Igreja Batista Regular em Jd. São Jorge"
	document   = "CNPJ 43.100.650/0001-11"
	address    = "Rua Gilberto Duarte de Azevedo, 392 - São Paulo/SP | CEP 05567-070"
	signCity   = "São Paulo"

	reportTitle = "Demonstrativo Financeiro"
	refersText  = "Referente a %s de %s"

	beforeBalanceText  = "Saldo Anterior"
	currentBalanceText = "Disponibilidade"
	monthBalanceText   = "Variação de Caixa (%s - %s)"
	resultBalanceText  = "Saldo Anterior + Variação no Período"

	resultText  = "Demonstrativo de Recursos"
	inFlowText  = "Entradas"
	OutFlowText = "Saídas"
	totalText   = "Total de %s"

	tableColumnOneText   = "Data"
	tableColumnTwoText   = "Descrição"
	tableColumnThreeText = "Valor"

	signatureOne   = "1º Tesoureiro"
	signatureTwo   = "2º Tesoureiro"
	signatureThree = "Comissão de Exame de Contas"

	breakResultPage = true

	verticalMargin   = 15.0
	horizontalMargin = 15.0

	timeLocale monday.Locale = monday.LocalePtBR
)

var ai18n = accounting.Accounting{
	Precision: 2,
	Symbol:    "R$ ",
	Thousand:  ".",
	Decimal:   ",",
}

var header = []string{
	tableColumnOneText,
	tableColumnTwoText,
	tableColumnThreeText,
}

func PDFInit(report interfaces.Report) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(horizontalMargin, verticalMargin, horizontalMargin)
	m.SetDefaultFontFamily(consts.Helvetica)

	cashTime := monday.Format(report.CashTime, "02/01", timeLocale)
	signTime := monday.Format(report.SignTime, "02 de January de 2006", timeLocale)
	reportMonth := monday.GetLongMonths(timeLocale)[report.ReportMonth]

	m.SetFirstPageNb(1)
	m.RegisterFooter(func() {
		m.Row(0, func() {
			m.Col(12, func() {
				m.Text(strconv.Itoa(m.GetCurrentPage()), props.Text{
					Align: consts.Right,
					Size:  10,
					Top:   10,
				})
			})
		})
	})

	m.Row(8, func() {
		m.Text(strings.ToUpper(entityName), props.Text{
			Size:  16,
			Style: consts.Bold,
			Align: consts.Center,
		})
	})
	m.Row(5.5, func() {
		m.Text(address, props.Text{
			Size:  10,
			Align: consts.Center,
		})
	})
	m.Row(8, func() {
		m.Text(document, props.Text{
			Size:  10,
			Align: consts.Center,
		})
	})

	m.Row(16.5, func() {
		m.Text(strings.ToUpper(reportTitle), props.Text{
			Size:  14,
			Style: consts.Bold,
			Align: consts.Center,
			Top:   10,
		})
	})
	m.Row(15, func() {
		m.Text(fmt.Sprintf(refersText, reportMonth, report.ReportYear), props.Text{
			Size:  12,
			Style: consts.Normal,
			Align: consts.Center,
		})
	})

	m.Line(2, props.Line{
		Color: color.Color{
			Red: 200, Green: 200, Blue: 200,
		},
	})

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text(beforeBalanceText+":", props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   2,
			})
		})
		m.Col(6, func() {
			m.Text(ai18n.FormatMoney(report.BeforeBalance), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Right,
				Top:   2,
			})
		})
	})

	m.Line(2, props.Line{
		Color: color.Color{
			Red: 200, Green: 200, Blue: 200,
		},
	})

	m.Row(22, func() {
		m.Text(strings.ToUpper(inFlowText), props.Text{
			Size:  14,
			Style: consts.Bold,
			Align: consts.Center,
			Top:   10,
		})
	})

	m.TableList(header, structureTransactions(report.Statements.InFlow), props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{2, 8, 2},
			Size:      12,
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{2, 8, 2},
			Size:      11,
		},
		VerticalContentPadding: 3,
		Line:                   true,
		LineProp: props.Line{
			Color: color.Color{
				Red: 200, Green: 200, Blue: 200,
			},
		},
	})

	m.Row(12, func() {
		m.Col(10, func() {
			m.Text(fmt.Sprintf(totalText, inFlowText), props.Text{
				Size:  12,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   3,
			})
		})
		m.Col(2, func() {
			m.Text(ai18n.FormatMoney(report.TotalInFlow), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   3,
			})
		})
	})

	m.Row(22, func() {
		m.Col(0, func() {
			m.Text(strings.ToUpper(OutFlowText), props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Center,
				Top:   10,
			})
		})
	})

	m.TableList(header, structureTransactions(report.Statements.OutFlow), props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{2, 8, 2},
			Size:      12,
		},
		HeaderProp: props.TableListContent{
			GridSizes: []uint{2, 8, 2},
			Size:      11,
		},
		VerticalContentPadding: 3,
		Line:                   true,
		LineProp: props.Line{
			Color: color.Color{
				Red: 200, Green: 200, Blue: 200,
			},
		},
	})

	m.Row(12, func() {
		m.Col(10, func() {
			m.Text(fmt.Sprintf(totalText, OutFlowText), props.Text{
				Size:  12,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   3,
			})
		})
		m.Col(2, func() {
			m.Text(ai18n.FormatMoney(report.TotalOutFlow), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   3,
			})
		})
	})

	if report.WrapPage {
		m.AddPage()
	}

	m.Row(32, func() {
		m.Col(0, func() {
			m.Text(strings.ToUpper(resultText), props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Center,
				Top:   16,
			})
		})
	})

	m.Row(10, func() {
		m.Col(10, func() {
			m.Text(fmt.Sprintf(monthBalanceText, inFlowText, OutFlowText), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   2,
			})
		})
		m.Col(2, func() {
			m.Text(ai18n.FormatMoney(report.MonthBalance), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Right,
				Top:   2,
			})
		})
	})

	m.Row(10, func() {
		m.Col(10, func() {
			m.Text(resultBalanceText, props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   2,
			})
		})
		m.Col(2, func() {
			m.Text(ai18n.FormatMoney(report.BeforeBalance.Add(report.MonthBalance)), props.Text{
				Size:  14,
				Style: consts.Normal,
				Align: consts.Right,
				Top:   2,
			})
		})
	})

	m.Line(2, props.Line{
		Color: color.Color{
			Red: 200, Green: 200, Blue: 200,
		},
	})

	m.Row(10, func() {
		m.Col(10, func() {
			m.Text(fmt.Sprintf("%s [%s]:", currentBalanceText, cashTime), props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   2,
			})
		})
		m.Col(2, func() {
			m.Text(ai18n.FormatMoney(report.CurrentBalance), props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Right,
				Top:   2,
			})
		})
	})

	m.Line(2, props.Line{
		Color: color.Color{
			Red: 200, Green: 200, Blue: 200,
		},
	})

	m.Row(50, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("%s, %s", signCity, signTime), props.Text{
				Top:   14,
				Style: consts.Normal,
				Size:  11,
			})
		})
	})

	m.Row(0, func() {
		m.Col(4, func() {
			m.Signature(signatureOne, props.Font{Size: 8})
		})

		m.Col(4, func() {
			m.Signature(signatureTwo, props.Font{Size: 8})
		})

		m.Col(4, func() {
			m.Signature(signatureThree, props.Font{Size: 8})
		})
	})

	var err = m.OutputFileAndClose("./" + report.FileName)

	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}
}

func structureTransactions(trxs []interfaces.Transaction) [][]string {
	return lo.Map(trxs, func(t interfaces.Transaction, _ int) []string {
		return []string{
			monday.Format(t.CreatedAt, "02 Jan", monday.LocalePtBR),
			t.Description,
			ai18n.FormatMoney(t.Amount),
		}
	})
}
