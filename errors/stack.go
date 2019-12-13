package errors

type Stack struct {
	Error
	Stack []Stack `json:"stack,omitempty"`
}

type StackOptions struct {
	excludePath     bool
	excludeInternal bool
}

var (
	FullStack = &StackOptions{
		excludePath:     false,
		excludeInternal: false,
	}
	InfoStack = &StackOptions{
		excludePath:     true,
		excludeInternal: true,
	}
)

func BuildStack(err error, opts *StackOptions) []Stack {
	stack := Stack{}

	switch err := err.(type) {
	case Errors:
		stacks := make([]Stack, 0)
		for _, err := range err {
			stacks = append(stacks, BuildStack(err, opts)...)
		}
		return stacks
	case Error:
		if opts.excludeInternal && err.Type == Internal {
			return []Stack{}
		}

		stack.Error = err

		if opts.excludePath {
			stack.Path = ""
		}

		if err.Cause != nil {
			causeStack := BuildStack(err.Cause, opts)
			if len(causeStack) > 0 {
				stack.Stack = causeStack
			}
			stack.Cause = nil
		}
	case error:
		stack.Error = Error{
			Type:    Unknown,
			Message: err.Error(),
		}
	}

	return []Stack{stack}
}
