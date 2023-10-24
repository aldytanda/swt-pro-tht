package handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aldytanda/swt-pro-tht/generated"
)

type ValidationErrorResp struct {
	Errors []ValidationError `json:"validation_errors"`
}

func (e ValidationErrorResp) Error() string {
	var errString string

	for _, v := range e.Errors {
		if len(v.ErrRules) == 0 {
			continue
		}
		if len(errString) > 0 {
			errString += ", "
		}
		errString += fmt.Sprintf("%s: %s", v.FieldName, v.ErrRules)

	}

	return errString
}

func (e *ValidationErrorResp) FromMap(v map[string][]string) bool {
	e.Errors = []ValidationError{}

	for field, errs := range v {
		if len(errs) > 0 {
			e.Errors = append(e.Errors, ValidationError{
				FieldName: field,
				ErrRules:  errs,
			})
		}
	}

	return len(e.Errors) > 0
}

func (e *ValidationErrorResp) ToResponseErrors() *[]generated.ValidationError {
	resp := []generated.ValidationError{}

	for _, v := range e.Errors {
		if len(v.ErrRules) > 0 {
			resp = append(resp, generated.ValidationError{
				FieldName: v.FieldName,
				ErrRules:  v.ErrRules,
			})
		}
	}

	return &resp
}

type ValidationError struct {
	FieldName string   `json:"field_name"`
	ErrRules  []string `json:"err_rules"`
}

func validateRegister(ctx context.Context, v generated.RegisterRequest) error {
	var errs ValidationErrorResp

	errMap := map[string][]string{
		"name":     {},
		"password": {},
		"phone":    {},
	}

	if len(v.Phone) < 10 {
		errMap["phone"] = append(errMap["phone"], "Minimum_Length_10")
	}

	if len(v.Phone) > 13 {
		errMap["phone"] = append(errMap["phone"], "Maximum_Length_13")
	}

	if !strings.HasPrefix(v.Phone, "+62") {
		errMap["phone"] = append(errMap["phone"], "Must_Start_with_'+62'")
	}

	if len(v.Name) < 3 {
		errMap["name"] = append(errMap["name"], "Minimum_Length_3")
	}

	if len(v.Name) > 60 {
		errMap["name"] = append(errMap["name"], "Maximum_Length_60")
	}

	if len(v.Password) < 6 {
		errMap["password"] = append(errMap["password"], "Minimum_Length_6")
	}

	if len(v.Password) > 64 {
		errMap["password"] = append(errMap["password"], "Maximum_Length_64")
	}

	reg := regexp.MustCompile(`^(.*[A-Z]+.*)$`)
	if !reg.MatchString(v.Password) {
		errMap["password"] = append(errMap["password"], "Must_Contain_Uppercase")
	}

	reg = regexp.MustCompile(`^(.*[0-9]+.*)$`)
	if !reg.MatchString(v.Password) {
		errMap["password"] = append(errMap["password"], "Must_Contain_Numeric")
	}

	reg = regexp.MustCompile(`^(.*[-+_!@#$%^&*., ?])+.*$`)
	if !reg.MatchString(v.Password) {
		errMap["password"] = append(errMap["password"], "Must_Contain_Special_Chars")
	}

	if errs.FromMap(errMap) {
		return errs
	}

	return nil
}

func validateUpdateUser(ctx context.Context, v generated.UpdateProfileRequest) error {
	var errs ValidationErrorResp

	errMap := map[string][]string{
		"payload": {},
		"name":    {},
		"phone":   {},
	}

	if v.Phone != nil {
		if len(*v.Phone) < 10 {
			errMap["phone"] = append(errMap["phone"], "Minimum Length 10")
		}
		if len(*v.Phone) > 13 {
			errMap["phone"] = append(errMap["phone"], "Maximum Length 10")
		}
		if !strings.HasPrefix(*v.Phone, "+62") {
			errMap["phone"] = append(errMap["phone"], "Must Start with '+62'")
		}
	}

	if v.Name != nil {
		if len(*v.Name) < 3 {
			errMap["name"] = append(errMap["name"], "Minimum Length 3")
		}

		if len(*v.Name) > 60 {
			errMap["name"] = append(errMap["name"], "Maximum Length 60")
		}
	}

	if v.Phone == nil && v.Name == nil {
		errMap["payload"] = append(errMap["payload"], "Must contain at least 1 field to update")
	}

	if errs.FromMap(errMap) {
		return errs
	}

	return nil
}
