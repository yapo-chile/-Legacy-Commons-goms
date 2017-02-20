package main

type GenericError struct {
	ErrorMessage string
}

const ResponseOK = "OK"
const ResponseNOK = "NOK"

const (
	CodeInvalidParams = iota + 410
	CodeMissingParam
	CodeForbidenExtraParams
	CodeNegativeBalance
	CodeMalformedRequest
	CodeInvalidParamValue
	CodeConsumeCredits
	CodeInvalidDate
	CodeEditCredits
)

const (
	CodeDbConnection = iota + 510
	CodeInternalError
)

const (
	MsgInvalidParams       = "INVALID_PARAMS"
	MsgMissingParam        = "MISSING_PARAMETER {%s}"
	MsgForbidenExtraParams = "FORBIDEN_EXTRA_PARAMS"
	MsgNegativeBalance     = "NEGATIVE_BALANCE"
	MsgMalformedRequest    = "MALFORMED_REQUEST"
	MsgInvalidParamValue   = "INVALID_PARAM_VALUE {%s}"
	MsgDbConnection        = "DB_CONNECTION"
	MsgInternalError       = "INTERNAL_ERROR"
	MsgConsumeCredits      = "ERROR_CONSUME_CREDITS_%s"
	MsgInvalidDate         = "ERROR_DATE_NOT_VALID"
	MsgEditCredits         = "ERROR_EDIT_CREDITS_%s"
)
