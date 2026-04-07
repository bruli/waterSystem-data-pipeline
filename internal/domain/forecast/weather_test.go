package forecast_test

import (
	"testing"
	"time"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
	"github.com/stretchr/testify/require"
)

func TestWeather_DryingFactor(t *testing.T) {
	weath := forecast.NewWeather(time.Now(), 22.5, 38, 0, 0, 771, time.Now())
	factor := weath.DryingFactor()
	require.NotEqual(t, 0.0, factor, "drying factor should not be 0")
}
