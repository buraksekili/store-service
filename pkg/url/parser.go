package url

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/buraksekili/store-service/users"
	"github.com/pkg/errors"
)

// ParseIntQueryParams parses integer query parameters
// in the given request's URL.
func ParseIntQueryParams(key string, def int, r *http.Request) (int, error) {
	strval := r.URL.Query().Get(key)
	if strings.TrimSpace(strval) == "" {
		return def, nil
	}
	val, err := strconv.ParseInt(strval, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, users.ErrInvalidRequestPath.Error())
	}
	return int(val), err
}
