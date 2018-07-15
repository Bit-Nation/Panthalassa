package callbacks

import (
	"testing"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestFuncRegistration(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	require.Equal(t, uint(0), m.reqLim.Count())
	funcId, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	require.Nil(t, err)
	require.Equal(t, uint(1), m.reqLim.Count())

	id, err := funcId.ToInteger()
	require.Nil(t, err)
	require.Equal(t, int64(1), id)

	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

}

func TestFuncUnRegisterSuccess(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	// register function

	funcId, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	require.Nil(t, err)

	// make sure the id is the one we expect it to be
	id, err := funcId.ToInteger()
	require.Nil(t, err)
	require.Equal(t, int64(1), id)

	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	// un-register function that has never been registered
	require.Equal(t, uint(1), m.reqLim.Count())
	returnedValue, err := vm.Call(`unRegisterFunction`, vm, 23565)
	require.Nil(t, err)
	returnedValueStr, err := returnedValue.ToString()
	require.Nil(t, err)
	require.Equal(t, "undefined", returnedValueStr)
	// request limitation must still be one since function that have not been registerd
	// can't decrease the counter
	require.Equal(t, uint(1), m.reqLim.Count())

	// successfully un-register function
	returnedValue, err = vm.Call(`unRegisterFunction`, vm, id)
	require.Nil(t, err)
	returnedValueStr, err = returnedValue.ToString()
	require.Nil(t, err)
	require.Equal(t, "undefined", returnedValueStr)
	// should be 0 since the the function id we passed in exists
	require.Equal(t, uint(0), m.reqLim.Count())

}

func TestFuncCallSuccess(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		_, err = cb.Call(cb)
		if err != nil {
			m.logger.Error(err.Error())
		}

		return otto.Value{}
	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	require.Nil(t, m.CallFunction(1, `{key: "value"}`))

}

func TestFuncCallError(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		_, err = cb.Call(cb, "I am an error")

		if err != nil {
			m.logger.Error(err.Error())
		}

		return otto.Value{}
	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	require.Equal(t, "I am an error", m.CallFunction(1, `{key: "value"}`).Error())

}

func TestFuncCallBackTwice(t *testing.T) {

	m := New(log.MustGetLogger(""))

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		val, err := cb.Call(cb)
		if err != nil {
			panic(err)
		}
		if !val.IsUndefined() {
			panic("expected value to be undefined")
		}

		val, err = cb.Call(cb)
		if err != nil {
			panic(err)
		}
		if val.String() != "Callback: Already called callback" {
			panic("Expected an error that tells me that I alrady called the callback")
		}

		return otto.Value{}

	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	m.CallFunction(1, `{key: "value"}`)

}
