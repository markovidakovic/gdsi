package auth

import (
	"regexp"
	"time"

	"github.com/markovidakovic/gdsi/server/pkg/response"
)

func validateSignup(model SignupRequestModel) []response.InvalidField {
	var invFields []response.InvalidField

	if model.FirstName == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "first_name",
			Error: "First name field is required",
		})
	}
	if model.LastName == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "last_name",
			Error: "Last name field is required",
		})
	}
	if model.Email == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "email",
			Error: "Email field is required",
		})
	} else if !isValidEmail(model.Email) {
		invFields = append(invFields, response.InvalidField{
			Field: "email",
			Error: "Invalid email",
		})
	}
	if model.Dob == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "dob",
			Error: "Date of birth field is required",
		})
	} else {
		if _, err := time.Parse("2006-01-02", model.Dob); err != nil {
			invFields = append(invFields, response.InvalidField{
				Field: "dob",
				Error: "Invalid date format",
			})
		}
	}
	if model.Gender == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "gender",
			Error: "Gender field is required",
		})
	} else if model.Gender != "male" && model.Gender != "female" {
		invFields = append(invFields, response.InvalidField{
			Field: "gender",
			Error: "Invalid gender, expected male or female",
		})
	}
	if model.PhoneNumber == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "phone_number",
			Error: "Phone number field is required",
		})
	} else if !regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,15}$`).MatchString(model.PhoneNumber) {
		invFields = append(invFields, response.InvalidField{
			Field: "phone_number",
			Error: "Invalid phone number",
		})
	}
	if model.Password == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "password",
			Error: "Password field required",
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
			Field: "email",
			Error: "Email field is required",
		})
	} else if !isValidEmail(model.Email) {
		invFields = append(invFields, response.InvalidField{
			Field: "email",
			Error: "Invalid email",
		})
	}
	if model.Password == "" {
		invFields = append(invFields, response.InvalidField{
			Field: "password",
			Error: "Password field is required",
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
