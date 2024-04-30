package validators

import (
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"net/mail"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	// Custom validation for Emails fields.
	_ = validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := mail.ParseAddress(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// FIXME
	//if err != nil {
	//	// Make error message for each invalid field.
	//	for _, err := range err.(validator.ValidationErrors) {
	//		fields[err.Field()] = err.Error()
	//	}
	//}

	return fields
}

// ParseAndValidate parser
func ParseAndValidate(ctx fiber.Ctx, out interface{}) error {
	resp := core_dtos.NewResponse(ctx)

	// Checking received data from JSON body. Return status 400 and error message.
	if err := ctx.Bind().Body(&out); err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		resp.SetMessage(err.Error())
		return resp.JSON()
	}

	// Create a new validator for a RegistrationRequestDto.
	validate := NewValidator()

	// Validate fields.
	if err := validate.Struct(&out); err != nil {
		resp.SetStatus(fiber.StatusBadRequest)
		resp.SetMessage(err.Error()) // ValidatorErrors(err)
	}

	return resp.JSON()
}
