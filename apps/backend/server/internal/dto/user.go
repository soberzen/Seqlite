package dto

type ProfileResponse struct {
	UserName  string  `json:"userName"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatarURL"`
}
