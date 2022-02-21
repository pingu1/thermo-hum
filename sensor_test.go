package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing RefTemperatureHumidity objects

func TestNewRefTemperatureHumidity_HappyPath(t *testing.T) {
	expectedTemperature := 1234.56
	expectedHumidity := 123.45

	res := NewRefTemperatureHumidity(expectedTemperature, expectedHumidity)

	assert.NotNil(t, res)
	assert.Equal(t, expectedTemperature, res.GetRefTemperature())
	assert.Equal(t, expectedHumidity, res.GetRefHumidity())
}

func TestRefTemperatureHumiditySetters_HappyPath(t *testing.T) {
	expectedTemperature := 1234.56
	expectedHumidity := 123.45

	res := NewRefTemperatureHumidity(0.0, 0.0)

	res.SetRefTemperature(expectedTemperature)
	res.SetRefHumidity(expectedHumidity)

	assert.Equal(t, expectedTemperature, res.GetRefTemperature())
	assert.Equal(t, expectedHumidity, res.GetRefHumidity())
}

func TestExtractRefValues_HappyPath(t *testing.T) {
	line := "reference 70.0 45.0"

	res, err := ExtractRef(line)

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, 70.0, res.GetRefTemperature())
	assert.Equal(t, 45.0, res.GetRefHumidity())
}

func TestExtractRefValues_NotEnoughElements(t *testing.T) {
	line := "hello world"

	_, err := ExtractRef(line)

	assert.NotNil(t, err)
	assert.Equal(t, "Error while parsing the header: not enough elements", err.Error())
}

func TestExtractRefValues_NotAReference(t *testing.T) {
	line := "Potato 17 42"

	_, err := ExtractRef(line)

	assert.NotNil(t, err)
	assert.Equal(t, "First line doesn't seem to contain the reference, stopping now", err.Error())
}

func TestExtractRefValues_RefTemperatureNotANumber(t *testing.T) {
	line := "reference hello 42"

	_, err := ExtractRef(line)

	assert.NotNil(t, err)
	assert.Equal(t, "strconv.ParseFloat: parsing \"hello\": invalid syntax", err.Error())
}

func TestExtractRefValues_RefHumidityNotANumber(t *testing.T) {
	line := "reference 17 world"

	_, err := ExtractRef(line)

	assert.NotNil(t, err)
	assert.Equal(t, "strconv.ParseFloat: parsing \"world\": invalid syntax", err.Error())
}

// Testing Sensor objects

func TestNewSensor_HappyPath(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato sensor"

	res := NewSensor(expectedType, expectedName)

	assert.NotNil(t, res)
	assert.Equal(t, expectedType, res.GetType())
	assert.Equal(t, expectedName, res.GetName())
	assert.Nil(t, res.GetValues())
}

func TestAppendData_HappyPathEmptySensor(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	expectedValue := 1.23

	sensor := NewSensor(expectedType, expectedName)
	assert.Nil(t, sensor.GetValues())

	lineData := []string{"2000-01-01T00:00:00", "Potato-sensor", "1.23"}
	sensor.AppendData(lineData)

	assert.NotNil(t, sensor.GetValues())
	assert.Equal(t, []float64{expectedValue}, sensor.GetValues())
}

func TestAppendData_HappyPathSensorWithData(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	expectedValue := 1.23

	sensor := NewSensor(expectedType, expectedName)

	lineData := []string{"2000-01-01T00:00:00", "Potato-sensor", "1.23"}
	sensor.AppendData(lineData)
	sensor.AppendData(lineData)

	assert.NotNil(t, sensor.GetValues())
	assert.Equal(t, []float64{expectedValue, expectedValue}, sensor.GetValues())
}

func TestAppendData_NotTheRightSensor(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"

	sensor := NewSensor(expectedType, expectedName)
	assert.Nil(t, sensor.GetValues())

	lineData := []string{"2000-01-01T00:00:00", "Potato-sensor-but-another-one", "1.23"}
	err := sensor.AppendData(lineData)

	assert.NotNil(t, err)
	assert.Equal(t, "Data is not for the right sensor", err.Error())
	assert.Nil(t, sensor.GetValues())
}

func TestAppendData_BadDataType(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"

	sensor := NewSensor(expectedType, expectedName)
	assert.Nil(t, sensor.GetValues())

	lineData := []string{"2000-01-01T00:00:00", "Potato-sensor", "this is not a number"}
	err := sensor.AppendData(lineData)

	assert.NotNil(t, err)
	assert.Equal(t, "Error while parsing the recorded measure for devide Potato-sensor :strconv.ParseFloat: parsing \"this is not a number\": invalid syntax", err.Error())
	assert.Nil(t, sensor.GetValues())
}

func TestAppendData_DataIsInteger(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	expectedValue := float64(1)

	sensor := NewSensor(expectedType, expectedName)
	assert.Nil(t, sensor.GetValues())

	lineData := []string{"2000-01-01T00:00:00", "Potato-sensor", "1"}
	err := sensor.AppendData(lineData)

	assert.Nil(t, err)
	assert.Equal(t, []float64{expectedValue}, sensor.GetValues())
}

func TestGetAverageValue_HappyPath(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	value1 := "1"
	value2 := "2"
	expectedAvg := 1.5

	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}

	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetAverageValue()

	assert.NotNil(t, res)
	assert.Equal(t, expectedAvg, res)
}

func TestGetAverageValue_NoData(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	expectedAvg := float64(0)
	
	sensor := NewSensor(expectedType, expectedName)

	res := sensor.GetAverageValue()
	
	assert.NotNil(t, res)
	assert.Equal(t, expectedAvg, res)
}

func TestGetStandardDeviation_HappyPath(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	value1 := "1"
	value2 := "2"
	expectedStdDev := 0.7071067811865476

	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}

	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetStandardDeviation()

	assert.NotNil(t, res)
	assert.Equal(t, expectedStdDev, res)
}

func TestGetStandardDeviation_NoData(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	expectedAvg := float64(0)
	
	sensor := NewSensor(expectedType, expectedName)

	res := sensor.GetStandardDeviation()
	
	assert.NotNil(t, res)
	assert.Equal(t, expectedAvg, res)
}

func TestGetMaxDeviationPercentage_HappyPath(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	refValue := 1.25
	value1 := "1"
	value2 := "2"
	expectedMaxDev := 0.6

	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetMaxDeviationPercentage(refValue)

	assert.NotNil(t, res)
	assert.Equal(t, expectedMaxDev, res)
}

func TestGetMaxDeviationPercentage_NoData(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	refValue := 1.25
	expectedAvg := float64(0)
	
	sensor := NewSensor(expectedType, expectedName)

	res := sensor.GetMaxDeviationPercentage(refValue)
	
	assert.NotNil(t, res)
	assert.Equal(t, expectedAvg, res)
}

func TestGetRating_HappyPath_ThermometerUltraPrecise(t *testing.T) {
	expectedType := "thermometer"
	expectedName := "therm-1"
	value1 := "69.5"
	value2 := "70.1"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, ThermometerUltraPrecise, res)
}

func TestGetRating_HappyPath_ThermometerVeryPrecise(t *testing.T) {
	expectedType := "thermometer"
	expectedName := "therm-1"
	value1 := "67.5"
	value2 := "72.5"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, ThermometerVeryPrecise, res)
}

func TestGetRating_HappyPath_ThermometerPreciseCorrectAvgButHugeStdDev(t *testing.T) {
	expectedType := "thermometer"
	expectedName := "therm-1"
	value1 := "60"
	value2 := "80"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, ThermometerPrecise, res)
}

func TestGetRating_HappyPath_ThermometerPreciseCorrectStdDevButTotallyOff(t *testing.T) {
	expectedType := "thermometer"
	expectedName := "therm-1"
	value1 := "40.1"
	value2 := "39.9"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, ThermometerPrecise, res)
}

func TestGetRating_HappyPath_HumidityAccepted(t *testing.T) {
	expectedType := "humidity"
	expectedName := "hum-1"
	value1 := "45.2"
	value2 := "44.7"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, HumidityAccepted, res)
}

func TestGetRating_HappyPath_HumidityRejected(t *testing.T) {
	expectedType := "humidity"
	expectedName := "hum-1"
	value1 := "45.2"
	value2 := "44.1"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, HumidityRejected, res)
}

func TestGetRating_WrongSensorType(t *testing.T) {
	expectedType := "Potato"
	expectedName := "Potato-sensor"
	value1 := "70.0"
	value2 := "70.0"
	refValues := NewRefTemperatureHumidity(70.0, 45.0)
	
	sensor := NewSensor(expectedType, expectedName)
	lineData1 := []string{"2000-01-01T00:00:00", expectedName, value1}
	lineData2 := []string{"2000-01-01T00:00:00", expectedName, value2}
	sensor.AppendData(lineData1)
	sensor.AppendData(lineData2)

	res := sensor.GetRating(refValues)

	assert.NotNil(t, res)
	assert.Equal(t, "Invalid sensor type for type " + expectedType + ". Checking next sensor", res)
}
