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
	startBalance = []string{"SALDO INICIAL", "SALDO ANTERIOR"}
	currBalance  = []string{"SALDO FINAL", "SALDO DO DIA"}

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

	return interfaces.Statement{
		InFlow:  getInFlow(transactions),
		OutFlow: getOutFlow(transactions),
		Balance: interfaces.Balance{
			Started:   entryValue(bankInfo.BankTranList.Transactions, startBalance, false),
			Current:   entryValue(bankInfo.BankTranList.Transactions, currBalance, true),
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

func entryValue(trxs []ofxgo.Transaction, term []string, isReverse bool) decimal.Decimal {
	var trx ofxgo.Transaction
	var found bool

	if isReverse {
		trx, _, found = lo.FindLastIndexOf(trxs, func(t ofxgo.Transaction) bool {
			return lo.Contains(term, t.Memo.String())
		})
	} else {
		trx, _, found = lo.FindIndexOf(trxs, func(t ofxgo.Transaction) bool {
			return lo.Contains(term, t.Memo.String())
		})
	}

	if found {
		value, _ := decimal.NewFromString(trx.TrnAmt.String())
		return value
	}

	return decimal.Zero
}
