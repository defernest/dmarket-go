package dmarket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	url := "https://api.dmarket.com"
	publicKey := "f7235e19236478f20b60b6240c49afb4d5a9970eb2228cb44b4123047eb89ec3"
	privateKey := "255ff758b9252c04d29ae19a88ef8be2dc8d3654a90037b8881937b81cfcf87bf7235e1a236478f20b60b6240c49afb4d5a9970eb2228cb44b4123047eb89ec3"
	t.Run("success", func(t *testing.T) {
		apiClient, err := NewClient(url, publicKey, privateKey)
		require.NoError(t, err)
		require.Equal(t, url, apiClient.DefaultClient.baseURL.Scheme+"://"+apiClient.DefaultClient.baseURL.Host)
		require.Equal(t, publicKey, apiClient.DefaultClient.publicKey)
		require.Equal(t, privateKey, apiClient.DefaultClient.privateKey)
	})
	t.Run("err: wrong keys len", func(t *testing.T) {
		_, err := NewClient("client://localhost", "", "")
		var keyErr errorBadKeys
		require.ErrorAs(t, err, &keyErr)
	})
	t.Run("err: wrong host url", func(t *testing.T) {
		_, err := NewClient("", publicKey, privateKey)
		require.Error(t, err)
	})
}
