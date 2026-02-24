package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/kataras/iris/v12"
)

var mutexSpan sync.RWMutex
var mutexMap sync.RWMutex

// GetCurrentFunctionName 2 current function - 3 previous function
func GetCurrentFunctionName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip, pc)
	f := runtime.FuncForPC(pc[0])
	values := strings.Split(f.Name(), ".")
	return values[len(values)-1]
}

func GetCallers(skip int) []string {
	pc := make([]uintptr, 40)
	runtime.Callers(skip, pc)
	var values []string
	for _, u := range pc {
		f := runtime.FuncForPC(u)
		name := f.Name()
		if name != "" {
			v := strings.Split(name, ".")
			values = append(values, v[len(v)-1])
		}
	}
	return values
}

// CheckIfValueInConstant check if the value exists on the constant value.
// Example: constant = CL,CO; Value = CO; Result = true
func CheckIfValueInConstant(value string, constant string) bool {
	isValid := false
	splitString := strings.Split(constant, ",")
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(splitString))
	for i, v := range splitString {
		go func(i int, v string) {
			defer waitGroup.Done()
			if strings.ReplaceAll(v, " ", "") == value {
				isValid = true
			}
		}(i, v)
	}
	waitGroup.Wait()
	return isValid
}

// GetObjectFromContext retrieve the object presents on the context request
func GetObjectFromContext(ctx iris.Context, v any) error {
	body, errRead := ctx.GetBody()
	if errRead != nil {
		return errRead
	}
	errUnmarshal := json.Unmarshal(body, &v)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	return nil
}

func UnmarshalToObject(toDecode []byte, v any) error {
	if json.Valid(toDecode) {
		bufferOfBytes := &bytes.Buffer{}
		json.HTMLEscape(bufferOfBytes, toDecode)
		if errUnmarshal := json.Unmarshal(bufferOfBytes.Bytes(), v); errUnmarshal != nil {
			return errUnmarshal
		}
		return nil
	} else {
		return fmt.Errorf("json data is not valid")
	}
}

func ParseTo(source any, destiny any) error {
	marshal, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = UnmarshalToObject(marshal, &destiny)
	if err != nil {
		return err
	}
	return nil
}
