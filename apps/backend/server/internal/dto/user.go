package dto

type ProfileResponse struct {
	UserName  string  `json:"userName"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatarURL"`
}

type UpdateProfileRequest struct {
	UserName  *string `json:"userName"  binding:"omitempty,min=6,max=32,displayname"`
	AvatarURL *string `json:"avatarURL" binding:"omitempty,url"`
}

type PublicProfileResponse struct {
	ID        uint    `json:"id"`
	UserName  string  `json:"userName"`
	AvatarURL *string `json:"avatarURL"`
}
