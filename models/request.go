package models

import (
	db "attendit/backend/models/db"
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

type CreateCompanyRequest struct {
	Author       string   `json:"author"`
	Name         string   `json:"name"`
	IPAddresses  []string `json:"ipAddresses"`
	CheckInTime  string   `json:"checkInTime"`
	CheckOutTime string   `json:"checkOutTime"`
}

func (a CreateCompanyRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Author, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.IPAddresses, validation.Required),
		validation.Field(&a.CheckInTime, validation.Required),
		validation.Field(&a.CheckOutTime, validation.Required),
	)
}

type InsertMembersToCompanyRequest struct {
	Members []db.Member `json:"members"`
}

func (a InsertMembersToCompanyRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Members, validation.Required),
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
