package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var passwordRule = []validation.Rule{
	validation.Required,
	validation.Length(8, 32),
	validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	UserName    string `json:"username"`
	Phone       string `json:"phone"`
}

func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
		validation.Field(&a.UserName, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.DisplayName, validation.Length(3, 64)),
		validation.Field(&a.Phone, validation.Length(11, 14)),
	)
}

type CheckInRequest struct {
	IpAddress string `json:"ipAddress"`
	Status    string `json:"status"`
}

func (a CheckInRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.IpAddress, validation.Required),
		validation.Field(&a.Status, validation.Required),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
	)
}

type ModifyUserRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	UserName    string `json:"username"`
	Phone       string `json:"phone"`
}

func (a ModifyUserRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, is.Email),
		validation.Field(&a.UserName, validation.Length(3, 64)),
		validation.Field(&a.DisplayName, validation.Length(3, 64)),
		validation.Field(&a.Phone, validation.Length(11, 14)),
	)
}
