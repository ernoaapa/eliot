package model

import (
	"log"
	"regexp"
	"strings"
	"sync"

	imageref "github.com/containerd/containerd/reference"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		validate.RegisterValidation("hasName", func(fl validator.FieldLevel) bool {
			return hasValidName(fl.Field().Interface().(Metadata))
		})
		validate.RegisterValidation("imageRef", func(fl validator.FieldLevel) bool {
			return isValidImageReference(fl.Field().Interface().(string))
		})
		validate.RegisterValidation("alphanumOrDash", func(fl validator.FieldLevel) bool {
			return isAlphanumericOrDash(fl.Field().Interface().(string))
		})
		validate.RegisterValidation("noSpaces", func(fl validator.FieldLevel) bool {
			return !containsSpaces(fl.Field().Interface().(string))
		})
		validate.RegisterValidation("empty", func(fl validator.FieldLevel) bool {
			return isEmpty(fl.Field().Interface())
		})
		validate.RegisterValidation("envKeyValuePair", func(fl validator.FieldLevel) bool {
			return IsValidEnvKeyValuePair(fl.Field().Interface().(string))
		})
	})
	return validate
}

func hasValidName(metadata Metadata) bool {
	if len(metadata.Name) == 0 {
		return false
	}
	return true
}

func isValidImageReference(ref string) bool {
	_, err := imageref.Parse(ref)
	return err == nil
}

func isAlphanumericOrDash(value string) bool {
	match, err := regexp.MatchString("^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$", value)
	if err != nil {
		log.Fatalf("Invalid regexp definition in isAlphanumericOrDash check: %s", err)
	}
	return match
}

func containsSpaces(value string) bool {
	return strings.Contains(value, " ")
}

func isEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return value.(string) == ""
	default:
		log.Fatalf("isempty validation supports only string, not %T", v)
	}
	return false
}

func IsValidEnvKeyValuePair(value string) bool {
	if value == "" {
		return false
	}

	parts := strings.SplitN(value, "=", 2)
	for _, part := range parts {
		if !isAlphanumericOrDash(part) {
			return false
		}
	}
	return true
}

// Validate validates given pod definitions
func Validate(pods []Pod) error {
	validate := getValidator()
	for _, pod := range pods {
		err := validate.Struct(pod)
		if err != nil {

			// this check is only needed when your code could produce
			// an invalid value for validation such as interface with nil
			// value most including myself do not usually have code like this.
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return err
			}

			return err.(validator.ValidationErrors)
		}
	}

	return nil
}
