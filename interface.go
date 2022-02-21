package main

type SensorInterface interface {
	AppendData(data []string) error
	GetType() string
	GetName() string
	GetValues() []float64
	GetAverageValue() float64
	GetStandardDeviation() float64
	GetMaxDeviationPercentage(refValue float64) float64
	GetRating(ref ReferenceInterface) string
	PrintRating(ref ReferenceInterface)
}

type ReferenceInterface interface {
	GetRefHumidity() float64
	GetRefTemperature() float64
	SetRefHumidity(hum float64)
	SetRefTemperature(tmp float64)
}