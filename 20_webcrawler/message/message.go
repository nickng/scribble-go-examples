package message

// URL is a URL reply for Master.
type URL struct {
	URL string
}

func (URL) GetOp() string {
	return "URL"
}

// Parse is a fragment of parsed HTML.
type Parse struct {
	Tokens []string
}

func (Parse) GetOp() string {
	return "Parse"
}

// Index is the index info of parsed HTML.
type Index struct {
}

func (Index) GetOp() string {
	return "Index"
}
