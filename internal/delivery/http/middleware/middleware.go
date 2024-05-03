package middleware

type Middleware struct {
	JWTSecretKey string
}

func NewMiddleware(jwtSecretKey string) *Middleware {
	return &Middleware{
		JWTSecretKey: jwtSecretKey,
	}
}
