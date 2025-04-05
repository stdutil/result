package result

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
