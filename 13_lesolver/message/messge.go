package message

type Data struct {
	V int
}

func (Data) GetOp() string {
	return "Data"
}
