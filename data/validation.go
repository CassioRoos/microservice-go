package data

import (
	"fmt"
	"github.com/go-playground/validator"
	"regexp"
)

// ValidationError wraps the validators FieldError so we do not
// expose this to out code
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field Validation for '%s' failed on the '%s' tag.",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

type Validation struct {
	validate *validator.Validate
}

func validateLicensePlate(fl validator.FieldLevel) bool {
	// here em Brasil the license plate should be something like ABC-1234
	// regex will ensure the this is respected
	rgx := regexp.MustCompile(`[A-Z]{3}-[0-9]{4}`)
	matches := rgx.FindAllString(fl.Field().String(), -1)
	// if there is more than one or no one should raise error
	if len(matches) == 1 {
		return true
	}
	return false
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("lcplt", validateLicensePlate)
	return &Validation{validate}
}

// Validate the item
// for more detail the returned error can be cast into a
// validator.ValidationErrors collection
//
// if ve, ok := err.(validator.ValidationErrors); ok {
//			fmt.Println(ve.Namespace())
//			fmt.Println(ve.Field())
//			fmt.Println(ve.StructNamespace())
//			fmt.Println(ve.StructField())
//			fmt.Println(ve.Tag())
//			fmt.Println(ve.ActualTag())
//			fmt.Println(ve.Kind())
//			fmt.Println(ve.Type())
//			fmt.Println(ve.Value())
//			fmt.Println(ve.Param())
//			fmt.Println()
//	}
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errsValidation := v.validate.Struct(i)
	if errsValidation == nil {
		return nil
	}

	errs := errsValidation.(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}
	var returnErrs []ValidationError
	for _, err := range errs {
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}
