package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

// Job represents a single Python download task.
type Job struct {
	ID    int
	URL   string
	Audio bool
}

// Result represents the outcome of a job.
type Result struct {
	JobID int
	URL   string
	Err   error
}

// worker processes URLs by calling the Python engine and returns results to the resultsChan.
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Check if we've been cancelled (Ctrl+C) before starting a new job
		select {
		case <-ctx.Done():
			results <- Result{JobID: job.ID, URL: job.URL, Err: fmt.Errorf("cancelled by user")}
			continue // Skip remaining jobs
		default:
		}

		fmt.Printf("-> Worker %d starting extraction for: %s\n", id, job.URL)
		
		// Build the python command
		args := []string{"downloader.py", job.URL}
		if job.Audio {
			args = append(args, "--audio")
		}

		pythonBin := "python"
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			pythonBin = "python3"
		}
		cmd := exec.CommandContext(ctx, pythonBin, args...)
		// Force Python to use UTF-8 encoding so emojis don't crash the script when piped through Go on Windows
		cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
		// We can pipe the python script's output to standard output so the user sees progress
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		results <- Result{JobID: job.ID, URL: job.URL, Err: err}
	}
}

func main() {
	// --- Command-Line Flags ---
	workersPtr := flag.Int("workers", 2, "Number of concurrent Python engines to run")
	urlsPtr := flag.String("urls", "", "Comma-separated list of URLs to download")
	audioPtr := flag.Bool("audio", false, "Extract audio only (MP3)")
	flag.Parse()

	fmt.Println("===================================================")
	fmt.Println("🚀 Starting Snaptube Clone Batch Manager v3.0")
	fmt.Printf("⚙️  Concurrency Limit: %d active engines\n", *workersPtr)
	if *audioPtr {
		fmt.Println("🎵 Mode: Audio-Only (MP3)")
	}
	fmt.Println("===================================================")

	if *urlsPtr == "" {
		fmt.Println("❌ Error: No URLs provided. Use -urls \"link1,link2\"")
		os.Exit(1)
	}

	jobURLs := strings.Split(*urlsPtr, ",")

	// --- Graceful Shutdown (Context) ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for Ctrl+C (Interrupt signal)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n⚠️  Ctrl+C pressed! Safely killing active Python engines...")
		cancel() // This instantly kills all running exec.CommandContext processes!
	}()

	// Channels
	jobsChan := make(chan Job, len(jobURLs))
	resultsChan := make(chan Result, len(jobURLs))
	var wg sync.WaitGroup

	// Start Worker Pool
	for w := 1; w <= *workersPtr; w++ {
		wg.Add(1)
		go worker(ctx, w, jobsChan, resultsChan, &wg)
	}

	// Distribute jobs
	for i, rawURL := range jobURLs {
		cleanURL := strings.TrimSpace(rawURL)
		if cleanURL == "" {
			continue
		}
		jobsChan <- Job{
			ID:    i + 1,
			URL:   cleanURL,
			Audio: *audioPtr,
		}
	}
	close(jobsChan)

	// Wait for workers in background
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// --- Fan-In Results Aggregation ---
	successCount := 0
	failCount := 0

	for res := range resultsChan {
		if res.Err != nil {
			failCount++
			fmt.Printf("❌ Job %d Failed: %v\n", res.JobID, res.Err)
		} else {
			successCount++
			fmt.Printf("✅ Job %d Successfully Processed by Python Engine\n", res.JobID)
		}
	}

	// Print Final Report
	fmt.Println("===================================================")
	fmt.Printf("📊 FINAL SUMMARY: %d Successful | %d Failed\n", successCount, failCount)
	fmt.Println("===================================================")
}
