package functions

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/kataras/iris/v12/context"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentFunctionName(t *testing.T) {
	testTable := []struct {
		name     string
		skip     int
		expected string
	}{
		{
			name:     "return func1",
			skip:     2,
			expected: "func1",
		},
		{
			name:     "return tRunner",
			skip:     3,
			expected: "tRunner",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actual := GetCurrentFunctionName(testCase.skip)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestCheckIfValueInConstant(t *testing.T) {
	testTable := []struct {
		name     string
		value    string
		constant string
		expected bool
	}{
		{
			name:     "return true",
			value:    "CO",
			constant: "CL,CO",
			expected: true,
		},
		{
			name:     "return false",
			value:    "CO",
			constant: "CL",
			expected: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actual := CheckIfValueInConstant(testCase.value, testCase.constant)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

type MockErrorReader struct {
	body string
}

func (m MockErrorReader) Read(p []byte) (n int, err error) {
	if m.body != "" {
		copied := copy(p, m.body)
		return copied, io.EOF
	}
	return 0, io.ErrUnexpectedEOF
}
func (m MockErrorReader) Close() error {
	return nil
}

func TestGetObjectFromContext(t *testing.T) {
	testTable := []struct {
		name           string
		body           string
		v              any
		expectedError  error
		expectedObject any
	}{
		{
			name:           "return body error",
			body:           "",
			v:              nil,
			expectedError:  io.ErrUnexpectedEOF,
			expectedObject: nil,
		},
		{
			name:           "return unmarshal error",
			body:           "invalid json",
			v:              nil,
			expectedError:  errors.New("invalid character 'i' looking for beginning of value"),
			expectedObject: nil,
		},
		{
			name: "return object",
			body: "{\"key\":\"value\"}",
			v: map[string]string{
				"key": "value",
			},
			expectedError: nil,
			expectedObject: map[string]string{
				"key": "value",
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			mockErrorReader := MockErrorReader{
				body: tt.body,
			}

			app := iris.New()
			ctx := context.NewContext(app)
			req := &http.Request{}
			req.Body = mockErrorReader
			ctx.ResetRequest(req)

			err := GetObjectFromContext(ctx, tt.v)
			if tt.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}

			assert.Equal(t, tt.expectedObject, tt.v)
		})
	}
}

func TestUnmarshalToObject(t *testing.T) {
	testTable := []struct {
		name           string
		toDecode       []byte
		v              any
		expectedError  error
		expectedObject any
	}{
		{
			name:          "return error json data is not valid",
			toDecode:      []byte("invalid json"),
			v:             nil,
			expectedError: errors.New("json data is not valid"),
		},
		{
			name:          "return error unmarshal",
			toDecode:      []byte("{\"key\":\"value\"}"),
			v:             map[string]string{},
			expectedError: errors.New("json: Unmarshal(non-pointer map[string]string)"),
		},
		{
			name:     "return object",
			toDecode: []byte(`{"key":"value"}`),
			v: &map[string]interface{}{
				"key": "value",
			},
			expectedError: nil,
			expectedObject: &map[string]interface{}{
				"key": "value",
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalToObject(tt.toDecode, tt.v)
			if tt.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestParseTo(t *testing.T) {
	testTable := []struct {
		name          string
		source        any
		destiny       any
		errorExpected error
	}{
		{
			name:          "return error marshal",
			source:        func() {},
			destiny:       nil,
			errorExpected: errors.New("json: unsupported type: func()"),
		},
		{
			name:          "return error json data is not valid",
			source:        "12345",
			destiny:       &map[string]string{},
			errorExpected: errors.New("json: cannot unmarshal string into Go value of type map[string]string"),
		},
		{
			name:          "return struct",
			source:        &map[string]string{"key": "value"},
			destiny:       &map[string]string{},
			errorExpected: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseTo(tt.source, tt.destiny)
			if tt.errorExpected == nil {
				assert.Nil(t, err)
				assert.Equal(t, tt.source, tt.destiny)
			} else {
				assert.Equal(t, tt.errorExpected.Error(), err.Error())
			}
		})
	}
}

func TestGenerateHttpResponse(t *testing.T) {
	testTable := []struct {
		name                   string
		status                 int
		body                   any
		modelsResponseExpected *rest.Response
	}{
		{
			name:   "return response with string body",
			status: 200,
			body:   "body",
			modelsResponseExpected: &rest.Response{
				Http: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
					Header:     make(http.Header),
				},
			},
		},
		{
			name:   "return response with json body",
			status: 200,
			body: map[string]string{
				"key": "value",
			},
			modelsResponseExpected: &rest.Response{
				Http: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte("{\"key\":\"value\"}"))),
					Header:     make(http.Header),
				},
			},
		},
		{
			name:   "return response with error body",
			status: 500,
			body:   func() {},
			modelsResponseExpected: &rest.Response{
				Http: &http.Response{
					Status:     http.StatusText(500),
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewReader([]byte("response can't be parsed"))),
					Header:     make(http.Header),
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			modelsResponse := GenerateHttpResponse(tt.status, tt.body)
			assert.Equal(t, tt.modelsResponseExpected, modelsResponse)
		})
	}
}
