package context

type context struct {
	Name string
	Data map[string]interface{}
}

// Add singleton contexts identified by name
var (
	ctxs map[string]*context
)

// Add context option 'WithData'
type ContextOption func(*context)

func WithData(data map[string]interface{}) ContextOption {
	return func(ctx *context) {
		ctx.Data = data
	}
}

// Context returns the singleton context
func Context(name string, opts ...ContextOption) *context {
	if ctxs == nil {
		ctxs = make(map[string]*context)
	}
	if ctx, ok := ctxs[name]; ok {
		return ctx
	}
	ctx := &context{
		Name: name,
		Data: make(map[string]interface{}),
	}
	// Loop through each option
	for _, opt := range opts {
		opt(ctx)
	}
	ctxs[name] = ctx
	return ctx
}

func Exists(name string) bool {
	if ctxs == nil {
		return false
	}
	_, ok := ctxs[name]
	return ok
}

// SetContext sets the singleton context
func (ctx *context) SetData(data map[string]interface{}) {
	ctx.Data = data
}

// Get returns the value for a key
func (ctx *context) Get(key string) (interface{}, bool) {
	value, ok := ctx.Data[key]
	return value, ok
}

// Set sets the value for a key
func (ctx *context) Set(key string, value interface{}) {
	ctx.Data[key] = value
}
