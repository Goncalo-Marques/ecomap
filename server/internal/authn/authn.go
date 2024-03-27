package authn

// service defines the authentication service structure.
type service struct {
	jwtSigningKey []byte
}

// New returns a new authentication service.
func New(jwtSigningKey []byte) *service {
	return &service{
		jwtSigningKey: jwtSigningKey,
	}
}
