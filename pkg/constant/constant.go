package constant

type BuildEnvironment string

const (
	PRODUCTION  BuildEnvironment = "production"
	DEVELOPMENT BuildEnvironment = "development"
	STAGING     BuildEnvironment = "staging"
)
