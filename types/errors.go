package types

import "fmt"

type ErrorCode string

const (
	ErrCodeNoLiquidityInHot    ErrorCode = "no_liquidity_in_hot"
	ErrCodeInsufficientBalance ErrorCode = "insufficient_balance"
	ErrCodeInvalidRequest      ErrorCode = "invalid_request"
	ErrCodeNotFound            ErrorCode = "not_found"
)

type CodedError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e CodedError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("error: %s", e.Code)
	}
	return fmt.Sprintf("error: %s: %s", e.Code, e.Message)
}
