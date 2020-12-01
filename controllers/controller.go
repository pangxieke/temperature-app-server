//    type MyController struct {
//      controller.Base
//    }
//
//    func (c *MyController) Index() error {
//		name := c.Params.ByName("name")
//      c.ResponseWriter.Write([]byte("URL Param name set to: " + name))
//      return nil
//    }
//
// To handle HTTP requests with this controller, use the controller.Action
// function:
//    router := httprouter.New()
//    router.GET("/hello/:name", controller.Action((*MyController).Index))
//
package controllers

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Controller interface {
	Init(http.ResponseWriter, *http.Request, httprouter.Params) error
	Destroy()
	Error(err error)
}

type Base struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Params         httprouter.Params
	RequestBody    string
	UUID           string
}

func (b *Base) Init(rw http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	b.Request, b.ResponseWriter, b.Params = r, rw, params
	if b.Request.ContentLength > 0 {
		buffer := make([]byte, b.Request.ContentLength)
		b.Request.Body.Read(buffer)
		b.RequestBody = string(buffer)
	}
	id, _ := uuid.NewUUID()
	b.UUID = fmt.Sprintf("%s", id)
	return nil
}

func (b *Base) Destroy() {
}

func (b *Base) Error(err error) {
	logs.Info(err.Error())
	http.Error(b.ResponseWriter, err.Error(), http.StatusInternalServerError)
}

// Action takes a method expression and translates it into a callable
// httprouter.Handle which, when called:
//
// 		1. Constructs a controller instance
// 		2. Initializes the controller via the Init function
// 		3. Invokes the Action method referenced by the method expression
// 		4. Calls destroy on the controller
//
// This flow allows for similar logic to be cleanly reused while data is no
// longer shared between requests. This is because a new Controller instance
// will be constructed every time the returned httprouter.Handle's ServeHTTP method
// is invoked.
//
// An example of a valid method expression is:
//
// 		controller.Action((*MyController).Index)
//
// Where MyController is an implementor of the Controller interface and Index
// is a method on MyController that takes no arguments and returns an err
func Action(action interface{}) httprouter.Handle {
	val := reflect.ValueOf(action)
	t, err := controllerType(val)
	if err != nil {
		panic(err)
	}

	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		v := reflect.New(t)
		c := v.Interface().(Controller)
		err = c.Init(rw, r, params)
		defer c.Destroy()
		if err != nil {
			c.Error(err)
			return
		}
		ret := val.Call([]reflect.Value{v})[0].Interface()
		if ret != nil {
			c.Error(ret.(error))
			return
		}
	}
}

func controllerType(action reflect.Value) (reflect.Type, error) {
	t := action.Type()

	if t.Kind() != reflect.Func {
		return t, errors.New("Action is not a function")
	}

	if t.NumIn() != 1 {
		return t, errors.New("Wrong Number of Arguments in action")
	}

	if t.NumOut() != 1 {
		return t, errors.New("Wrong Number of return values in action")
	}

	out := t.Out(0)
	if !out.Implements(interfaceOf((*error)(nil))) {
		return t, errors.New("Action return type invalid")
	}

	t = t.In(0)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if !reflect.PtrTo(t).Implements(interfaceOf((*Controller)(nil))) {
		return t, errors.New("Controller does not implement ctrl.Controller interface")
	}

	return t, nil
}

func interfaceOf(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}
