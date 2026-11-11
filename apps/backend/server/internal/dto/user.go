package dto

import "encoding/json"

type ProfileResponse struct {
	UserName  string  `json:"userName"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatarURL"`
}

// PatchString records whether a string field was omitted or explicitly set to null.
type PatchString struct {
	Present bool
	Value   *string
}

func (s *PatchString) UnmarshalJSON(data []byte) error {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	s.Present = true
	s.Value = value
	return nil
}

type UpdateProfileRequest struct {
	UserName  PatchString `json:"userName"`
	AvatarURL PatchString `json:"avatarURL"`
}

type PublicProfileResponse struct {
	ID        uint    `json:"id"`
	UserName  string  `json:"userName"`
	AvatarURL *string `json:"avatarURL"`
}
