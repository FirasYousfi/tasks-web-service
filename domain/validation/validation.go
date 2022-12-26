package validation

import (
	"errors"
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
)

var (
	// ErrEmptyField when fields cannot be empty
	ErrEmptyField = errors.New("field cannot be empty")
	// ErrInvalidLength when the content of the field is invalid
	ErrInvalidLength = errors.New("field length is invalid")
)

// ValidateParams Validates the parameters given in the request, req is returned also in case in the future we want to set some default values here
func ValidateParams(req *entity.TaskDescription) (*entity.TaskDescription, error) {
	err := ValidateTitle(req.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %v", err)
	}

	err = ValidateDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid description: %v", err)
	}

	err = ValidatePriority(req.Priority)
	if err != nil {
		return nil, fmt.Errorf("invalid priority: %v", err)
	}

	req.Status, err = ValidateStatus(req.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %v", err)
	}
	return req, nil
}

func ValidatePriority(priority int) error {
	if priority < 0 || priority > 10 {
		return fmt.Errorf("invalid priority range, should be a value from 0 to 10")
	}
	return nil
}

func ValidateDescription(description string) error {
	if len(description) > 500 {
		return fmt.Errorf("%s: description length should be under 500 characters", ErrInvalidLength)
	}
	return nil
}

func ValidateTitle(title string) error {
	if title == "" {
		return ErrEmptyField
	}
	if len(title) > 100 {
		return fmt.Errorf("%s: title length should be under 100 characters", ErrInvalidLength)
	}
	return nil
}

func ValidateStatus(status entity.Status) (entity.Status, error) {
	if status == "" {
		return entity.OnHold, nil
	}
	switch status {
	case entity.New, entity.Active, entity.Closed, entity.OnHold:
		return status, nil
	}
	return "", errors.New("invalid status type")
}
