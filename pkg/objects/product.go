package objects

import (
	"errors"
	"fmt"
	"time"
)

type Product struct {
	Id        int32
	Name      *string
	LogoUrl   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const requiredField string = "missing field \"%s\" required for product"

func (p Product) Validate() error {
	if err := p.ValidateName(); err != nil {
		return err
	}
	if err := p.ValidateLogoUrl(); err != nil {
		return err
	}
	return nil
}

func (p Product) ValidateName() error {
	if p.Name == nil {
		errMsg := fmt.Sprintf(requiredField, "Name")
		return errors.New(errMsg)
	}
	return nil
}

func (p Product) ValidateLogoUrl() error {
	if p.LogoUrl == nil {
		errMsg := fmt.Sprintf(requiredField, "LogoUrl")
		return errors.New(errMsg)
	}
	return nil
}
