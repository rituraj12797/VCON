package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
)

// job represent a single piece of work for a worker goroutine
type job struct {
	index int
	text  string
}

// result => holds the outcome of a job.
type result struct {
	index int
	hash  string
}

// HashStringsInParallel calculates the sha256 hash for an array of strings using true parallelism ( multi processing ).
func Hasher(lines []string) []string {
	numLines := len(lines)

	if numLines == 0 {
		return []string{}
	}

	// we find the number of cpu cores and create that much number of worker go rouines to achieve true parallelism
	numWorkers := runtime.NumCPU()

	// job channel to allow go routines to take work from it
	jobs := make(chan job, numLines)
	// result go routine to alow workers to deposit their results here
	results := make(chan result, numLines)

	// start the go routines and be assured they wont run untill somethign is nto in the channel
	for w := 0; w < numWorkers; w++ {
		go worker(jobs, results)
	}

	// put work in job channel and workers will pick them up and execute parallely
	for i, line := range lines {
		jobs <- job{index: i, text: line}
	}

	// work distributio done so close the channel to avoid any king of dead lock
	close(jobs)

	// collect the result from the result chanel where go routines are depositing their completed work

	// result arrat of same size to show hashes of statement at corrosponding indices
	hashes := make([]string, numLines)
	for i := 0; i < numLines; i++ {
		res := <-results
		hashes[res.index] = res.hash
	}

	return hashes
}

// it reads from the jobs channel performs the hashing and writes to the results channel.
func worker(jobs <-chan job, results chan<- result) {
	// will terminate  once go routine becomes empty as the channel is blocked now
	for j := range jobs {

		h := sha256.New()
		h.Write([]byte(j.text))
		hash := hex.EncodeToString(h.Sum(nil))

		// send the result
		results <- result{index: j.index, hash: hash}
	}
}
