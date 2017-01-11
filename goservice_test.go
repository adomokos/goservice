package goservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type ConvertsMessageToUpperCaseAction struct{}

func (action ConvertsMessageToUpperCaseAction) Execute(ctx Context) Context {
	ctx.SetMessage(strings.ToUpper(ctx.Message()))

	return ctx
}

type AddsACharacterAction struct{}

func (action AddsACharacterAction) Execute(ctx Context) Context {
	ctx.SetMessage(fmt.Sprintf("%sa", ctx.Message()))

	return ctx
}

type FailsAction struct{}

func (action FailsAction) Execute(ctx Context) Context {
	ctx.SetFailure()

	return ctx
}

func Test_Action_Call(t *testing.T) {
	context := MakeContext()
	context.SetMessage("message")

	organizer := MakeOrganizer(
		ConvertsMessageToUpperCaseAction{},
		AddsACharacterAction{})

	result := organizer.Call(context)

	assert.Equal(t, "MESSAGEa", result.Message())
}

func Test_FailAction(t *testing.T) {
	context := MakeContext()
	context.SetMessage("message")

	organizer := MakeOrganizer(
		FailsAction{},
		AddsACharacterAction{})
	result := organizer.Call(context)

	assert.Equal(t, "message", result.Message())
	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsFailure())
}
