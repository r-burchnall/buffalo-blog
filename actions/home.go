package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
// Used to denote the health endpoint
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(`{"status": "ok"}`))
}
