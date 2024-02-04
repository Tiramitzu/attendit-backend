package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var Required = validation.Required.Error("tidak boleh kosong")

var passwordRule = []validation.Rule{
	Required,
	validation.Length(8, 32).Error("harus memiliki 8-32 karakter"),
	validation.Match(regexp.MustCompile("^\\S+$")).Error("tidak boleh mengandung spasi"),
}

type PaidLeaveRequest struct {
	UserId    string `json:"userId"`
	Reason    string `json:"reason"`
	StartDate string `json:"startDate"`
	Days      int    `json:"days"`
}

func (a PaidLeaveRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UserId, Required),
		validation.Field(&a.Reason, Required),
		validation.Field(&a.StartDate, Required),
		validation.Field(&a.Days, Required, is.Digit.Error("harus berupa angka"), validation.Min(1).Error("harus lebih dari 0 hari")),
	)
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

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

func (a CreateUser) Validate() error {
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
	Password string `json:"password"`
}

func (a ModifyUserRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, is.Email.Error(" tidak valid")),
		validation.Field(&a.FullName, validation.Length(3, 0).Error("harus lebih dari 3 karakter")),
		validation.Field(&a.Phone, is.Digit.Error(" tidak valid"), validation.Length(11, 14).Error("harus terdiri dari 11-14 karakter")),
		validation.Field(&a.Password, passwordRule...),
	)
}
