package model

type RegisterRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	UserID      uint64 `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FilterFindUser struct {
	Username *string
	ID       []uint64
}
