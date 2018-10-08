package messages

// Data is the default payload type.
type Data struct {
	V int
}

// GetOp returns the message label.
func (Data) GetOp() string { return "Data" }
