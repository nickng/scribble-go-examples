package message

// Data is Data(int) signature.
type Data struct {
	V int
}

// GetOp returns the label.
func (*Data) GetOp() string { return "Data" }
