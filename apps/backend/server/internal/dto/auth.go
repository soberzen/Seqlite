package dto

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type RegisterRequest struct {
	UserName  string  `json:"userName"  binding:"required"`
	Email     string  `json:"email"     binding:"required,email"`
	Password  string  `json:"password"  binding:"required"`
	AvatarURL *string `json:"avatarURL"`
}

type RegisterResponse struct {
	UserName  string  `json:"userName"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatarURL"`
}
