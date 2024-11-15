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

// Status type
type Status string

// Status items
const (
	OK        Status = `OK`
	EXCEPTION Status = `EXCEPTION`
	VALID     Status = `VALID`
	INVALID   Status = `INVALID`
	YES       Status = `YES`
	NO        Status = `NO`
)

type (
	// Result - standard result structure
	Result struct {
		Messages     []string     `json:"messages"`                // Accumulated messages as a result from Add methods. Do not append messages using append()
		Status       string       `json:"status"`                  // OK, ERROR, VALID or any status
		Operation    string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
		TaskID       *string      `json:"task_id,omitempty"`       // ID of the task and of the result
		WorkerID     *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
		FocusControl *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
		Page         *int         `json:"page,omitempty"`          // Current Page
		PageCount    *int         `json:"page_count,omitempty"`    // Page Count
		PageSize     *int         `json:"page_size,omitempty"`     // Page Size
		Tag          *interface{} `json:"tag,omitempty"`           // Miscellaneous result
		Prefix       string       `json:"prefix,omitempty"`        // Prefix of the message to return
		ln           log.Log      // Internal note
		eventVerb    string       // event verb related to the name of the operation
		osIsWin      bool
	}
	// ResultAny struct with generic type data
	ResultAny[T any] struct {
		Result
		Data T `json:"data"`
	}
	// InitResultParam are optional parameters for initiating a Result
	InitResultParam struct {
		Status  Status // Initial status
		Prefix  string // Prefix
		Message string // Message
	}
	// InitResultOption for initial result parameters
	InitResultOption func(opt *InitResultParam) error
)

// WithStatus sets the status of the Result as an option
func WithStatus(status Status) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.Status = status
		return nil
	}
}

// WithPrefix sets the prefix of the Result as an option
func WithPrefix(prefix string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.Prefix = prefix
		return nil
	}
}

// WithMessage sets the message of the Result as an option
func WithMessage(message string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.Message = message
		return nil
	}
}

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
	if irp.Message != "" {
		switch irp.Status {
		case OK, VALID, YES:
			res.AddInfo(irp.Message)
		case EXCEPTION, INVALID, NO:
			res.AddError(irp.Message)
		default:
			res.ln.AddAppMsg(irp.Message)
		}
	}

	// Auto-detect function that called this function
	if pc, _, _, ok := runtime.Caller(1); ok {
		if details := runtime.FuncForPC(pc); details != nil {
			nm := details.Name()
			if pos := strings.LastIndex(nm, `.`); pos != -1 {
				nm = nm[pos+1:]
			}
			res.Operation = strings.ToLower(nm)
			res.eventVerb = res.Operation
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
	r.ln.AddInfo(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddWarning adds a formatted warning message and returns itself
func (r *Result) AddWarning(fmtMsg string, a ...interface{}) Result {
	r.ln.AddWarning(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddError adds a formatted error message and returns itself
func (r *Result) AddError(fmtMsg string, a ...interface{}) Result {
	r.ln.AddError(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddErr adds a error-typed value and returns itself.
func (r *Result) AddErr(err error) Result {
	r.AddError(err.Error())
	return *r
}

// AddSuccess adds a formatted success message and returns itself
func (r *Result) AddSuccess(fmtMsg string, a ...interface{}) Result {
	r.ln.AddSuccess(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddAppMsg adds a formatted application message and returns itself
func (r *Result) AddAppMsg(fmtMsg string, a ...interface{}) Result {
	r.ln.AddAppMsg(fmt.Sprintf(fmtMsg, a...))
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

// ToString adds a formatted error message and returns itself
func (r *Result) MessagesToString() string {
	// if r.Messages is not empty, it can be because it was unmarshalled from result bytes
	if len(r.Messages) > 0 {
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

// RowsAffectedInfo - a function to simplify adding information for rows affected
func (r *Result) RowsAffectedInfo(rowsaff int64) {
	if rowsaff != 0 {
		r.AddInfo(fmt.Sprintf("%d rows affected", rowsaff))
	} else {
		r.AddInfo("No rows affected")
	}
}

func (r *Result) updateMessage() {
	// get current notes to update the messages
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
}

// AddInfo adds an information message and returns itself
func (r *ResultAny[T]) AddInfo(fmtMsg string, a ...interface{}) ResultAny[T] {
	r.Result.AddInfo(fmtMsg, a...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddWarning adds a warning message and returns itself
func (r *ResultAny[T]) AddWarning(fmtMsg string, a ...interface{}) ResultAny[T] {
	r.Result.AddWarning(fmtMsg, a...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddError adds an error message and returns itself
func (r *ResultAny[T]) AddError(fmtMsg string, a ...interface{}) ResultAny[T] {
	r.Result.AddError(fmtMsg, a...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddErr adds a error-typed value and returns itself.
func (r *ResultAny[T]) AddErr(err error) ResultAny[T] {
	r.Result.AddErr(err)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddSuccess adds an success message and returns itself
func (r *ResultAny[T]) AddSuccess(fmtMsg string, a ...interface{}) ResultAny[T] {
	r.Result.AddSuccess(fmtMsg, a...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// Stuff adds or appends the messages of a Result.
func (r *ResultAny[T]) Stuff(rs Result) ResultAny[T] {
	r.Result.Stuff(rs)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddErrWithAlt adds an error-typed value, and an alternate error
// message if the err happens to be nil. It returns itself.
func (r *ResultAny[T]) AddErrWithAlt(err error, altMsg string, altMsgValues ...any) ResultAny[T] {
	r.Result.AddErrWithAlt(err, altMsg, altMsgValues...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// AddErrorWithAlt appends the messages of a Result.
// And an alternative message if the Result is other than OK or VALID status.
func (r *ResultAny[T]) AddErrorWithAlt(rs Result, altMsg string, altMsgValues ...any) ResultAny[T] {
	r.Result.AddErrorWithAlt(rs, altMsg, altMsgValues...)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}

// Return sets the current status of a result
func (r *ResultAny[T]) Return(status Status) ResultAny[T] {
	r.Result.Return(status)
	return ResultAny[T]{
		Result: r.Result,
		Data:   r.Data,
	}
}
