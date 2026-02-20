package constants_headers

const (
	REQUEST_ID          = "X-Request-Id"
	CONTENT_TYPE        = "Content-Type"
	ENVIRONMENT         = "X-Environment"
	CHREF               = "X-Chref"
	CMREF               = "X-Cmref"
	COUNTRY             = "X-Country"
	COMMERCE            = "X-Commerce"
	B3_TRACE_ID         = "X-B3-RequestId"
	B3_SPAN_ID          = "X-B3-SpanId"
	B3_PARENT_SPAN_ID   = "X-B3-Parent-SpanId"
	B3_SAMPLED          = "X-B3-Sampled"
	EXT_REQUEST_ID      = "Ext-Request-Id"
	EXT_REQUEST_TIMEOUT = "Ext-Request-Timeout"
	EXT_SERVICE         = "Ext-Service"
	SITE_ID             = "X-Site-Id"

	// deprecated
	DATADOG_PARENT_ID         = "X-Datadog-Parent-Id"
	DATADOG_SAMPLING_PRIORITY = "X-Datadog-Sampling-Priority"
	DATADOG_TRACE_ID          = "X-Datadog-Trace-Id"
	TRACE_ID                  = "X-Trace-Id"
)
