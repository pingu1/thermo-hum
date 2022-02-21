package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/**
 * General Note:
 *   For this code, we assume that logs are farily small and will all fit in memory.
 *   We also assume they're copied and pasted directly in the console as inputs.
 *   Obviously, for more complex logs (files stored on disk, weighting several Gb), other
 *   solutions would fit better (https://pkg.go.dev/io/ioutil#ReadAll or https://pkg.go.dev/bufio)
 */

/**
 * Reding stdin from console line
 */
func ReadInput(stdin io.Reader) []string {
	var lines []string

	scan := bufio.NewScanner(stdin)

	// Scan until break char
	for {

		// Prompt for input
		fmt.Println("Enter log content:")
		for scan.Scan() {

			line := scan.Text()
			if len(line) == 1 {
				// Set break char as Ctrl+]
				if line[0] == '\x1D' {
					return lines
				}
			}
			// aggregate lines in an array
			lines = append(lines, line)
		}

		// In case something goes wrong, stop reading the input
		if err := scan.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return lines
}

func ExtractSensorData(lines []string) []SensorInterface {

	/**
	 * Probably a better approach for this function would be to have a map of all sensors
	 * and adding data as we read lines. However, this means accessing data on each and
	 * every line we get and possibly updating the inmemory store, which is ressource
	 * intensive.
	 * Instead, we'll be using a buffering approach where we assume we read lines until we
	 * find a new sensor, then we add data to this sensor until we find another line with
	 * a new sensor. By buffering the data like this, we make sure we're only writing to
	 * the slice.
	 */

	// Are we expecting a new sensor or some sensor data?
	expectingSensor := true
	expectingData := false

	// Init sensors
	sensors := make([]SensorInterface, 0)
	var currentSensor SensorInterface

	// We arbitrarilly decide that all data from a given sensor are within the sensor's block
	for _, line := range lines {
		data := strings.Split(strings.TrimSpace(line), " ")

		if expectingSensor && len(data) == 3 {
			// We are looking for some new Sensor and we are reading data, skip the line
			continue
		}

		if expectingData && len(data) == 2 {
			// We are getting a new sensor, which means the data from last sensor are all fetched
			// 1. append the sensor data to the returned slice
			sensors = append(sensors, currentSensor)

			// 2. swap flags
			expectingSensor = true
			expectingData = false
		}

		if expectingData {
			// 1. Append values to current sensor, don't care about errors (arbitrary choice)
			err := currentSensor.AppendData(data)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}

		if expectingSensor {
			// 1. Create a new sensor
			currentSensor = NewSensor(data[0], data[1])
			// 2. swap flags
			expectingSensor = false
			expectingData = true
		}
	}

	// Don't forget to add the last sensor to the list
	sensors = append(sensors, currentSensor)

	return sensors
}
