package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadInput_HappyPath(t *testing.T) {
	var stdin bytes.Buffer

	stdin.Write([]byte("reference 70.0 45.0\nthermometer temp-1\n2007-04-05T22:00 temp-1 72.4\n\x1D"))

	res := ReadInput(&stdin)
	assert.Equal(t, 3, len(res))
}

func TestReadInput_NoLines(t *testing.T) {
	var stdin bytes.Buffer

	stdin.Write([]byte("\x1D"))

	res := ReadInput(&stdin)
	assert.Equal(t, 0, len(res))
}

func TestExtractSensorData_HappyPath(t *testing.T) {
	var lines []string

	sensorLine := "thermometer temp-1"
	dataLine1 := "2007-04-05T22:00 temp-1 72.4 "
	dataLine2 := "2007-04-05T22:01 temp-1 76.0"

	lines = append(lines, sensorLine)
	lines = append(lines, dataLine1)
	lines = append(lines, dataLine2)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	
	sensor := res[0]
	expectedType := Thermometer
	expectedName := "temp-1"
	expectedValues := []float64{72.4, 76.0}

	assert.Equal(t, expectedType, sensor.GetType())
	assert.Equal(t, expectedName, sensor.GetName())
	assert.Equal(t, 2, len(sensor.GetValues()))
	assert.Equal(t, expectedValues, sensor.GetValues())
}

func TestExtractSensorData_HappyPathWithIntegers(t *testing.T) {
	var lines []string

	sensorLine := "thermometer temp-1"
	dataLine1 := "2007-04-05T22:00 temp-1 72 "
	dataLine2 := "2007-04-05T22:01 temp-1 76"

	lines = append(lines, sensorLine)
	lines = append(lines, dataLine1)
	lines = append(lines, dataLine2)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	
	sensor := res[0]
	expectedType := Thermometer
	expectedName := "temp-1"
	expectedValues := []float64{72.0, 76.0}

	assert.Equal(t, expectedType, sensor.GetType())
	assert.Equal(t, expectedName, sensor.GetName())
	assert.Equal(t, 2, len(sensor.GetValues()))
	assert.Equal(t, expectedValues, sensor.GetValues())
}

func TestExtractSensorData_LinesInBadOrder(t *testing.T) {
	var lines []string

	sensorLine := "thermometer temp-1"
	dataLine1 := "2007-04-05T22:00 temp-1 72.4 "
	dataLine2 := "2007-04-05T22:01 temp-1 76.0"

	// Here the first line of data is befor the header so shouldn't be counted
	lines = append(lines, dataLine1)
	lines = append(lines, sensorLine)
	lines = append(lines, dataLine2)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	
	sensor := res[0]
	expectedType := Thermometer
	expectedName := "temp-1"
	expectedValues := []float64{76.0}

	assert.Equal(t, expectedType, sensor.GetType())
	assert.Equal(t, expectedName, sensor.GetName())
	assert.Equal(t, 1, len(sensor.GetValues()))
	assert.Equal(t, expectedValues, sensor.GetValues())
}

func TestExtractSensorData_LinesForBadSensor(t *testing.T) {
	var lines []string

	sensorLine := "thermometer temp-1"
	// We inject data from another sensor
	dataLine1 := "2007-04-05T22:00 temp-2 72.4 "
	dataLine2 := "2007-04-05T22:01 temp-1 76.0"

	lines = append(lines, sensorLine)
	lines = append(lines, dataLine1)
	lines = append(lines, dataLine2)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	
	sensor := res[0]
	expectedType := Thermometer
	expectedName := "temp-1"
	expectedValues := []float64{76.0}

	assert.Equal(t, expectedType, sensor.GetType())
	assert.Equal(t, expectedName, sensor.GetName())
	assert.Equal(t, 1, len(sensor.GetValues()))
	assert.Equal(t, expectedValues, sensor.GetValues())
}

func TestExtractSensorData_HappyPathWithTwoSensors(t *testing.T) {
	var lines []string

	sensorLine1 := "thermometer temp-1"
	dataLine1 := "2007-04-05T22:00 temp-1 72.4 "
	dataLine2 := "2007-04-05T22:01 temp-1 76.0"
	// Make sure trailing spaces are trimmed
	sensorLine2 := "humidity hum-1                                      "
	dataLine3 := "2007-04-05T22:04 hum-1 45.2                            "
	dataLine4 := "2007-04-05T22:05 hum-1 45.3"

	lines = append(lines, sensorLine1)
	lines = append(lines, dataLine1)
	lines = append(lines, dataLine2)
	lines = append(lines, sensorLine2)
	lines = append(lines, dataLine3)
	lines = append(lines, dataLine4)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))
	
	sensor1 := res[0]
	expectedType1 := Thermometer
	expectedName1 := "temp-1"
	expectedValues1 := []float64{72.4, 76.0}

	sensor2 := res[1]
	expectedType2 := HumiditySensor
	expectedName2 := "hum-1"
	expectedValues2 := []float64{45.2, 45.3}

	assert.Equal(t, expectedType1, sensor1.GetType())
	assert.Equal(t, expectedName1, sensor1.GetName())
	assert.Equal(t, 2, len(sensor1.GetValues()))
	assert.Equal(t, expectedValues1, sensor1.GetValues())

	assert.Equal(t, expectedType2, sensor2.GetType())
	assert.Equal(t, expectedName2, sensor2.GetName())
	assert.Equal(t, 2, len(sensor2.GetValues()))
	assert.Equal(t, expectedValues2, sensor2.GetValues())
}

func TestExtractSensorData_TwoSensorsNoData(t *testing.T) {
	var lines []string

	sensorLine1 := "thermometer temp-1"
	sensorLine2 := "thermometer temp-2"

	lines = append(lines, sensorLine1)
	lines = append(lines, sensorLine2)

	res := ExtractSensorData(lines)

	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))

	sensor1 := res[0]
	expectedType1 := Thermometer
	expectedName1 := "temp-1"
	var expectedValues1 []float64

	sensor2 := res[1]
	expectedType2 := Thermometer
	expectedName2 := "temp-2"
	var expectedValues2 []float64

	assert.Equal(t, expectedType1, sensor1.GetType())
	assert.Equal(t, expectedName1, sensor1.GetName())
	assert.Equal(t, 0, len(sensor1.GetValues()))
	assert.Equal(t, expectedValues1, sensor1.GetValues())

	assert.Equal(t, expectedType2, sensor2.GetType())
	assert.Equal(t, expectedName2, sensor2.GetName())
	assert.Equal(t, 0, len(sensor2.GetValues()))
	assert.Equal(t, expectedValues2, sensor2.GetValues())
}