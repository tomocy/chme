package chme

import (
	"net/http"
)

// DefaultInputName is used as the key of the value to which request method is changed.
const DefaultInputName = "_method"

var defaultChme = NewChme(DefaultInputName)

var changeableMethods = map[string]bool{
	http.MethodPut:    true,
	http.MethodPatch:  true,
	http.MethodDelete: true,
}

// Chme provides methods which change request method to others.
type Chme interface {
	ChangePostToHiddenMethod(next http.Handler) http.Handler
}

type chme struct {
	inputName string
}

// NewChme returns new instance which implements Che interface.
func NewChme(name string) Chme {
	return &chme{
		inputName: name,
	}
}

// ChangePostToHiddenMethod changes POST to the method set in FormValue named "_method".
func ChangePostToHiddenMethod(next http.Handler) http.Handler {
	return defaultChme.ChangePostToHiddenMethod(next)
}

// ChangePostToHiddenMethod changes POST to the method set in FormValue named when NewChme.
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
