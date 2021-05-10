package chttp

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gocopper/copper/cerrors"
)

type htmlComponentReflector struct {
	component HTMLComponent
}

func newHTMLComponentReflector(c HTMLComponent) (*htmlComponentReflector, error) {
	return &htmlComponentReflector{component: c}, nil // todo: validate component
}

func (c *htmlComponentReflector) getRenderMethod() reflect.Value {
	return reflect.ValueOf(c.component).MethodByName("Render")
}

func (c *htmlComponentReflector) getActionMethod(name string) (reflect.Value, error) {
	return reflect.ValueOf(c.component).MethodByName(name), nil // todo: validate method
}

func (c *htmlComponentReflector) createEmptyPropsStruct() reflect.Value {
	return reflect.New(c.getRenderMethod().Type().In(1)).Elem()
}

func (c *htmlComponentReflector) createActionMethodArgs(action string, argValues []json.RawMessage) ([]reflect.Value, error) {
	const MinNumActionArgs = 2

	args := make([]reflect.Value, 0)

	method, err := c.getActionMethod(action)
	if err != nil {
		return nil, cerrors.New(err, "failed to get action method", map[string]interface{}{
			"component": c.component.Name(),
			"action":    action,
		})
	}

	if method.Type().NumIn() <= MinNumActionArgs {
		return args, nil
	}

	for i := range argValues {
		field := reflect.New(method.Type().In(MinNumActionArgs + i)).Interface()

		err := json.Unmarshal(argValues[i], field)
		if err != nil {
			return nil, cerrors.New(err, "failed to unmarshal arg value", map[string]interface{}{
				"component": c.component.Name(),
				"action":    action,
				"arg":       method.Type().In(i).Name(),
			})
		}

		args = append(args, reflect.ValueOf(field).Elem())
	}

	return args, nil
}

func (c *htmlComponentReflector) createPropsStructFromInterfaceValues(propValues []interface{}) reflect.Value {
	props := c.createEmptyPropsStruct()

	for i := range propValues {
		props.Field(i).Set(reflect.ValueOf(propValues[i]).Convert(props.Field(i).Type()))
	}

	return props
}

func (c *htmlComponentReflector) createPropsInterfaceValuesFromJSONValues(propValues []json.RawMessage) ([]interface{}, error) {
	props, err := c.createPropsStructFromJSONValues(propValues)
	if err != nil {
		return nil, cerrors.New(err, "failed to create props struct from json", nil)
	}

	interfaceValues := make([]interface{}, props.NumField())

	for i := 0; i < props.NumField(); i++ {
		interfaceValues[i] = props.Field(i).Interface()
	}

	return interfaceValues, nil
}

func (c *htmlComponentReflector) createPropsStructFromJSONValues(propValues []json.RawMessage) (reflect.Value, error) {
	props := c.createEmptyPropsStruct()

	for i := range propValues {
		field := reflect.New(props.Field(i).Type()).Interface()

		err := json.Unmarshal(propValues[i], field)
		if err != nil {
			return reflect.ValueOf(nil), cerrors.New(err, "failed to unmarshal props", map[string]interface{}{
				"component": c.component.Name(),
				"field":     props.Field(i).Type().Name(),
			})
		}

		props.Field(i).Set(reflect.ValueOf(field).Elem())
	}

	return props, nil
}

func (c *htmlComponentReflector) callRenderMethod(r *http.Request, props reflect.Value) (interface{}, error) {
	var (
		method = c.getRenderMethod()
		args   = []reflect.Value{
			reflect.ValueOf(r),
			props,
		}
	)

	result := method.Call(args)

	if !result[1].IsNil() {
		return nil, result[1].Interface().(error)
	}

	return result[0].Interface(), nil
}

func (c *htmlComponentReflector) callActionMethod(r *http.Request, action string, props reflect.Value, args []reflect.Value) ([]interface{}, error) {
	method, err := c.getActionMethod(action)
	if err != nil {
		return nil, cerrors.New(err, "failed to get action method", map[string]interface{}{
			"component": c.component,
			"action":    action,
		})
	}

	in := []reflect.Value{
		reflect.ValueOf(r),
		props,
	}

	in = append(in, args...)

	out := method.Call(in)

	if !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}

	newProps := make([]interface{}, 0)
	for i := 0; i < out[0].NumField(); i++ {
		newProps = append(newProps, out[0].Field(i).Interface())
	}

	return newProps, nil
}
