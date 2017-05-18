package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d : %s", e.Code, e.Message)
}

// Bind parses request body into obj structure
func Bind(req *http.Request, obj interface{}) error {
	var err error

	switch getContentType(req) {
	case "application/json":
		err = BindJSON(req, obj)
	case "application/xml":
		err = BindXML(req, obj)
	}
	if err != nil {
		return err
	}

	if err = parseRequestBinding(req, obj); err != nil {
		return err
	}
	return nil
}

// BindSkipBody parses request header and URI into obj structure, skipping request body
func BindSkipBody(req *http.Request, obj interface{}) error {
	if err := parseRequestBinding(req, obj); err != nil {
		return err
	}
	return nil
}

// BindJSON parses JSON request body
func BindJSON(req *http.Request, obj interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
		return &Error{Code: 400, Message: fmt.Sprintf("malformed json: %s", err)}
	}
	return nil
}

// BindXML parses XML request body
func BindXML(req *http.Request, obj interface{}) error {
	if err := xml.NewDecoder(req.Body).Decode(obj); err != nil {
		return &Error{Code: 400, Message: "MalformedXML"}
	}
	return nil
}

func getContentType(req *http.Request) string {
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		return req.Header.Get("Content-Type")
	}
	return ""
}

func parseRequestBinding(req *http.Request, obj interface{}) error {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			err := parseRequestBinding(req, val.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		name := field.Tag.Get("header")

		if name != "" {
			value := req.Header.Get(name)
			if value != "" {
				err := parseVar(typ, &val, i, value)
				if err != nil {
					return err
				}
			}
		}

		name = field.Tag.Get("mux")

		if name != "" {
			value := mux.Vars(req)[name]
			if value != "" {
				err := parseVar(typ, &val, i, value)
				if err != nil {
					return err
				}
			}
		}

		name = field.Tag.Get("query")

		if name != "" {
			value := req.URL.Query().Get(name)
			if value != "" {
				err := parseVar(typ, &val, i, value)
				if err != nil {
					return err
				}
			}
		}

		name = typ.Field(i).Tag.Get("preprocess")

		if name == "urldecode" {
			value, err := url.QueryUnescape(val.Field(i).String())
			if err != nil {
				return err
			}
			val.Field(i).SetString(value)
		} else if name != "" {
			panic(fmt.Sprintf("Unknown preprocess type: %s", name))
		}

	}

	return nil
}

func parseVar(typ reflect.Type, val *reflect.Value, index int, value string) error {
	field := typ.Field(index)

	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || val.Field(index).OverflowInt(n) {
			return &Error{Code: 400, Message: "wrong value for field " + field.Name}
		}

		val.Field(index).SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil || val.Field(index).OverflowUint(n) {
			return &Error{Code: 400, Message: "wrong value for field " + field.Name}
		}
		val.Field(index).SetUint(n)
		return nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(value, val.Field(index).Type().Bits())
		if err != nil || val.Field(index).OverflowFloat(n) {
			return &Error{Code: 400, Message: "wrong value for field " + field.Name}
		}
		val.Field(index).SetFloat(n)
		return nil
	case reflect.Bool:
		n, err := strconv.ParseBool(value)
		if err != nil {
			return &Error{Code: 400, Message: "wrong value for field " + field.Name}
		}
		val.Field(index).SetBool(n)
		return nil
	case reflect.String:
		val.Field(index).SetString(value)
		return nil
	default:
		if field.Type == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return &Error{Code: 400, Message: "wrong value for field " + field.Name + ": " + err.Error()}
			}
			val.Field(index).Set(reflect.ValueOf(t))
			return nil
		}
	}

	panic("unknown binding type")
}
