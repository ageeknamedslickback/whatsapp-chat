package rest

// Handlers defines our REST layer interface
type Handlers interface{}

// Rest ..
type Rest struct {
}

// NewRestHandlers ..
func NewRestHandlers() *Rest {
	return &Rest{}
}
