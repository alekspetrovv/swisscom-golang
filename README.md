# Go HTTP Batch Client

A simple command-line HTTP client written in Go that sends GET requests in configurable parallel batches and steps. It measures request times, logs success/failure status, and exports detailed results to a CSV file.

## Overview

This tool is designed to send a series of HTTP GET requests to a specified URL. It allows you to control the concurrency by defining:
* The number of requests to send in parallel within each "step" (batch).
* The total number of "steps" to execute.

It's useful for basic load testing, checking endpoint responsiveness, or any scenario where you need to make multiple HTTP requests in a structured way.

## Features

* Sends HTTP GET requests.
* Configurable number of parallel requests per step.
* Configurable number of steps (batches).
* Configurable target URL.
* Default target Spring Boot API URL: `http://localhost:4005/api/services` (intended to retrieve all services).
* Measures and reports the duration for each individual request.
* Provides a summary on the console including:
    * Total requests launched.
    * Number of successful requests.
    * Number of failed requests.
    * Total time taken for all operations.
* Exports detailed results for every request to a CSV file (`request_results.csv`), including URL, success status, HTTP status code, error messages, and duration.

## Prerequisites

* Go (version 1.20 or later recommended) installed and configured in your system's PATH. You can download Go from [golang.org](https://golang.org/).

## Setup

1.  Clone this repository or download the `http_client.go` and `run_client.sh` files into a directory on your local machine.
    ```bash
    # Example if you have a git repository:
    # git clone <your-go-client-repo-url>
    # cd <your-go-client-project-directory>
    ```
2.  Ensure the Bash script is executable (if you plan to use it):
    ```bash
    chmod +x run_client.sh
    ```

## Important Note: Backend API Setup

**The default target URL for this Go client is `http://localhost:4005/api/services`.**

This URL is intended to connect to the backend Spring Boot API component of the "Swisscom Service Management" application. Before running this Go client with its default URL, please ensure that the corresponding Spring Boot backend application is:
1.  Cloned from its repository.
2.  Set up according to its own `README.md` instructions.
3.  Running and accessible at `http://localhost:4005` (or the port your backend is configured to use).

You can find the backend project and its setup instructions here:
[https://github.com/alekspetrovv/swisscom-java](https://github.com/alekspetrovv/swisscom-java)
## Running the Client

You can run the client using the provided Bash script (recommended for ease of use) or directly with `go run`.

### Using the Bash Script (`run_client.sh`)

The Bash script provides a convenient way to run the Go program with default or custom parameters.

* **Run with default settings:**
  (Defaults: 100 parallel requests per step, 10 steps, URL: `http://localhost:4005/api/services`)
    ```bash
    ./run_client.sh
    ```

* **Run with custom parameters:**
  The script accepts parameters in the following order: `parallel_tasks`, `steps`, `target_url`.
    ```bash
    ./run_client.sh <parallel_tasks> <steps> [target_url]
    ```
    * `<parallel_tasks>`: Number of requests to send in parallel within each step.
    * `<steps>`: Total number of steps to execute.
    * `[target_url]`: (Optional) The full URL to send requests to. If omitted, the default is used.

  **Examples:**
    * Send 50 parallel requests per step, for 3 steps, to the default URL:
        ```bash
        ./run_client.sh 50 3
        ```
    * Send 20 parallel requests per step, for 2 steps, to a custom URL:
        ```bash
        ./run_client.sh 20 2 "[https://jsonplaceholder.typicode.com/posts](https://jsonplaceholder.typicode.com/posts)"
        ```

### Directly with `go run`

You can also run the Go program directly and pass command-line flags to it.

```bash
go run http_client.go -parallel <parallel_tasks> -steps <steps> -url <url> 