package goservice

type Context map[string]interface{}

func (ctx Context) Message() string {
	return ctx["__message"].(string)
}

func (ctx Context) SetMessage(msg string) {
	ctx["__message"] = msg
}

func (ctx Context) SetSuccess() {
	ctx["__success"] = true
}

func (ctx Context) SetFailure() {
	ctx["__success"] = false
}

func (ctx Context) IsSuccess() bool {
	return ctx["__success"].(bool)
}

func (ctx Context) IsFailure() bool {
	return ctx["__success"].(bool) == false
	// return ctx.Success == false
}

func MakeContext() Context {
	ctx := Context{}
	ctx.SetSuccess()
	return ctx
}

type Organizer struct {
	Actions []func(Context) Context
	Ctx     Context
}

func MakeOrganizer(actions ...func(Context) Context) Organizer {
	organizer := Organizer{Actions: actions}
	return organizer
}

func (organizer Organizer) Call(ctx Context) Context {
	return organizer.With(ctx).Reduce(organizer.Actions)
}

func (organizer Organizer) With(ctx Context) Organizer {
	organizer.Ctx = ctx
	return organizer
}

func (organizer Organizer) Reduce(actions []func(Context) Context) Context {
	ctx := organizer.Ctx
	for _, action := range actions {
		if ctx.IsSuccess() {
			ActionHandler(action).Execute(ctx)
		}
	}
	return ctx
}

// The ActionHandler type is an adapter to allow the use of
// ordinary functions as Action handlers. If f is a function
// with the appropriate signature, ActionHandler(f) is a
// Handler that calls f.
type ActionHandler func(ctx Context) Context

func (f ActionHandler) Execute(ctx Context) Context {
	return f(ctx)
}
