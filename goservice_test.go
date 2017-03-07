package goservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func convertsMessageToUpperCase(ctx *Context) *Context {
	ctx.Message = strings.ToUpper(ctx.Message)
	return ctx
}

func addsACharacter(ctx *Context) *Context {
	ctx.Message = fmt.Sprintf("%sa", ctx.Message)
	return ctx
}

func failsContext(ctx *Context) *Context {
	ctx.Fail("I don't like this")
	return ctx
}

func addsNumberToContext(ctx *Context) *Context {
	ctx.Set("number", 1)
	return ctx
}

func addsOneToNumber(ctx *Context) *Context {
	number := ctx.Get("number").(int)
	ctx.Set("number", number+1)
	return ctx
}

func handlesError(ctx *Context) *Context {
	_, err := os.Open("nofile.txt")
	if err != nil {
		ctx.Error = err
	}

	return ctx
}

func Test_AlteringMessage(t *testing.T) {
	context := NewContext()
	context.Message = "message"

	organizer := NewOrganizer(
		convertsMessageToUpperCase,
		addsACharacter)

	result := organizer.Call(context)

	assert.Equal(t, "MESSAGEa", result.Message)
}

func Test_FailContext(t *testing.T) {
	context := NewContext()

	organizer := NewOrganizer(
		failsContext,
		addsACharacter)
	result := organizer.Call(context)

	assert.False(t, result.Success)
	assert.Equal(t, "I don't like this", context.Message)
}

func Test_AddItemToContext(t *testing.T) {
	organizer := NewOrganizer(
		addsNumberToContext,
		addsOneToNumber)
	result := organizer.Call(NewContext())

	assert.Equal(t, 2, result.Get("number").(int))
}

func Test_CapturesError(t *testing.T) {
	organizer := NewOrganizer(
		handlesError,
		addsNumberToContext)
	result := organizer.Call(NewContext())

	assert.NotNil(t, result.Error)
}

// Call it with:
// $: go test -v -run="none" -benchtime="3s" -bench="BenchmarkOrganizer" -benchmem
func BenchmarkOrganizer(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		organizer := NewOrganizer(
			addsNumberToContext,
			addsOneToNumber)
		organizer.Call(NewContext())
	}
}
