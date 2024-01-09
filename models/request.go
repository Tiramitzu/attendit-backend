package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var Required = validation.Required.Error("tidak boleh kosong")

var passwordRule = []validation.Rule{
	validation.Required,
	validation.Length(8, 32),
	validation.Match(regexp.MustCompile("^\\S+$")).Error("tidak boleh mengandung spasi"),
}

type ScheduleRequest struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func (a ScheduleRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, Required),
		validation.Field(&a.StartTime, Required),
		validation.Field(&a.EndTime, Required),
	)
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, Required, is.Email.Error("tidak valid")),
		validation.Field(&a.Password, passwordRule...),
		validation.Field(&a.FullName, validation.Length(3, 0).Error("harus lebih dari 3 karakter")),
		validation.Field(&a.Phone, validation.Length(11, 14).Error("harus terdiri dari 11-14 karakter")),
	)
}

type CheckInRequest struct {
	Status string `json:"status"`
}

func (a CheckInRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Status, Required),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, Required, is.Email.Error("tidak valid")),
		validation.Field(&a.Password, passwordRule...),
	)
}

type ModifyUserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

func (a ModifyUserRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, is.Email.Error(" tidak valid")),
		validation.Field(&a.FullName, validation.Length(3, 0).Error("harus lebih dari 3 karakter")),
		validation.Field(&a.Phone, validation.Length(11, 14).Error("harus terdiri dari 11-14 karakter")),
	)
}
