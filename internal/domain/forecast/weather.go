package forecast

import "time"

type Weather struct {
	hour                     time.Time
	temperature              float64
	relativeHumidity         int
	precipitationProbability float64
	cloudCover               int
	shortwaveRadiation       float64
	generatedAt              time.Time
}

func (w Weather) GeneratedAt() time.Time {
	return w.generatedAt
}

func (w Weather) Hour() time.Time {
	return w.hour
}

func (w Weather) Temperature() float64 {
	return w.temperature
}

func (w Weather) RelativeHumidity() int {
	return w.relativeHumidity
}

func (w Weather) PrecipitationProbability() float64 {
	return w.precipitationProbability
}

func (w Weather) CloudCover() int {
	return w.cloudCover
}

func (w Weather) ShortwaveRadiation() float64 {
	return w.shortwaveRadiation
}

func (w Weather) DryingFactor() float64 {
	return (w.temperature / 40.0) * (1.0 - (float64(w.relativeHumidity) / 100)) * (w.shortwaveRadiation / 1000)
}

func NewWeather(
	hour time.Time,
	temperature float64,
	relativeHumidity int,
	precipitationProbability float64,
	cloudCover int,
	shortwaveRadiation float64,
	generatedAt time.Time,
) *Weather {
	return &Weather{
		hour:                     hour,
		temperature:              temperature,
		relativeHumidity:         relativeHumidity,
		precipitationProbability: precipitationProbability,
		cloudCover:               cloudCover,
		shortwaveRadiation:       shortwaveRadiation,
		generatedAt:              generatedAt,
	}
}
