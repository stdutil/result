package result

import (
	"encoding/json"
	"testing"

	l "github.com/stdutil/log"
)

func TestResultAnyReturnWithMsg(t *testing.T) {
	ra := ResultAny[any]{
		Result: InitResult(
			WithPrefix("PREFIX"),
		),
		Data: nil,
	}
	ra.ReturnWithMsg(OK, l.Warn, "This is a message upon return, %s", "Sir")

	b, err := json.Marshal(ra)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(string(b))
	}
}
