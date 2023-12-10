package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	fmt.Println("Server listening on :8000...")
	http.ListenAndServe(":8000", nil)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, false)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, true)
}

func handleRequest(w http.ResponseWriter, r *http.Request, concurrent bool) {
	// Decode JSON payload
	var input struct {
		ToSort [][]int `json:"to_sort"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Process and measure time
	startTime := time.Now()
	var sortedArrays [][]int

	if concurrent {
		sortedArrays = processConcurrently(input.ToSort)
	} else {
		sortedArrays = processSequentially(input.ToSort)
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Send response
	response := struct {
		SortedArrays [][]int `json:"sorted_arrays"`
		TimeNs       string  `json:"time_ns"`
	}{
		SortedArrays: sortedArrays,
		TimeNs:       fmt.Sprint(elapsedTime.Nanoseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processSequentially(arrays [][]int) [][]int {
	var result [][]int
	for _, arr := range arrays {
		time.Sleep(1 * time.Second)
		sort.Ints(arr)
		result = append(result, arr)
	}
	return result
}

func processConcurrently(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	length := len(arrays)
	var result = make([][]int, length)

	for index, arr := range arrays {
		wg.Add(1)
		go func(a []int, index int) {
			time.Sleep(1 * time.Second)

			sort.Ints(a)
			result[index] = a

			wg.Done()
		}(arr, index)
	}

	wg.Wait()
	return result
}
