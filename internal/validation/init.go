package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func Instance() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})

	return validate
}
