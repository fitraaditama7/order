package customerror

import "fmt"

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Errors  error  `json:"-"`
}

func (ce Err) Error() string {
	if ce.Errors != nil {
		message := fmt.Sprintf("%s - %s", ce.Message, ce.Errors.Error())
		return message
	}
	return ce.Message
}

func Error(code string, message string) error {
	return &Err{
		Code:    code,
		Message: message,
	}
}
