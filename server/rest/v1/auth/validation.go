package auth

import (
	"regexp"
	"time"

	"github.com/markovidakovic/gdsi/server/response"
)

func validateSignup(model SignupRequestModel) []response.InvalidField {
	var invFields []response.InvalidField

	if model.Name == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "name",
			Message:  "Name field is required",
			Location: "body",
		})
	}
	if model.Email == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "email",
			Message:  "Email field is required",
			Location: "body",
		})
	} else if !isValidEmail(model.Email) {
		invFields = append(invFields, response.InvalidField{
			Field:    "email",
			Message:  "Invalid email",
			Location: "body",
		})
	}
	if model.Dob == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "dob",
			Message:  "Date of birth field is required",
			Location: "body",
		})
	} else {
		if _, err := time.Parse("2006-01-02", model.Dob); err != nil {
			invFields = append(invFields, response.InvalidField{
				Field:    "dob",
				Message:  "Invalid date format",
				Location: "body",
			})
		}
	}
	if model.Gender == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "gender",
			Message:  "Gender field is required",
			Location: "body",
		})
	} else if model.Gender != "male" && model.Gender != "female" {
		invFields = append(invFields, response.InvalidField{
			Field:    "gender",
			Message:  "Invalid gender, expected male or female",
			Location: "body",
		})
	}
	if model.PhoneNumber == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "phone_number",
			Message:  "Phone number field is required",
			Location: "body",
		})
	} else if !regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,15}$`).MatchString(model.PhoneNumber) {
		invFields = append(invFields, response.InvalidField{
			Field:    "phone_number",
			Message:  "Invalid phone number",
			Location: "body",
		})
	}
	if model.Password == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "password",
			Message:  "Password field required",
			Location: "body",
		})
	}

	if len(invFields) > 0 {
		return invFields
	}

	return nil
}

func validateLogin(model LoginRequestModel) []response.InvalidField {
	var invFields []response.InvalidField

	if model.Email == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "email",
			Message:  "Email field is required",
			Location: "body",
		})
	} else if !isValidEmail(model.Email) {
		invFields = append(invFields, response.InvalidField{
			Field:    "email",
			Message:  "Invalid email",
			Location: "body",
		})
	}
	if model.Password == "" {
		invFields = append(invFields, response.InvalidField{
			Field:    "password",
			Message:  "Password field is required",
			Location: "body",
		})
	}

	if len(invFields) > 0 {
		return invFields
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
