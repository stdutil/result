package result

import "github.com/stdutil/log"

type (
	Status string
	// Result - standard result structure
	Result struct {
		Messages          []string     `json:"messages"`                // Accumulated messages as a result from Add methods. Do not append messages using append()
		Status            string       `json:"status"`                  // OK, ERROR, VALID or any status
		Operation         string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
		TaskID            *string      `json:"task_id,omitempty"`       // ID of the task and of the result
		WorkerID          *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
		FocusControl      *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
		Page              *int         `json:"page,omitempty"`          // Current Page
		PageCount         *int         `json:"page_count,omitempty"`    // Page Count
		PageSize          *int         `json:"page_size,omitempty"`     // Page Size
		Tag               *interface{} `json:"tag,omitempty"`           // Miscellaneous result
		Prefix            string       `json:"prefix,omitempty"`        // Prefix of the message to return
		ln                log.Log      // Internal note
		eventVerb         string       // event verb related to the name of the operation
		osIsWin           bool         // checks for OS to determine carriage return line feed
		useOperationInMsg bool         // use Operation value in messages
		initFc            string       // original focus control
	}
	// ResultAny struct with generic type data
	ResultAny[T any] struct {
		Result
		Data T `json:"data"`
	}
	// InitResultParam are optional parameters for initiating a Result
	InitResultParam struct {
		EventVerb         string // Custom event verb or id
		Status            Status // Initial status
		Prefix            string // Prefix
		Message           string // Message
		InitialFocusID    string // Initial Focus Control id
		UseOperationInMsg bool   // Use Operation tag in messages
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
func WithPrefix(pfx string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.Prefix = pfx
		return nil
	}
}

// WithMessage sets the message of the Result as an option
func WithMessage(msg string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.Message = msg
		return nil
	}
}

// WithFocusControl sets the message of the Result as an option
func WithFocusControl(focusId string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.InitialFocusID = focusId
		return nil
	}
}

// WithEventVerb sets the custom event verb of the Result as an option
func WithEventVerb(eventVerb string) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.EventVerb = eventVerb
		return nil
	}
}

// UseOperationInMessage sets to include the Operation tag in messages
func UseOperationInMessage(on bool) InitResultOption {
	return func(irp *InitResultParam) error {
		irp.UseOperationInMsg = on
		return nil
	}
}
