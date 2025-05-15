package transport

import (
	"fmt"
	"net/http"
)

// GET /api/forms
func handler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")

	if key != "" {
		fmt.Fprintln(w, "key: ", key)
	} else {
		fmt.Fprintln(w, "No Quary parm")
	}
}

// POST /api/forms

// GET /api/forms/{id}

// GET /api/forms/{id}/schema

// POST /api/forms/{id}/submit

// GET /api/forms/{id}/responses

// PATCH /api/responses/{response_id}
