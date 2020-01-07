package middlewares

// Interface
type Auth interface {
}

// Implementation
type authImpl struct {
}

func NewAuth() Auth {
	return &authImpl{}
}
