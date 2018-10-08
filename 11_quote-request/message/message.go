package message

// QuoteReq is a quote request sent to supplier or manufacturer.
type QuoteReq struct {
	ItemID int
}

func (QuoteReq) GetOp() string {
	return "QuoteReq"
}

// Quote is a reply of QuoteRequest
type Quote struct {
	ItemID int
	Quote  int
}

func (Quote) GetOp() string {
	return "Quote"
}
