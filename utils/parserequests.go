package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// reads `field` parameter for logs
func ReadFieldParam(r *http.Request) (field string, err error) {
	params := httprouter.ParamsFromContext(r.Context())

	field = params.ByName("field")
	if field == "" {
		err = errors.New("invalid field parameter")
	}

	return field, err
}

// reads `value` parameter for log field
func ReadValueParam(r *http.Request) (value string, err error) {
	params := httprouter.ParamsFromContext(r.Context())

	value = params.ByName("value")
	if value == "" {
		err = errors.New("invalid value parameter")
	}

	return value, err
}

// helper to parse specific query params into their native type.
// currently only supports a time query param specifying row age
// but more can be appended to return statement accordingly.
// alternatively, could make a custom params struct with appropriate types
func ParseQueryParams(r *http.Request) (time.Time, error) {
	rowAge := r.URL.Query().Get("age")
	since, err := parseRowAge(rowAge)
	if err != nil {
		err = fmt.Errorf("malformed 'age' query parameter: %s, %w", rowAge, err)
	}
	return since, err
}

// parses request path parameters as Unix time to time.Time
func parseRowAge(age string) (time.Time, error) {
	var since time.Time
	if age == "" {
		return since, nil
	}

	seconds, err := strconv.ParseInt(age, 10, 64)

	// return 0 Time if query param
	if err != nil || seconds == 0 {
		return since, err
	}

	eventTime := time.Now().Unix() - seconds

	return time.Unix(eventTime, 0), nil
}

// in order to accommodate shell interpolation within the CURL command,
// the data payload is enclosed in double quotes. This requires that the
// payload object's key strings are enclosed in single quotes.
// Go's default unmarshal method requires double quotes to enclose
// the key strings.
// This function makes the replacement. Feels hacky.
func PrepareRequest(req []byte) []byte {
	r := strings.Replace(string(req), "'", "\"", -1)
	return []byte(r)
}
