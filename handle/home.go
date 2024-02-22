package handle

import (
	"net/http"

	"github.com/coreycole/go_htmx/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
