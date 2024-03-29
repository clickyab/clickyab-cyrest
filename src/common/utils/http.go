package utils

import (
	"common/config"
	"strconv"

	"net/http"
)

// GetPageAndCount return the p and c variable from the request, if not available
// return the default value
func GetPageAndCount(r *http.Request, offset bool) (int, int) {
	p64, err := strconv.ParseInt(r.URL.Query().Get("p"), 10, 0)
	p := int(p64)
	if err != nil || p < 1 {
		p = 1
	}

	c64, err := strconv.ParseInt(r.URL.Query().Get("c"), 10, 0)
	c := int(c64)
	if err != nil || c > config.Config.Page.MaxPerPage || c < config.Config.Page.MinPerPage {
		c = config.Config.Page.PerPage
	}

	if offset {
		// If i need to make it to offset model then do it here
		p = (p - 1) * c
	}

	return p, c
}
