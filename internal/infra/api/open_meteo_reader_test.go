package api_test

import (
	"testing"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/api"
	"github.com/stretchr/testify/require"
)

func TestOpenMeteoReader(t *testing.T) {
	read := apiinfra.NewOpenMeteoReader()
	slot := forecast.Tomorrow()
	weath, err := read.Read(t.Context(), slot)
	require.NoError(t, err)
	require.Len(t, weath, 24)
}
