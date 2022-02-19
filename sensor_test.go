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