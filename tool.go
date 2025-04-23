package agent

type Tool struct {
	Name        string
	Description string
	Parameters  []ToolParameters
	Handler     func(map[string]any) (any, error)
}

type ToolParameters struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

type ToolCall struct {
	CallId    string
	ToolName  string
	Arguments map[string]any
}

type ToolCallResult struct {
	CallId   string
	ToolName string
	Result   any
	Error    error
}
