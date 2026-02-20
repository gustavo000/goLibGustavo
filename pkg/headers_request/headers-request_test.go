package headers_request

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	constants_headers "github.com/gustavo000/goLibGustavo/constants/constants_headers"
	"github.com/gustavo000/goLibGustavo/resources/properties"
)

func TestRequestHeaders_ModifyHeaders(t *testing.T) {
	requestHeaders := RequestHeaders{}
	requestHeaders = append(requestHeaders, RequestHeader{Key: "test", DefaultValue: "test"})

	rq := requestHeaders.ModifyHeaders(RequestHeader{Key: "test", DefaultValue: "test"})
	assert.Len(t, rq, 1)
	assert.Equal(t, RequestHeader{Key: "test", DefaultValue: "test"}, rq[0])
}

func TestRequestHeaders_AddHeaders(t *testing.T) {
	requestHeaders := RequestHeaders{}
	requestHeaders = append(requestHeaders, RequestHeader{Key: "test", DefaultValue: "test"})

	rq := requestHeaders.AddHeaders(RequestHeader{Key: "test2", DefaultValue: "test2"})
	assert.Len(t, rq, 2)
	assert.Equal(t, RequestHeader{Key: "test", DefaultValue: "test"}, rq[0])
	assert.Equal(t, RequestHeader{Key: "test2", DefaultValue: "test2"}, rq[1])
}

func TestGetBaseRequestHeaders(t *testing.T) {
	properties.NewProperties(
		properties.WithEnvironment("test"),
	)

	tableTest := []struct {
		Expected RequestHeader
	}{
		{Expected: RequestHeader{Key: constants_headers.ENVIRONMENT, DefaultValue: "test"}},
		{Expected: RequestHeader{Key: constants_headers.CHREF, DefaultValue: "F_COM"}},
		{Expected: RequestHeader{Key: constants_headers.CMREF, DefaultValue: "F_COM"}},
		{Expected: RequestHeader{Key: constants_headers.COUNTRY, DefaultValue: "CL"}},
		{Expected: RequestHeader{Key: constants_headers.REQUEST_ID}},
		{Expected: RequestHeader{Key: constants_headers.SITE_ID}},
		{Expected: RequestHeader{Key: constants_headers.DATADOG_PARENT_ID}},
		{Expected: RequestHeader{Key: constants_headers.DATADOG_SAMPLING_PRIORITY}},
		{Expected: RequestHeader{Key: constants_headers.DATADOG_TRACE_ID}},
		{Expected: RequestHeader{Key: constants_headers.B3_PARENT_SPAN_ID}},
		{Expected: RequestHeader{Key: constants_headers.B3_SAMPLED}},
		{Expected: RequestHeader{Key: constants_headers.B3_TRACE_ID}},
		{Expected: RequestHeader{Key: constants_headers.B3_SPAN_ID}},
		{Expected: RequestHeader{Key: constants_headers.CONTENT_TYPE, ForceValue: "application/json"}},
	}

	requestHeaders := GetBaseRequestHeaders()

	for _, tt := range tableTest {
		t.Run(tt.Expected.Key, func(t *testing.T) {
			if tt.Expected.Key == constants_headers.REQUEST_ID {
				for _, rqs := range requestHeaders {
					if rqs.Key == tt.Expected.Key {
						_, err := uuid.Parse(rqs.DefaultValue)
						assert.Nil(t, err)
						return
					}
				}
			}

			assert.Contains(t, requestHeaders, tt.Expected)
		})
	}
}
