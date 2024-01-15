package mnemo

type (
	MnemoError struct {
		Err error
	}
	ConnError struct {
		StatusCode int
		Err        error
	}
	RequestError struct {
		StatusCode int
		Err        error
	}
	StoreError struct {
		Err error
	}
)

func (e MnemoError) Error() string {
	return "Mnemo Error: " + e.Err.Error()
}
func (e ConnError) Error() string {
	return "Mnemo Conn Error: " + e.Err.Error()
}
func (e RequestError) Error() string {
	return "Mnemo Request Error: " + e.Err.Error()
}
func (e StoreError) Error() string {
	return "Mnemo Store Error: " + e.Err.Error()
}
