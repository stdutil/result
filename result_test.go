package result

import (
	"encoding/json"
	"testing"

	l "github.com/stdutil/log"
)

func TestReturnWithMsg(t *testing.T) {
	r := InitResult(
		WithPrefix("PREFIX"),
	)
	r.ReturnWithMsg(OK, l.Warn, "This is a message upon return, %s", "Sir")

	b, err := json.Marshal(r)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(string(b))
	}
}
