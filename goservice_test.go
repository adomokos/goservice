package goservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func convertsMessageToUpperCase(ctx Context) Context {
	ctx.SetMessage(strings.ToUpper(ctx.Message()))
	return ctx
}

func addsACharacter(ctx Context) Context {
	ctx.SetMessage(fmt.Sprintf("%sa", ctx.Message()))
	return ctx
}

func failsContext(ctx Context) Context {
	ctx.SetFailure()
	return ctx
}

func addsNumberToContext(ctx Context) Context {
	ctx["number"] = 1
	return ctx
}

func addsOneToNumber(ctx Context) Context {
	number := ctx["number"].(int)
	ctx["number"] = number + 1
	return ctx
}

func handlesError(ctx Context) Context {
	_, err := os.Open("nofile.txt")
	if err != nil {
		ctx.SetError(err)
	}

	return ctx
}

func Test_AlteringMessage(t *testing.T) {
	context := MakeContext()
	context.SetMessage("message")

	organizer := MakeOrganizer(
		convertsMessageToUpperCase,
		addsACharacter)

	result := organizer.Call(context)

	assert.Equal(t, "MESSAGEa", result.Message())
}

func Test_FailContext(t *testing.T) {
	context := MakeContext()
	context.SetMessage("message")

	organizer := MakeOrganizer(
		failsContext,
		addsACharacter)
	result := organizer.Call(context)

	assert.Equal(t, "message", result.Message())
	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsFailure())
}

func Test_AddItemToContext(t *testing.T) {
	organizer := MakeOrganizer(
		addsNumberToContext,
		addsOneToNumber)
	result := organizer.Call(MakeContext())

	assert.Equal(t, 2, result["number"].(int))
}

func Test_CapturesError(t *testing.T) {
	organizer := MakeOrganizer(
		handlesError,
		addsNumberToContext)
	result := organizer.Call(MakeContext())

	assert.False(t, result.IsSuccess())
}

// Call it with:
// $: go test -v -run="none" -benchtime="3s" -bench="BenchmarkOrganizer" -benchmem
func BenchmarkOrganizer(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		organizer := MakeOrganizer(
			addsNumberToContext,
			addsOneToNumber)
		organizer.Call(MakeContext())
	}
}
