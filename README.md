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

## Running the tool

To run the analyzer, either use

```shell
go build && go run .
```

or

```shell
./sensor
```

Type the values directly in the console, or copy/paste the log data.

**Note:** log data must comply to the format given, else errors will be thrown

When you're done typing, end the log capture using `Ctrl+]`

## Testing the tool

Simply run

```shell
go test
```
