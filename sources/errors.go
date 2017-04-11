package sources

// THIS FILE IS FOR THE SPECIFIC API ERRORS
//Here goes the error codes
const (
	CodeExampleError1 = iota + 420
	CodeExampleError2
)

//Here goes the error messages, they can be hardcoded strings or you can pass '%s'
//to format some parts of the message later on runtime
const (
	MsgExampleError1 = "ERROR_SOMETHING_IS_WRONG DETAILS: '%s'"
	MsgExampleError2 = "ERROR_HARDCODED_MESSAGE"
)
