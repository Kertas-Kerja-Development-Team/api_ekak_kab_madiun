package user

type UserChangePasswordRequest struct {
	Nip         string `json:"nip"`
	OldPassword string `json:"old_password"`
	Password1   string `json:"password1" validate:"required"`
	Password2   string `json:"password2" validate:"required"`
}
