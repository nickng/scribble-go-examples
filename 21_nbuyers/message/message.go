package message

import "time"

// Address is the delivery address.
type Address struct {
	Line1    string
	Line2    string
	Country  string
	PostCode string
}

func (Address) GetOp() string {
	return "Address"
}

// Date is the delivery date.
type Date struct {
	D time.Time
}

func (Date) GetOp() string {
	return "Date"
}
