package msg

// Meta represents metadata retrieved from server.
// Contains full Size of the file to retrieve.
type Meta struct {
	Size int
}

// GetOp returns the message label Meta.
func (Meta) GetOp() string {
	return "Meta"
}

// Job is the fetch task issued by Master to Fetchers.
type Job struct {
	URL       string
	RangeFrom int
	RangeTo   int
}

// GetOp returns the message label Job.
func (Job) GetOp() string {
	return "Job"
}

// Data is the fetch result which should be a fragment
// of the overall data.
type Data struct {
	Data []byte
}

// GetOp returns the message label Data.
func (Data) GetOp() string {
	return "Data"
}

// Done is a fetch complete signal.
type Done struct {
}

// GetOp returns the message label Done.
func (Done) GetOp() string {
	return "Done"
}
