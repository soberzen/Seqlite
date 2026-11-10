package dto

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	UserName  string `json:"user_name"  binding:"required"`
	Email     string `json:"email"      binding:"required,email"`
	Password  string `json:"password"   binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

type RegisterResponse struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}
