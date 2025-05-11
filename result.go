// Package result for standard result for REST-API responses
//
//	Author: Elizalde G. Baguinon
//	Created: October 17, 2019
package result

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/stdutil/log"
)

// Status items
const (
	OK        Status = `OK`
	EXCEPTION Status = `EXCEPTION`
	VALID     Status = `VALID`
	INVALID   Status = `INVALID`
	YES       Status = `YES`
	NO        Status = `NO`
)

// InitResult - initialize result for API query. This is the recommended initialization of this object.
// The variadic arguments of InitResultOption will modify default status.
// Depending on the current status (default is EXCEPTION), the message type is automatically set to that type
func InitResult(opts ...InitResultOption) Result {
	res := Result{
		Status:  string(EXCEPTION),
		ln:      log.Log{},
		osIsWin: runtime.GOOS == "windows",
	}
	res.Messages = make([]string, 0)
	irp := InitResultParam{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&irp)
	}
	if irp.Status != "" {
		res.Status = string(irp.Status)
	}
	res.SetPrefix(irp.Prefix)
	res.eventVerb = irp.EventVerb
	res.initFc = irp.InitialFocusID // preserve initial focus control
	res.SetFocusControl(res.initFc, false)

	// Auto-detect function that called this function
	if pc, _, _, ok := runtime.Caller(1); ok {
		if details := runtime.FuncForPC(pc); details != nil {
			nm := details.Name()
			if pos := strings.LastIndex(nm, `.`); pos != -1 {
				nm = nm[pos+1:]
			}
			res.Operation = strings.ToLower(nm)
			if res.eventVerb == "" {
				res.eventVerb = res.Operation
			}
		}
	}

	if irp.Message != "" {
		msg := irp.Message
		if irp.UseOperationInMsg && res.Operation != "" {
			msg = fmt.Sprintf(" %s: %s", res.Operation, irp.Message)
		}
		switch irp.Status {
		case OK, VALID, YES:
			res.AddInfo("%s", msg)
		case EXCEPTION, INVALID, NO:
			res.AddError("%s", msg)
		default:
			res.ln.AddAppMsg(msg)
		}
	}

	return res
}

// MessageManager returns the internal message manager
func (r *Result) MessageManager() *log.Log {
	return &r.ln
}

// Return sets the current status of a result
func (r *Result) Return(status Status) Result {
	r.Status = string(status)
	return *r
}

// OK returns true if the status is OK.
func (r *Result) OK() bool {
	return r.Status == string(OK)
}

// Error returns true if the status is EXCEPTION.
func (r *Result) Error() bool {
	return r.Status == string(EXCEPTION)
}

// Valid returns true if the status is VALID.
func (r *Result) Valid() bool {
	return r.Status == string(VALID)
}

// Invalid returns true if the status is INVALID.
func (r *Result) Invalid() bool {
	return r.Status == string(INVALID)
}

// Yes returns true if the status is YES.
func (r *Result) Yes() bool {
	return r.Status == string(YES)
}

// No returns true if the status is No.
func (r *Result) No() bool {
	return r.Status == string(NO)
}

// AddInfo adds a formatted information message and returns itself
func (r *Result) AddInfo(fmtMsg string, a ...interface{}) Result {
	msg := fmtMsg
	if len(a) > 0 {
		msg = fmt.Sprintf(fmtMsg, a...)
	}
	if r.useOperationInMsg && r.Operation != "" {
		msg = fmt.Sprintf(" %s: ", r.Operation) + msg
	}
	r.ln.AddInfo(msg)
	r.updateMessage()
	return *r
}

// AddWarning adds a formatted warning message and returns itself
func (r *Result) AddWarning(fmtMsg string, a ...interface{}) Result {
	msg := fmtMsg
	if len(a) > 0 {
		msg = fmt.Sprintf(fmtMsg, a...)
	}
	if r.useOperationInMsg && r.Operation != "" {
		msg = fmt.Sprintf(" %s: ", r.Operation) + msg
	}
	r.ln.AddWarning(msg)
	r.updateMessage()
	return *r
}

// AddError adds a formatted error message and returns itself
func (r *Result) AddError(fmtMsg string, a ...interface{}) Result {
	msg := fmtMsg
	if len(a) > 0 {
		msg = fmt.Sprintf(fmtMsg, a...)
	}
	if r.useOperationInMsg && r.Operation != "" {
		msg = fmt.Sprintf(" %s: ", r.Operation) + msg
	}
	r.ln.AddError(msg)
	r.updateMessage()
	return *r
}

// AddErr adds a error-typed value and returns itself.
func (r *Result) AddErr(err error) Result {
	r.AddError("%s", err)
	return *r
}

// AddSuccess adds a formatted success message and returns itself
func (r *Result) AddSuccess(fmtMsg string, a ...interface{}) Result {
	msg := fmtMsg
	if len(a) > 0 {
		msg = fmt.Sprintf(fmtMsg, a...)
	}
	if r.useOperationInMsg && r.Operation != "" {
		msg = fmt.Sprintf(" %s: ", r.Operation) + msg
	}
	r.ln.AddSuccess(msg)
	r.updateMessage()
	return *r
}

// AddRawMsg adds a formatted application message and returns itself
func (r *Result) AddRawMsg(fmtMsg string, a ...interface{}) Result {
	msg := fmtMsg
	if len(a) > 0 {
		msg = fmt.Sprintf(fmtMsg, a...)
	}
	r.ln.AddAppMsg(msg)
	r.updateMessage()
	return *r
}

// AddErrWithAlt adds an error-typed value, and an alternate error
// message if the err happens to be nil. It returns itself.
func (r *Result) AddErrWithAlt(err error, altMsg string, altMsgValues ...any) Result {
	if err != nil {
		return r.AddErr(err)
	}
	if altMsg != "" {
		return r.AddError(altMsg, altMsgValues...)
	}
	return *r
}

// AddErrorWithAlt appends the messages of a Result.
// And an alternative message if the Result is other than OK or VALID status.
func (r *Result) AddErrorWithAlt(rs Result, altMsg string, altMsgValues ...any) Result {
	if !(rs.OK() || rs.Valid()) {
		for _, n := range rs.ln.Notes() {
			r.ln.Append(n)
		}
		r.updateMessage()
		return *r
	}
	if altMsg == "" {
		return *r
	}
	r.ln.Append(
		log.LogInfo{
			Type:    log.Error,
			Message: fmt.Sprintf(altMsg, altMsgValues...),
			Prefix:  r.ln.Prefix,
		})
	r.updateMessage()
	return *r
}

// AppendErr copies the messages of the Result parameter and append an error message
func (r *Result) AppendErr(rs Result, err error) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddErr(err)
}

// AppendErrorf copies the messages of the Result parameter and append a formatted error message
func (r *Result) AppendError(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddError(fmtMsg, a...)
}

// AppendInfof copies the messages of the Result parameter and append a formatted information message
func (r *Result) AppendInfo(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddInfo(fmtMsg, a...)
}

// AppendWarning copies the messages of the Result parameter and append a formatted warning message
func (r *Result) AppendWarning(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddWarning(fmtMsg, a...)
}

// Stuff adds or appends the messages of a Result.
func (r *Result) Stuff(rs Result) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	r.updateMessage()
	return *r
}

// EventID returns the past tense of Operation
func (r *Result) EventID() string {
	ev := r.eventVerb
	if ev == "" {
		return "unknown"
	}
	// simple past tenser
	if !strings.HasSuffix(ev, "e") {
		return ev + "ed"
	}
	return ev + "d"
}

// MessagesToString returns all errors in a string separated by carriage return and/or line feed
func (r *Result) MessagesToString() string {
	// The r.Messages might have been unmarshalled from result bytes so we should process.
	if len(r.Messages) == 1 {
		return r.Messages[0]
	}
	if len(r.Messages) > 1 {
		lf := "\n"
		if r.osIsWin {
			lf = "\r\n"
		}
		sb := strings.Builder{}
		for _, v := range r.Messages {
			vlf := v + lf // prevents escape to the heap
			sb.Write([]byte(vlf))
		}
		return sb.String()
	}
	return r.ln.ToString()
}

// SetPrefix changes the prefix
func (r *Result) SetPrefix(pfx string) {
	r.ln.Prefix = pfx
	r.Prefix = pfx
}

// SetFocusControl sets the control to focus when an issue is encountered
//
// When appendOnly is true, it only appends to the present FocusControl field
// To reset the focus control, call ResetFocusControl method
func (r *Result) SetFocusControl(ctrl string, appendOnly bool) {
	if r.FocusControl == nil {
		r.FocusControl = new(string)
	}
	if !appendOnly {
		r.initFc = ctrl
		r.FocusControl = &ctrl
		return
	}
	*r.FocusControl = r.initFc + "_" + ctrl
}

// ResetFocusControl resets the focus control to the initial value
func (r *Result) ResetFocusControl() {
	r.FocusControl = &r.initFc
}

// RowsAffectedInfo - a function to simplify adding information for rows affected
func (r *Result) RowsAffectedInfo(rowsaff int64) {
	if rowsaff != 0 {
		r.AddInfo("%d rows affected", rowsaff)
	} else {
		r.AddInfo("No rows affected")
	}
}

func (r *Result) updateMessage() {
	// get current notes to update the messages array
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
}
