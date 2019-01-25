package chme

import (
	"net/http"
)

const DefaultInputName = "_method"

var defaultChme = NewChme(DefaultInputName)

var changeableMethods = map[string]bool{
	http.MethodPut:    true,
	http.MethodPatch:  true,
	http.MethodDelete: true,
}

type Chme interface {
	ChangePostToHiddenMethod(next http.Handler) http.Handler
}

type chme struct {
	inputName string
}

func NewChme(name string) Chme {
	return &chme{
		inputName: name,
	}
}

func ChangePostToHiddenMethod(next http.Handler) http.Handler {
	return defaultChme.ChangePostToHiddenMethod(next)
}

func (c chme) ChangePostToHiddenMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		method := r.FormValue(c.inputName)
		if ok := changeableMethods[method]; ok {
			r.Method = method
		}

		next.ServeHTTP(w, r)
	})
}
