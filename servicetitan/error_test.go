package servicetitan

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestError_Error(t *testing.T) {
	err := Error{
		StatusCode:  401,
		RequestPath: "some/path",
		Message:     "missing ST-App-Key",
	}

	assert.Equal(t, err.Error(), `ServiceTitan error: missing ST-App-Key got response code 401 for request path "some/path"`)
}
