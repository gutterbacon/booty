package service_test

import (
	ks "github.com/kardianos/service"
	"github.com/stretchr/testify/require"
	"go.amplifyedge.org/booty-v2/internal/service"
	"testing"
)

func TestNewService(t *testing.T) {
	_, err := service.NewService(&ks.Config{
		Name:        "blah",
		DisplayName: "blash",
		Description: "blash blash blash",
		Option:      map[string]interface{}{},
	})
	require.NoError(t, err)
}
