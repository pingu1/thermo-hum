package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
    "strings"
)

// Valid devices
var Thermometer = "thermometer"
var HumiditySensor = "humidity"

// Ratings
var ThermometerUltraPrecise = "ultra precise"
var ThermometerVeryPrecise = "very precise"
var ThermometerPrecise = "precise"
var HumidityAccepted = "accepted"
var HumidityRejected = "rejected"

// Control Values
var ThermometerAvgRange = 0.5
var ThermometerUltraPreciseSD = 3
var ThermometerVeryPreciseSD = 5
var HumidityAcceptedRange = 0.01

/** Defining reference Temperature and Humidity **/
type RefTempteratureHumidity struct {
    refTemperature float64
    refHumidity  float64
}

func NewRefTemperatureHumidity(temp float64, hum float64) ReferenceInterface {
    return &RefTempteratureHumidity{
        refTemperature: temp,
        refHumidity: hum,
    }
}

func (rth *RefTempteratureHumidity) GetRefTemperature() float64 {
    return rth.refTemperature
}

func (rth *RefTempteratureHumidity) SetRefTemperature(temp float64) {
    rth.refTemperature = temp
}

func (rth *RefTempteratureHumidity) GetRefHumidity() float64 {
    return rth.refHumidity
}

func (rth *RefTempteratureHumidity) SetRefHumidity(hum float64) {
    rth.refHumidity = hum
}

/**
 * Extracting reference values
 */
 func ExtractRef(refLine string) (ReferenceInterface, error) {
    refTH := &RefTempteratureHumidity{}

    // Make sure header has the right number of elements
    ref := strings.Split(strings.TrimSpace(refLine), " ")
    if len(ref) != 3 {
		error := errors.New("Error while parsing the header: not enough elements")
		return refTH, error
    }

	header := ref[0]
	if header != "reference" {
		error := errors.New("First line doesn't seem to contain the reference, stopping now")
		return refTH, error
	}

    // Make sure the ref. Temperature is set properly
    temp, err := strconv.ParseFloat(ref[1], 64)
    if err != nil {
		return refTH, err
    }

    // Make sure the ref. Humidity is set properly
    hum, err := strconv.ParseFloat(ref[2], 64)
    if err != nil {
		return refTH, err
    }

    refTH.SetRefTemperature(temp)
    refTH.SetRefHumidity(hum)

    return refTH, nil
}

/** Defining sensors **/
type Sensor struct {
    sensorType string
    sensorName string
    sensorValues []float64
}

func NewSensor(sType string, sName string) SensorInterface {
	return &Sensor{
        sensorType: sType,
        sensorName: sName,
		sensorValues: nil,
    }
}

func (s *Sensor) AppendData(data []string) {
	// input param is composed of 3 elements: date, sensor name and value recorded
	if data[1] != s.sensorName {
		return
	}

	value, err := strconv.ParseFloat(data[2], 64)
	if err != nil {
		errorMsg := "Error while parsing the recorded measure for devide " + s.sensorName + " :"
		fmt.Println(errorMsg + err.Error())
        return
	}

	s.sensorValues = append(s.sensorValues, value)
}

func (s *Sensor) GetType() string {
	return s.sensorType
}

func (s *Sensor) GetName() string {
	return s.sensorName
}

func (s *Sensor) GetValues() []float64 {
	return s.sensorValues
}

func (s *Sensor) GetAverageValue() float64 {
	var sum float64
	nbrValues := len(s.sensorValues)

	for i := 0; i < nbrValues; i++ {
		sum += s.sensorValues[i]
	}

	return sum/float64(nbrValues)
}

func (s *Sensor) GetStandardDeviation() float64 {
	var sd float64

	avg := s.GetAverageValue()
	nbrValues := len(s.sensorValues)

	for i := 0; i < nbrValues; i++ {
		sd += math.Pow(s.sensorValues[i] - avg, 2)
	}

	return math.Sqrt(sd/float64(nbrValues))
}

func (s *Sensor) GetMaxDeviationPercentage(refValue float64) float64 {
	maxDeviation := 0.0

	nbrValues := len(s.sensorValues)

	for i := 0; i < nbrValues; i++ {
		deviation := getDeviation(refValue, s.sensorValues[i])/refValue
		if deviation > maxDeviation {
			maxDeviation = deviation
		}
	}

	return maxDeviation
}

func (s *Sensor) GetRating(ref ReferenceInterface) string {
	// Reject invalid sensor types
	if !s.isValidSensorType() {
		return "Invalid sensor type for type " + s.sensorType + ". Checking next sensor"
	}

	if s.sensorType == Thermometer {
		return s.getThermometerRating(ref.GetRefTemperature())
	}

	return s.getHumiditySensorRating(ref.GetRefHumidity())
}

func (s *Sensor) PrintRating(ref ReferenceInterface) {
	rating := s.GetRating(ref)
	fmt.Printf("%s: %s\n", s.sensorName, rating)
}

func (s *Sensor) isValidSensorType() bool {
	return s.sensorType == Thermometer || s.sensorType == HumiditySensor
}

func (s *Sensor) getThermometerRating(refTemperature float64) string {
	// We want the average temperature to be lower than 0.5 degrees from ref
	
	deviation := getDeviation(refTemperature, s.GetAverageValue())
	if deviation > ThermometerAvgRange {
		return ThermometerPrecise
	}

	// We also need to check the Standard Deviation
	if s.GetStandardDeviation() <= float64(ThermometerUltraPreciseSD) {
		return ThermometerUltraPrecise
	}

	if s.GetStandardDeviation() <= float64(ThermometerVeryPreciseSD) {
		return ThermometerVeryPrecise
	}

	return ThermometerPrecise
}

func (s *Sensor) getHumiditySensorRating(refHumidity float64) string {
	// For Humidity sensors, we only care about the readings accuracy

	maxDeviation := s.GetMaxDeviationPercentage(refHumidity)
	if maxDeviation <= HumidityAcceptedRange {
		return HumidityAccepted
	}

	return HumidityRejected
}


// Additional helper
func getDeviation(refValue float64, value float64) float64 {
	return math.Abs(refValue - value)
}