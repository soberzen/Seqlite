package validation

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var rules = []struct {
	tag string
	fn  validator.Func
}{
	{tag: "emoji", fn: emoji},
	{tag: "displayname", fn: displayname},
}

var registerValidation = sync.OnceValue(func() error {
	engine := binding.Validator.Engine()
	v, ok := engine.(*validator.Validate)
	if !ok {
		return fmt.Errorf("unexpected Gin validator engine %T", engine)
	}

	for _, rule := range rules {
		if err := v.RegisterValidation(rule.tag, rule.fn); err != nil {
			return fmt.Errorf("register validation %q: %w", rule.tag, err)
		}
	}
	return nil
})

func RegisterValidation() error {
	return registerValidation()
}
