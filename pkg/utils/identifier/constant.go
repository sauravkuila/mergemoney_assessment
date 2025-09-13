package identifier

type IdentifierType string

const (
	IDENTIFIER_BASE62    IdentifierType = "base62"
	IDENTIFIER_SNOWFLAKE IdentifierType = "snowflake"
)
