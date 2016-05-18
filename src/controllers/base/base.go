package base

import "gen"

// Controller is the base controller for the mqtt payload handler
type Controller struct {
}

// SimpleError return a field structure for error data
func (c *Controller) SimpleError(cid string, err error) *heliumgen.SimpleError {
	return &heliumgen.SimpleError{
		CorrelationId: cid,
		Error:         err.Error(),
	}
}

// MultiError is the complex version of error
func (c *Controller) MultiError(cid string, errs map[string]error) *heliumgen.ComplexError {
	errsSS := make(map[string]string, len(errs))
	for i := range errs {
		errsSS[i] = errs[i].Error()
	}
	return &heliumgen.ComplexError{
		CorrelationId: cid,
		Errors:        errsSS,
	}
}
