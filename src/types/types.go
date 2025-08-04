package types

type LoginRequest struct {
	Login      string
	Password   string
	RememberMe string
}

type LoginResponse struct {
	RememberMe string
	PhpSession string
}
