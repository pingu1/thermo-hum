package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	lines := ReadInput(os.Stdin)
	var header string

	/** debugging **/
	// fmt.Printf("Found %d lines\n", len(lines))

	if len(lines) > 1 {
		header, lines = lines[0], lines[1:]
	} else {
		header = lines[0]
		lines = nil
	}

	ref, err := ExtractRef(header)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/** debugging **/
	fmt.Printf("\nRef. Temperature is %f | Ref. Humidity is %f\n\n", ref.GetRefTemperature(), ref.GetRefHumidity())

	if lines == nil {
		fmt.Println("No content found for sensors, exiting now")
		os.Exit(1)
	}

	sensors := ExtractSensorData(lines)

	/** debugging **/
	fmt.Printf("Found %d sensors\n\n", len(sensors))

	/** debugging **/
	/* Testing synchroneously
	   for _, sensor := range sensors {
	       fmt.Printf("Data for sensor %s\n", sensor.GetName())
	       fmt.Println(sensor.GetValues())
	       fmt.Printf("Average value: %f\n", sensor.GetAverageValue())
	       fmt.Printf("Standard deviation: %f\n", sensor.GetStandardDeviation())
	       fmt.Printf("Max deviation: %f%%\n", sensor.GetMaxDeviationPercentage(ref.GetRefHumidity()))
	       fmt.Printf("Rating: %s\n\n", sensor.CalculateRating(ref))
	   }
	*/

	ComputeResults(sensors, ref)

	// Printing results
	for _, sensor := range sensors {
		fmt.Printf("%s: %s\n", sensor.GetName(), sensor.GetRating())
	}
}

func ComputeResults(sensors []SensorInterface, ref ReferenceInterface) {
	// To make sure we're doing this as fast as possible, calculate the ratings in an async manner
	var wg sync.WaitGroup
	for i := 0; i < len(sensors); i++ {
		wg.Add(1)
		sensor := sensors[i]
		go func() {
			defer wg.Done()
			sensor.SetRating(ref)
		}()
	}
	wg.Wait()
}
