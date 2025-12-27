package junocashd

import (
	"fmt"
)

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	if e == nil {
		return "junocashd: rpc error <nil>"
	}
	if e.Message == "" {
		return fmt.Sprintf("junocashd: rpc error %d", e.Code)
	}
	return fmt.Sprintf("junocashd: rpc error %d: %s", e.Code, e.Message)
}
