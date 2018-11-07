package message

// Initial is the initial Bid.
type Initial struct {
	Bid int
}

func (Initial) GetOp() string {
	return "InitialBid"
}

// Highest is the highest Bid.
type Highest struct {
	Bid int
}

func (Highest) GetOp() string {
	return "HighestBid"
}

// BidOrSkip is a bid by bidder.
type BidOrSkip struct {
	Bid     int
	MakeBid bool // false if Skip
}

func (BidOrSkip) GetOp() string {
	return "BidOrSkip"
}

// Winner is a signal for a winner.
type Winner struct {
	BidderID int
}

func (Winner) GetOp() string {
	return "Winner"
}
