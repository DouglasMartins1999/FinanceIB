package extractors

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"dotins.eu.org/financeIB/src/interfaces"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

var (
	sheetName = "Demonstrativo"
	regexDot  = regexp.MustCompile("\\.")
	regexTime = regexp.MustCompile("\\d{5}")
)

func XLSXInit(filename string, outputname string) interfaces.Report {
	var sheet = lo.Must(excelize.OpenFile(filename))
	var rows = lo.Must(sheet.GetRows(sheetName))[2:]

	var inFlow = parseRows(rows, true, 0, 1, 2)
	var outFlow = parseRows(rows, false, 3, 4, 5)

	var sheetTotalInFlow = toAmountInCents(lo.Must(sheet.GetCellValue(sheetName, "I16")))
	var sheetTotalOutFlow = toAmountInCents(lo.Must(sheet.GetCellValue(sheetName, "I17")))
	var sheetMonthBalance = toAmountInCents(lo.Must(sheet.GetCellValue(sheetName, "I18")))

	var calcTotalInFlow = totalInCents(inFlow)
	var calcTotalOutFlow = totalInCents(outFlow)

	if lo.Must(strconv.ParseBool(calcOrGet(sheet, "I8"))) {
		inFlow = groupEquals(inFlow)
		outFlow = groupEquals(outFlow)
	}

	return interfaces.Report{
		ReportMonth: lo.Must(strconv.ParseInt(calcOrGet(sheet, "I4"), 10, 64)) - 1,
		ReportYear:  calcOrGet(sheet, "I5"),

		WrapPage: lo.Must(strconv.ParseBool(calcOrGet(sheet, "I7"))),
		FileName: lo.Ternary(outputname != "", outputname, strings.ReplaceAll(filename, ".xlsx", ".pdf")),

		BeforeBalance:  toAmountInCents(lo.Must(sheet.GetCellValue(sheetName, "I11"))),
		CurrentBalance: toAmountInCents(lo.Must(sheet.GetCellValue(sheetName, "I12"))),

		CashTime: toTime(calcOrGet(sheet, "I13")),
		SignTime: toTime(calcOrGet(sheet, "I6")),

		TotalInFlow:  lo.Must(lo.Coalesce(sheetTotalInFlow, calcTotalInFlow)),
		TotalOutFlow: lo.Must(lo.Coalesce(sheetTotalOutFlow, calcTotalOutFlow)),
		MonthBalance: lo.Must(lo.Coalesce(sheetMonthBalance, calcTotalInFlow.Sub(calcTotalOutFlow))),

		Statements: interfaces.Statement{
			InFlow:  inFlow,
			OutFlow: outFlow,
		},
	}
}

func parseRows(rows [][]string, isCredit bool, dateIndex int, amountIndex int, descIndex int) []interfaces.Transaction {
	var arrIndex = lo.Max([]int{dateIndex, amountIndex, descIndex})
	var tempRows = lo.Filter(rows, func(r []string, _ int) bool {
		return len(r) >= arrIndex && r[dateIndex] != "" && r[amountIndex] != ""
	})

	return lo.Map(tempRows, func(r []string, _ int) interfaces.Transaction {
		return interfaces.Transaction{
			IsCredit:    isCredit,
			CreatedAt:   toTime(r[dateIndex]),
			Amount:      lo.Must(decimal.NewFromString(r[amountIndex])),
			Description: r[descIndex],
		}
	})
}

func toAmountInCents(value string) (result decimal.Decimal) {
	if value != "" {
		return lo.Must(decimal.NewFromString(value))
	}

	return result
}

func toTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	if regexTime.FindString(value) == "" {
		return lo.Must(time.Parse("01-02-06", value))
	}

	return time.UnixMilli((lo.Must(strconv.ParseInt(strings.Split(value, ".")[0], 10, 64)) - 25569) * 86400000)
}

func groupEquals(trxs []interfaces.Transaction) []interfaces.Transaction {
	var groups = lo.GroupBy(trxs, func(t interfaces.Transaction) string {
		return t.Description
	})

	var result = lo.MapToSlice(groups, func(key string, values []interfaces.Transaction) interfaces.Transaction {
		return interfaces.Transaction{
			Description: key,
			Amount:      totalInCents(values),
			CreatedAt: lo.MaxBy(values, func(v interfaces.Transaction, max interfaces.Transaction) bool {
				return v.CreatedAt.After(max.CreatedAt)
			}).CreatedAt,
		}
	})

	sort.Slice(result, func(a int, b int) bool {
		return result[a].CreatedAt.Before(result[b].CreatedAt)
	})

	return result
}

func calcOrGet(sheet *excelize.File, cell string) string {
	var getValue, _ = sheet.GetCellValue(sheetName, cell)

	if getValue != "" {
		return getValue
	}

	return lo.Must(sheet.CalcCellValue(sheetName, cell))
}
