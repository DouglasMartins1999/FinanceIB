package interfaces

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Amount      decimal.Decimal
	IsCredit    bool
	Description string
	CreatedAt   time.Time
}

type Balance struct {
	Started   decimal.Decimal
	Current   decimal.Decimal
	UpdatedAt time.Time
}

type Statement struct {
	InFlow  []Transaction
	OutFlow []Transaction
	Balance Balance
}

type Report struct {
	ReportMonth int64
	ReportYear  string
	SignTime    time.Time
	CashTime    time.Time

	WrapPage bool
	FileName string

	BeforeBalance  decimal.Decimal
	CurrentBalance decimal.Decimal

	TotalInFlow  decimal.Decimal
	TotalOutFlow decimal.Decimal

	MonthBalance decimal.Decimal

	Statements Statement
}
