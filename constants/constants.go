package constants

const (
	EXTERNAL_API = "EXTERNAL_API"
	AMB          = "AMB"
	MS           = "MS"
	CORE         = "CORE"
	DL           = "DL"
	SEC          = "Sec"
	COUNTRY      = "COUNTRY"
	CHANNEL      = "CHANNEL"
	CONSUMER     = "CONSUMER"

	QUERY_GENERATE = "query.generate"
	QUERY_JOB      = "query.job"
	QUERY_SYNTAX   = "query.syntax"

	// Public routes allow this many requests per minute for one visitor.
	PUBLIC_RATE_LIMIT = 40
	// Every route, including private ones, allows this many requests per minute for one visitor.
	PRIVATE_RATE_LIMIT = 40
)
