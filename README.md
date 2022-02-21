```
  _______ _    _ ______ _____  __  __  ____         _    _ _    _ __  __ 
 |__   __| |  | |  ____|  __ \|  \/  |/ __ \       | |  | | |  | |  \/  |
    | |  | |__| | |__  | |__) | \  / | |  | |______| |__| | |  | | \  / |
    | |  |  __  |  __| |  _  /| |\/| | |  | |______|  __  | |  | | |\/| |
    | |  | |  | | |____| | \ \| |  | | |__| |      | |  | | |__| | |  | |
    |_|  |_|  |_|______|_|  \_\_|  |_|\____/       |_|  |_|\____/|_|  |_|
```

![workflow](https://github.com/pingu1/thermo-hum/actions/workflows/go.yml/badge.svg)

# THERMO-HUM log analyzer tool

## Context

```
37Widgets makes inexpensive thermometers and humidity sensors. In order to spot check the manufacturing process, some units are put in a test environment (for an unspecified amount of time) and their readings are logged. The test environment has a static, known temperature and relative humidity, but the sensors are expected to fluctuate a bit.                                                                                                          
 
As a developer, your task is to process the logs and automate the quality control evaluation. The evaluation criteria are as follows:
 
1) For a thermometer, it is branded "ultra precise" if the mean of the readings is within 0.5 degrees of the known temperature, and the standard deviation is less than 3. It is branded "very precise" if the mean is within 0.5 degrees of the room, and the standard deviation is under 5. Otherwise, it's sold as "precise".                  
 
2) For a humidity sensor, it must be discarded unless it is within 1% of the reference value for all readings.  
 
An example log looks like the following. The first line means that the room was held at a constant 70 degrees, 45% relative humidity. Subsequent lines either identify a sensor (<type> <name>) or give a reading (<time> <name> <value>).                                                
 
reference 70.0 45.0                                    
thermometer temp-1                                      
2007-04-05T22:00 temp-1 72.4                            
2007-04-05T22:01 temp-1 76.0                            
2007-04-05T22:02 temp-1 79.1                            
2007-04-05T22:03 temp-1 75.6                            
2007-04-05T22:04 temp-1 71.2                            
2007-04-05T22:05 temp-1 71.4                            
2007-04-05T22:06 temp-1 69.2                            
2007-04-05T22:07 temp-1 65.2                            
2007-04-05T22:08 temp-1 62.8                            
2007-04-05T22:09 temp-1 61.4                            
2007-04-05T22:10 temp-1 64.0                            
2007-04-05T22:11 temp-1 67.5                            
2007-04-05T22:12 temp-1 69.4                            
thermometer temp-2                                      
2007-04-05T22:01 temp-2 69.5                            
2007-04-05T22:02 temp-2 70.1                            
2007-04-05T22:03 temp-2 71.3                            
2007-04-05T22:04 temp-2 71.5                            
2007-04-05T22:05 temp-2 69.8                            
humidity hum-1                                          
2007-04-05T22:04 hum-1 45.2                            
2007-04-05T22:05 hum-1 45.3                            
2007-04-05T22:06 hum-1 45.1                            
humidity hum-2                                          
2007-04-05T22:04 hum-2 44.4                            
2007-04-05T22:05 hum-2 43.9                            
2007-04-05T22:06 hum-2 44.9                            
2007-04-05T22:07 hum-2 43.8                            
2007-04-05T22:08 hum-2 42.1                            
 
Output                                                  
temp-1: precise                                        
temp-2: ultra precise                                  
hum-1: OK                                              
hum-2: discard                                          
 
The log should be read from stdin. In the end, you will own this process, so you should solve the problem as described, but feel free to advocate for any changes you think would make sense to improve the process (split into multiple files, change log format, etc).
```

## A few assumptions

* We will assume that the log data is small enough to be injected via the console line and fits in the memory of the computer on which the program runs
* We will assume no web-server is required, else we would probably want to have endpoints accepting some sort of JSON formatted data instead of a raw log file
* We will assume the log data has always the same format and that a block of data is included between 2 sensors definitions
* We will assume that we're testing a small sample of the entire production, hence the standard devidations formula is SD = SQRT(SUM(POW(xi - avg, 2)) / (N-1)) where xi is the data at index i, avg is the average value of all data, and N is the number of points
* We will assume that if the code encounters an error in the data provided, it should discard the line (and possibly log the error to Stderr)
* We will assume that the code should be optimized for speed of execution

## Running the tool

To run the analyzer, either use

```shell
go run .
```

or

```shell
go build && ./sensor
```

Type the values directly in the console, or copy/paste the log data.

**Note:** log data must comply to the format given, else errors will be thrown

When you're done typing, end the log capture using `Ctrl+]`

## Testing the tool

Simply run

```shell
go test
```

Or check GitHub actions.