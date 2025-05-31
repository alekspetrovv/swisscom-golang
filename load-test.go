package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const outputFileName = "request_results.csv"

type RequestResult struct {
	OverallRequestNum int
	Step              int
	RequestIDInStep   int
	URL               string
	Success           bool
	HTTPStatusCode    int
	ErrorMsg          string
	Duration          time.Duration
}

func sendHTTPRequest(
	targetURL string,
	requestIDInStep int,
	overallRequestNum int,
	currentStep int,
	wg *sync.WaitGroup,
	resultsChan chan<- RequestResult,
) {
	defer wg.Done()

	startTime := time.Now()
	var success bool
	var statusCode int
	var errMsg string

	resp, err := http.Get(targetURL)
	if err != nil {
		errMsg = err.Error()
		success = false
	} else {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
		if statusCode >= 200 && statusCode < 300 {
			success = true
		} else {
			success = false
			errMsg = fmt.Sprintf("HTTP status %s", resp.Status)
		}

		_, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			if errMsg != "" {
				errMsg += "; "
			}
			errMsg += "Error reading body: " + readErr.Error()
			success = false
		}
	}

	duration := time.Since(startTime)
	resultsChan <- RequestResult{
		OverallRequestNum: overallRequestNum,
		Step:              currentStep,
		RequestIDInStep:   requestIDInStep,
		URL:               targetURL,
		Success:           success,
		HTTPStatusCode:    statusCode,
		ErrorMsg:          errMsg,
		Duration:          duration,
	}
}

func main() {
	parallelFlag := flag.Int("parallel", 100, "Number of parallel requests per step")
	stepsFlag := flag.Int("steps", 10, "Number of steps (batches) to run")
	urlFlag := flag.String("url", "http://localhost:4005/api/services", "Target URL for HTTP requests (optional)")

	flag.Parse()

	parallelRequestsPerStep := *parallelFlag
	numberOfSteps := *stepsFlag
	targetURL := *urlFlag
	overallStartTime := time.Now()
	totalRequestsToLaunch := parallelRequestsPerStep * numberOfSteps

	if totalRequestsToLaunch == 0 {
		fmt.Println("No requests to launch (parallel or steps is zero). Exiting.")
		return
	}

	resultsChan := make(chan RequestResult, totalRequestsToLaunch)
	allResults := make([]RequestResult, 0, totalRequestsToLaunch)
	totalRequestsLaunched := 0

	for step := 1; step <= numberOfSteps; step++ {
		var wg sync.WaitGroup

		for i := 1; i <= parallelRequestsPerStep; i++ {
			totalRequestsLaunched++
			wg.Add(1)
			go sendHTTPRequest(targetURL, i, totalRequestsLaunched, step, &wg, resultsChan)
		}

		wg.Wait()
	}

	for i := 0; i < totalRequestsToLaunch; i++ {
		result := <-resultsChan
		allResults = append(allResults, result)
	}
	close(resultsChan)

	overallDuration := time.Since(overallStartTime)

	var successfulRequests, failedRequests int
	for _, res := range allResults {
		if res.Success {
			successfulRequests++
		} else {
			failedRequests++
		}
	}

	fmt.Println("\n--- Overall Summary ---")
	fmt.Printf("Total Requests Launched: %d\n", totalRequestsToLaunch)
	fmt.Printf("Successful Requests:     %d\n", successfulRequests)
	fmt.Printf("Failed Requests:         %d\n", failedRequests)
	fmt.Printf("Total Time Taken:        %s\n", overallDuration)
	fmt.Printf("Detailed per-request logs saved to: %s\n", outputFileName)

	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"OverallRequestNum", "Step", "RequestIDInStep", "URL", "Success", "HTTPStatusCode", "ErrorMsg", "DurationMillis"}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		return
	}

	for _, res := range allResults {
		row := []string{
			strconv.Itoa(res.OverallRequestNum),
			strconv.Itoa(res.Step),
			strconv.Itoa(res.RequestIDInStep),
			res.URL,
			strconv.FormatBool(res.Success),
			strconv.Itoa(res.HTTPStatusCode),
			res.ErrorMsg,
			strconv.FormatInt(res.Duration.Milliseconds(), 10),
		}
		if err := writer.Write(row); err != nil {
			fmt.Printf("Error writing CSV row for request %d: %v\n", res.OverallRequestNum, err)
		}
	}
	fmt.Printf("Results export to %s complete.\n", outputFileName)
}
