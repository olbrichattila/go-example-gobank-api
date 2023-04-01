package types

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRespose struct {
	Token string `json:"token"`
}
