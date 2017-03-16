package goservice

import (
	"sync"
)

var (
	mutex sync.RWMutex
)

type Context struct {
	data    map[string]interface{}
	Success bool
	Message string
	Error   error
}

func NewContext() *Context {
	return &Context{
		data:    make(map[string]interface{}),
		Success: true,
		Message: ""}
}

func (ctx *Context) Set(key string, val interface{}) {
	mutex.Lock()
	ctx.data[key] = val
	mutex.Unlock()
}

func (ctx *Context) Get(key string) interface{} {
	mutex.RLock()
	if val := ctx.data[key]; val != nil {
		value := ctx.data[key]
		mutex.RUnlock()
		return value
	}
	mutex.RUnlock()
	return nil
}

func (ctx *Context) Fail(message string) {
	mutex.Lock()
	ctx.Message = message
	ctx.Success = false
	mutex.Unlock()
}

func (ctx *Context) Delete(key string) {
	mutex.Lock()
	if ctx.data[key] != nil {
		delete(ctx.data, key)
	}
	mutex.Unlock()
}

func (ctx *Context) Clear() {
	mutex.Lock()
	ctx.data = make(map[string]interface{})
	mutex.Unlock()
}

type Organizer struct {
	Actions []func(*Context)
	Ctx     *Context
}

func NewOrganizer(actions ...func(*Context)) Organizer {
	organizer := Organizer{Actions: actions}
	return organizer
}

func (organizer Organizer) Call(ctx *Context) {
	organizer.With(ctx).Reduce(organizer.Actions)
}

func (organizer Organizer) With(ctx *Context) Organizer {
	organizer.Ctx = ctx
	return organizer
}

func (organizer Organizer) Reduce(actions []func(*Context)) {
	ctx := organizer.Ctx
	for _, action := range actions {
		if ctx.Success {
			ActionHandler(action).Execute(ctx)
		}
	}
}

// The ActionHandler type is an adapter to allow the use of
// ordinary functions as Action handlers. If f is a function
// with the appropriate signature, ActionHandler(f) is a
// Handler that calls f.
type ActionHandler func(ctx *Context)

func (f ActionHandler) Execute(ctx *Context) {
	f(ctx)
}
