package goservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	context := MakeContext()

	organizer := MakeOrganizer(
		addsNumberToContext,
		addsOneToNumber)
	result := organizer.Call(context)

	assert.Equal(t, 2, result["number"].(int))
}
