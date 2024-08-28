package main

import (
	"encoding/csv"
	"strings"
	"strconv"
	"errors"
	"fmt"
	"log"
	"os"
)

type TemperatureScale struct {
	name string
	abrv string
	meltingPoint float64
	boilingPoint float64
}

func loadTemperatures() []TemperatureScale {
	// Read the CSV file
	file, err := os.Open("scales.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	scales := rawCSVData[1:]
	// Parse the CSV data
	var temperatureScales []TemperatureScale
	for _, record := range scales {
		// fmt.Println(record)
		fields := strings.Split(record[0], ";")
		meltingPoint, _ := strconv.ParseFloat(fields[2], 64)
		boilingPoint, _ := strconv.ParseFloat(fields[3], 64)

		temperatureScale := TemperatureScale{
			name: fields[0],
			abrv: fields[1],
			meltingPoint: meltingPoint,
			boilingPoint: boilingPoint,
		}
		temperatureScales = append(temperatureScales, temperatureScale)
	}
	return temperatureScales
}

func SelectTemperatureScale(
	scales []TemperatureScale, context string, preselected *TemperatureScale,
) (*TemperatureScale, error) {
	// Read from std input
	fmt.Println(context, ": Select a temperature scale")
	for i, scale := range scales {
		if preselected != nil && scale == *preselected { continue }
		fmt.Printf("%d: %s (º%s)\n", i, scale.name, scale.abrv)
	}
	var selectedScale string
	fmt.Scanln(&selectedScale)
	index, _ := strconv.Atoi(selectedScale)
	if index < 0 || index >= len(scales) {
		return nil, errors.New("Invalid selection")
	} else {
		return &scales[index], nil
	}
}

func ConvertTemperature(
	fromScale *TemperatureScale, toScale *TemperatureScale,
) float64 {
	var temperature float64
	fmt.Println("Enter the temperature in", fromScale.abrv)
	fmt.Scanln(&temperature)
	answer := toScale.meltingPoint
	return answer + ((temperature - fromScale.meltingPoint) * (toScale.boilingPoint - toScale.meltingPoint) / (fromScale.boilingPoint - fromScale.meltingPoint))
	// return (
	// 	(temperature - fromScale.meltingPoint)
	// 	* (toScale.boilingPoint - toScale.meltingPoint)
	// 	/ (fromScale.boilingPoint - fromScale.meltingPoint)
	// ) + toScale.meltingPoint
}

func main() {
	temperatureScales := loadTemperatures()
	for {
		fromScale, fromErr := SelectTemperatureScale(temperatureScales, "From", nil)
		if fromErr != nil {
			fmt.Println(fromErr)
			break
		}
		toScale, toErr := SelectTemperatureScale(temperatureScales, "To", fromScale)
		if toErr != nil {
			fmt.Println(toErr)
			break
		}
		if fromScale == toScale {
			fmt.Println("The scales are the same\n")
			continue
		}
		fmt.Printf("Converting from %s to %s\n\n", fromScale.name, toScale.name)
		convertedTemperature := ConvertTemperature(fromScale, toScale)
		fmt.Printf("The converted temperature is: %.2f\n\n", convertedTemperature)
	}
}
