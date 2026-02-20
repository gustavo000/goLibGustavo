package headers_request

import (
	"github.com/google/uuid"
	"github.com/gustavo000/goLibGustavo/constants/constants_headers"
	"github.com/gustavo000/goLibGustavo/models/properties"
)

type RequestHeaders []RequestHeader

func (p RequestHeaders) ModifyHeaders(headersParser ...RequestHeader) RequestHeaders {
	for _, parser := range headersParser {
		for i, headerParser := range p {
			if headerParser.Key == parser.Key {
				p[i] = parser
				break
			}
		}
	}
	return p
}

func (p RequestHeaders) AddHeaders(headersParser ...RequestHeader) RequestHeaders {
	for _, parser := range headersParser {
		p = append(p, parser)
	}
	return p
}

type RequestHeader struct {
	Key          string
	Translate    string
	ForceValue   string
	DefaultValue string
}

func GetBaseRequestHeaders() RequestHeaders {
	return RequestHeaders{
		{Key: constants_headers.ENVIRONMENT, DefaultValue: properties.GetProperty().GetEnv()},
		{Key: constants_headers.CHREF, DefaultValue: "F_COM"},
		{Key: constants_headers.CMREF, DefaultValue: "F_COM"},
		{Key: constants_headers.COUNTRY, DefaultValue: "CL"},
		{Key: constants_headers.REQUEST_ID, DefaultValue: uuid.NewString()},
		{Key: constants_headers.TRACE_ID, DefaultValue: uuid.NewString()},
		{Key: constants_headers.SITE_ID},
		{Key: constants_headers.DATADOG_SAMPLING_PRIORITY},
		{Key: constants_headers.DATADOG_TRACE_ID},
		{Key: constants_headers.B3_PARENT_SPAN_ID},
		{Key: constants_headers.B3_SAMPLED},
		{Key: constants_headers.B3_TRACE_ID},
		{Key: constants_headers.B3_SPAN_ID},
		{Key: constants_headers.CONTENT_TYPE, ForceValue: "application/json"},
	}
}
