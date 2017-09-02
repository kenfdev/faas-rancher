// Copyright (c) 2017 Ken Fukuyama
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// VarsHandler a wrapper type for mux.Vars
type VarsHandler func(w http.ResponseWriter, r *http.Request, vars map[string]string)

// ServeHTTP a wrapper function for mux.Vars
func (vh VarsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vh(w, r, vars)
}
