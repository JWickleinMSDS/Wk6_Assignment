package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat"
)

// Grab the CSV file
func LoadCSV(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data [][]float64
	for i, record := range rawCSVData {
		if i == 0 {
			continue
		}
		var floatRow []float64
		for j, value := range record {
			if j == 0 {
				continue
			}
			floatVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			floatRow = append(floatRow, floatVal)
		}
		data = append(data, floatRow)
	}
	return data, nil
}

// Build the simple linear regression.
func SimpleLinearRegression(x, y []float64) (mse, aic, bic float64) {
	// Fit the model
	alpha, beta := stat.LinearRegression(x, y, nil, false)

	// Calculate residuals and RSS
	var rss float64
	for i, xi := range x {
		pi := alpha + beta*xi
		residual := y[i] - pi
		rss += residual * residual
	}
	n := float64(len(x))
	mse = rss / n

	// Calculate AIC and BIC
	l := n * math.Log(rss/n)
	k := 2.0
	aic = l + 2*k
	bic = l + math.Log(n)*k

	return mse, aic, bic
}

func main() {
	// Start the timer
	startTime := time.Now()

	// Load the dataset from the CSV file
	data, err := LoadCSV("boston.csv")
	if err != nil {
		log.Fatalf("Error loading CSV: %v", err)
	}

	x := make([]float64, len(data))
	y := make([]float64, len(data))
	for i, row := range data {
		x[i] = row[0]
		y[i] = row[len(row)-1] // The last column is the dependent variable
	}

	// Calculate the variables
	var totalMSE, totalAIC, totalBIC float64

	// Perform the operation 100 times
	for i := 0; i < 100; i++ {
		// Fit the model and calculate criteria for each iteration
		mse, aic, bic := SimpleLinearRegression(x, y)

		// Accumulate results
		totalMSE += mse
		totalAIC += aic
		totalBIC += bic
	}

	// Calculate averages
	avgMSE := totalMSE / 100
	avgAIC := totalAIC / 100
	avgBIC := totalBIC / 100

	// Stop the timer
	duration := time.Since(startTime)

	// Print average results and time taken
	fmt.Printf("Average MSE: %f, Average AIC: %f, Average BIC: %f\n", avgMSE, avgAIC, avgBIC)
	fmt.Printf("Time taken: %s\n", duration)
}
