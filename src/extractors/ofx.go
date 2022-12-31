package extractors

import (
	"os"
	"strings"

	"dotins.eu.org/financeIB/src/interfaces"
	"github.com/aclindsa/ofxgo"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var (
	startBalance = "SALDO ANTERIOR"

	filterExps = []string{
		"SDO CTA/APL",
		"SALDO ANTERIOR",
		"SALDO DO DIA",
	}
)

func OFXInit(filename string) interfaces.Statement {
	var file = lo.Must(os.Open(filename))
	var parsed = parse(file)

	return parsed
}

func parse(file *os.File) interfaces.Statement {
	var content = lo.Must(ofxgo.ParseResponse(file))
	var bankInfo = content.Bank[0].(*ofxgo.StatementResponse)
	var transactions = structTransactions(filterTransactions(bankInfo.BankTranList.Transactions))
	var balance, _ = decimal.NewFromString(bankInfo.BalAmt.String())

	return interfaces.Statement{
		InFlow:  getInFlow(transactions),
		OutFlow: getOutFlow(transactions),
		Balance: interfaces.Balance{
			Started:   startedBalance(bankInfo.BankTranList.Transactions),
			Current:   balance.Div(decimal.NewFromInt(100)),
			UpdatedAt: bankInfo.DtAsOf.Time,
		},
	}
}

func totalInCents(trxs []interfaces.Transaction) decimal.Decimal {
	return lo.Reduce(trxs, func(agg decimal.Decimal, v interfaces.Transaction, _ int) decimal.Decimal {
		return decimal.Sum(agg, v.Amount)
	}, decimal.Zero)
}

func getInFlow(trxs []interfaces.Transaction) []interfaces.Transaction {
	return lo.Filter(trxs, func(t interfaces.Transaction, _ int) bool {
		return t.IsCredit
	})
}

func getOutFlow(trxs []interfaces.Transaction) []interfaces.Transaction {
	return lo.Filter(trxs, func(t interfaces.Transaction, _ int) bool {
		return !t.IsCredit
	})
}

func filterTransactions(trxs []ofxgo.Transaction) []ofxgo.Transaction {
	return lo.Filter(trxs, func(t ofxgo.Transaction, _ int) bool {
		return lo.EveryBy(filterExps, func(f string) bool {
			return !strings.Contains(t.Memo.String(), f)
		})
	})
}

func structTransactions(trxs []ofxgo.Transaction) []interfaces.Transaction {
	return lo.Map(trxs, func(trx ofxgo.Transaction, index int) interfaces.Transaction {
		amount, _ := decimal.NewFromString(trx.TrnAmt.String())

		return interfaces.Transaction{
			Amount:      amount.Abs(),
			IsCredit:    amount.Sign() > 0,
			Description: string(trx.Memo),
			CreatedAt:   trx.DtPosted.Time,
		}
	})
}

func startedBalance(trxs []ofxgo.Transaction) decimal.Decimal {
	var elm, found = lo.Find(trxs, func(t ofxgo.Transaction) bool {
		return t.Memo.String() == startBalance
	})

	if found {
		value, _ := decimal.NewFromString(elm.TrnAmt.String())
		return value
	}

	return decimal.Zero
}
