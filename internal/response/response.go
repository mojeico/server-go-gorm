package response

//Response is used for static shape json return

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

//EmptyObj object is used when data doesnt want to be null on json
type EmptyObj struct{}

//BuildResponse method is to inject data value to dynamic success response
func BuildResponse(message string, data interface{}) *Response {
	return &Response{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

//BuildErrorResponse method is to inject data value to dynamic failed response

func BuildErrorResponse(message string, err string, data interface{}) *Response {

	return &Response{
		Status:  false,
		Message: message,
		Errors:  err,
		Data:    data,
	}
}
