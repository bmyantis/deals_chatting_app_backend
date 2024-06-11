package data

type BaseResponse struct {
	ProcessStatus string `validate:"oneof=S F U" json:"processStatus"`
	TxnRef        string `validate:"required,min=1,max=200" json:"txnRef"`
}

type BaseErrorResponse struct {
	BaseResponse
	Error string `validate:"required,min=1" json:"error"`
}
