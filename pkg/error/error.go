package customerr

import "fmt"

var (
	ErrNoneOfSeedDataInserted = fmt.Errorf("error on inserting seed data")
	ErrWordsNotFound          = fmt.Errorf("words not found")
	ErrWordAlreadyExist       = fmt.Errorf("word already exist")
	ErrWordIDNotFound         = fmt.Errorf("word has an id as sent not found")
	ErrAuthorizationFailed    = fmt.Errorf("authorization failed")
	ErrInvalidParameter       = fmt.Errorf("invalid request")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) (ErrorResponse, int) {
	resp := ErrorResponse{Error: err.Error()}

	code := 500
	switch err {
	case ErrNoneOfSeedDataInserted:
		code = 500
	case ErrWordsNotFound:
		code = 404
	case ErrWordAlreadyExist:
		code = 409
	case ErrWordIDNotFound:
		code = 404
	case ErrAuthorizationFailed:
		code = 401
	case ErrInvalidParameter:
		code = 400
	}

	return resp, code
}
