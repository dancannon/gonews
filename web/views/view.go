package views

import (
	"github.com/codegangsta/martini-contrib/binding"
)

type ViewErrors binding.Errors

func (errors ViewErrors) FieldError(field string) string {
	err, ok := errors.Fields[field]
	if !ok {
		return ""
	}

	return err
}
