package resolvers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type resolver struct {
	function interface{}
}

func (r *resolver) hasArguments() bool {
	return reflect.TypeOf(r.function).NumIn() == 1
}

func (r *resolver) hasArgumentsAndIdentity() bool {
	return reflect.TypeOf(r.function).NumIn() == 2
}

func (r *resolver) call(p json.RawMessage, i json.RawMessage) (interface{}, error) {
	var args []reflect.Value
	var identityArg reflect.Value
	var err error

	if r.hasArguments() {
		pld := payload{p}
		args, err = pld.parse(reflect.TypeOf(r.function).In(0))

		if err != nil {
			return nil, err
		}
	} else if r.hasArgumentsAndIdentity() {
		pld := payload{p}
		args, err = pld.parse(reflect.TypeOf(r.function).In(0))

		identity := reflect.New(reflect.TypeOf(r.function).In(1))
		if err := json.Unmarshal(i, identity.Interface()); err != nil {
			return nil, fmt.Errorf("Unable to prepare payload: %s", err.Error())
		}
		identityArg = identity.Elem()

		if err != nil {
			return nil, err
		}
	}

	returnValues := reflect.ValueOf(r.function).Call(append(args, identityArg))
	var returnData interface{}
	var returnError error

	if len(returnValues) == 2 {
		returnData = returnValues[0].Interface()
	}

	if err := returnValues[len(returnValues)-1].Interface(); err != nil {
		returnError = returnValues[len(returnValues)-1].Interface().(error)
	}

	return returnData, returnError
}
