package endpoints

type Request interface {
	validate() error
}

// CountRequest collects the request parameters for the Count method.
type CountRequest struct {
}

func (r CountRequest) validate() error {
	return nil // TBA
}
