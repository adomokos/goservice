package goservice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func convertsMessageToUpperCase(ctx *Context) {
	ctx.Message = strings.ToUpper(ctx.Message)
}

func addsACharacter(ctx *Context) {
	ctx.Message = fmt.Sprintf("%sa", ctx.Message)
}

func failsContext(ctx *Context) {
	ctx.Fail("I don't like this")
}

func addsNumberToContext(ctx *Context) {
	ctx.Set("number", 1)
}

func addsOneToNumber(ctx *Context) {
	number := ctx.Get("number").(int)
	ctx.Set("number", number+1)
}

func handlesError(ctx *Context) {
	_, err := os.Open("nofile.txt")
	if err != nil {
		ctx.Error = err
	}
}

func Test_AlteringMessage(t *testing.T) {
	context := NewContext()
	context.Message = "message"

	organizer := NewOrganizer(
		convertsMessageToUpperCase,
		addsACharacter)

	organizer.Call(context)

	assert.Equal(t, "MESSAGEa", context.Message)
}

// func Test_FailContext(t *testing.T) {
// context := NewContext()

// organizer := NewOrganizer(
// failsContext,
// addsACharacter)
// result := organizer.Call(context)

// assert.False(t, result.Success)
// assert.Equal(t, "I don't like this", context.Message)
// }

// func Test_AddItemToContext(t *testing.T) {
// organizer := NewOrganizer(
// addsNumberToContext,
// addsOneToNumber)
// result := organizer.Call(NewContext())

// assert.Equal(t, 2, result.Get("number").(int))
// }

// func Test_CapturesError(t *testing.T) {
// organizer := NewOrganizer(
// handlesError,
// addsNumberToContext)
// result := organizer.Call(NewContext())

// assert.NotNil(t, result.Error)
// }

// Call it with:
// $: go test -v -run="none" -benchtime="3s" -bench="BenchmarkOrganizer" -benchmem
// func BenchmarkOrganizer(b *testing.B) {
// b.ResetTimer()

// for i := 0; i < b.N; i++ {
// organizer := NewOrganizer(
// addsNumberToContext,
// addsOneToNumber)
// organizer.Call(NewContext())
// }
// }
