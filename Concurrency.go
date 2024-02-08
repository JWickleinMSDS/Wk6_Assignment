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

// Do a simple linear regression and calculate MSE, AIC, and BIC.
func SimpleLinearRegression(x, y []float64) (mse, aic, bic float64) {
	// Fit it
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
	startTime := time.Now()

	data, err := LoadCSV("boston.csv")
	if err != nil {
		log.Fatalf("Error loading CSV: %v", err)
	}

	x := make([]float64, len(data))
	y := make([]float64, len(data))
	for i, row := range data {
		x[i] = row[0]
		y[i] = row[len(row)-1] // THe last column is the dependent variable.
	}

	results := make(chan [3]float64, 100)

	for i := 0; i < 100; i++ {
		go func() {
			mse, aic, bic := SimpleLinearRegression(x, y)
			results <- [3]float64{mse, aic, bic}
		}()
	}

	var totalMSE, totalAIC, totalBIC float64
	for i := 0; i < 100; i++ {
		result := <-results
		totalMSE += result[0]
		totalAIC += result[1]
		totalBIC += result[2]
	}

	avgMSE := totalMSE / 100
	avgAIC := totalAIC / 100
	avgBIC := totalBIC / 100

	duration := time.Since(startTime)
	fmt.Printf("Average MSE: %f, Average AIC: %f, Average BIC: %f\n", avgMSE, avgAIC, avgBIC)
	fmt.Printf("Time taken: %s\n", duration)
}
