package servicetitan

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestClient_New(t *testing.T) {
	t.Run("returns new client", func(t *testing.T) {
		metadata := validClientOptions()
		c, err := New(metadata)
		assert.NilError(t, err)

		assert.Assert(t, c.client != nil)
		assert.Assert(t, c.session == nil)
		assert.Equal(t, c.metadata, metadata)

		authServ := c.AuthService.(authService)
		assert.Equal(t, authServ.client, c)
		assert.Equal(t, authServ.baseURL, "https://auth.servicetitan.io")
	})

	t.Run("returns error with client options are not valid", func(t *testing.T) {
		_, err := New(ClientInfo{})
		assert.ErrorIs(t, err, errMissingAppID)
	})
}
