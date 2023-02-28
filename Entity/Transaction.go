package Entity

type ArticleCredit struct {
	Account
	Bilan
}

type ArticleDebit struct {
	Account
	Bilan
}

type Transaction struct {
	Credit ArticleCredit
	Debit  ArticleDebit
}
