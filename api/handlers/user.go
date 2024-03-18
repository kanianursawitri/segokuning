package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"shopifyx/api/responses"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"shopifyx/internal/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Database *functions.User
}

type CredentialType string

const (
	Phone CredentialType = "phone"
	Email CredentialType = "email"
)

// Struct to define validation rules for both registration and login
type AuthRequest struct {
	CredentialType  CredentialType `json:"credentialType"`
	CredentialValue string         `json:"credentialValue"`
	Password        string         `json:"password"`
}

type RegisterRequest struct {
	CredentialType  CredentialType `json:"credentialType"`
	CredentialValue string         `json:"credentialValue"`
	Name            string         `json:"name"`
	Password        string         `json:"password"`
}

func (a AuthRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CredentialType, validation.Required, validation.By(func(value interface{}) error {
			if a.CredentialType != "email" && a.CredentialType != "phone" {
				return errors.New("invalid credential type")
			}
			return nil
		})),
		validation.Field(&a.CredentialValue, validation.Required),
		validation.Field(&a.Password, validation.Required, validation.Length(5, 15)),
		validation.Field(&a.CredentialValue, validation.By(func(value interface{}) error {
			strValue := value.(string)
			if a.CredentialType == "email" {
				if !isValidEmail(strValue) {
					return errors.New("invalid email format")
				}
			} else if a.CredentialType == "phone" {
				if !isValidPhoneNumber(strValue) {
					return errors.New("invalid phone number format")
				}
			}
			return nil
		})),
	)
}

func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CredentialType, validation.Required, validation.By(func(value interface{}) error {
			if a.CredentialType != "email" && a.CredentialType != "phone" {
				return errors.New("invalid credential type")
			}
			return nil
		})),
		validation.Field(&a.CredentialValue, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(5, 15)),
		validation.Field(&a.Password, validation.Required, validation.Length(5, 15)),
		validation.Field(&a.CredentialValue, validation.By(func(value interface{}) error {
			strValue := value.(string)
			if a.CredentialType == "email" {
				if !isValidEmail(strValue) {
					return errors.New("invalid email format")
				}
			} else if a.CredentialType == "phone" {
				if !isValidPhoneNumber(strValue) {
					return errors.New("invalid phone number format")
				}
			}
			return nil
		})),
	)
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) <= 7 && len(phoneNumber) >= 13 {
		return false
	}
	regex := `^\+\d{10}$`

	pattern := regexp.MustCompile(regex)

	return pattern.MatchString(phoneNumber)
}

func (u *User) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	var userValue *string
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	usr := entity.User{
		CredentialValue: req.CredentialValue,
		CredentialType:  string(req.CredentialType),
		Name:            req.Name,
		Password:        req.Password,
	}

	result, err := u.Database.Register(ctx.UserContext(), usr)
	if err != nil {
		if err.Error() == "EXISTING_USERNAME" {
			return responses.ErrorConflict(ctx, err.Error())
		}

		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	accessToken, err := utils.GenerateAccessToken(result.CredentialValue, result.Id)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	if req.CredentialType == "phone" {
		userValue = result.Phone
	} else {
		userValue = result.Email
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"name":                     result.Name,
			string(req.CredentialType): userValue,
			"accessToken":              accessToken,
		},
	})
}

func (u *User) Login(ctx *fiber.Ctx) error {
	// Parse request body
	var req AuthRequest
	var userValue *string
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	// Validate request body
	if err := req.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	usr := entity.User{
		CredentialValue: req.CredentialValue,
		CredentialType:  string(req.CredentialType),
		Password:        req.Password,
	}

	// login user
	result, err := u.Database.Login(ctx.UserContext(), usr)
	if err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			return responses.ErrorNotFound(ctx, err.Error())
		}

		if err.Error() == "INVALID_PASSWORD" {
			return responses.ErrorBadRequest(ctx, err.Error())
		}

		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.CredentialValue, result.Id)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	if req.CredentialType == "phone" {
		userValue = result.Phone
	} else {
		userValue = result.Email
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"name":                     result.Name,
			string(req.CredentialType): userValue,
			"accessToken":              accessToken,
		},
	})
}
