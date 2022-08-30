package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func WrapHandlerWithPersistence(handler interface{}) interface{} {
	return func(ctx context.Context, msg json.RawMessage) (interface{}, error) {
		return callHandler(ctx, msg, handler)
	}
}

func callHandler(ctx context.Context, msg json.RawMessage, handler interface{}) (response interface{}, errorResponse error) {
	event, err := unmarshalEventForHandler(msg, handler)
	if err != nil {
		return nil, err
	}

	handlerType := reflect.TypeOf(handler)
	arguments := []reflect.Value{}

	if handlerType.NumIn() == 1 {
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		firstArgType := handlerType.In(0)
		if firstArgType.Implements(contextType) {
			arguments = []reflect.Value{reflect.ValueOf(ctx)}
		} else {
			arguments = []reflect.Value{event.Elem()}
		}
	} else if handlerType.NumIn() == 2 {
		arguments = []reflect.Value{reflect.ValueOf(ctx), event.Elem()}
	}

	handlerValue := reflect.ValueOf(handler)
	output := handlerValue.Call(arguments)

	if len(output) > 0 {
		val := output[len(output)-1].Interface()
		if errVal, ok := val.(error); ok {
			errorResponse = errVal
		}
	}

	if len(output) > 1 {
		response = output[0].Interface()
	}

	return
}

func unmarshalEventForHandler(msg json.RawMessage, handler interface{}) (reflect.Value, error) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() == 0 {
		return reflect.ValueOf(nil), nil
	}

	messageType := handlerType.In(handlerType.NumIn() - 1)
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	firstArgType := handlerType.In(0)

	if handlerType.NumIn() == 1 && firstArgType.Implements(contextType) {
		return reflect.ValueOf(nil), nil
	}

	newMessage := reflect.New(messageType)
	err := json.Unmarshal(msg, newMessage.Interface())
	if err != nil {
		return reflect.ValueOf(nil), err
	}

	fmt.Printf("%#v\n", newMessage.Elem().FieldByName("Records"))

	if newMessage.Elem().FieldByName("Records").Len() > 0 {
		for idx := 0; idx < newMessage.Elem().FieldByName("Records").Len(); idx++ {
			fmt.Printf("%d - %#v\n", idx, newMessage.Elem().FieldByName("Records").Index(idx).FieldByName("MessageId"))

			ddi, err := getAttributeFromMsg(newMessage.Elem().FieldByName("Records").Index(idx), "MessageDeduplicationId")
			if err != nil {
				fmt.Printf("error getting attribute: %s\n", err.Error())
				continue
			}

			fmt.Printf("%d - %#v\n", idx, ddi)
		}
	}

	return newMessage, err
}

func getAttributeFromMsg(msg reflect.Value, key string) (value string, err error) {
	if msg.CanInterface() {
		if msg.FieldByName("Attributes").Kind() == reflect.Map {
			value = msg.FieldByName("Attributes").MapIndex(reflect.ValueOf(key)).String()
		} else {
			err = errors.New("msg does not contain Attributes map")
		}
	} else {
		err = errors.New("msg is not addressable")
	}

	return
}
