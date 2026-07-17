package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Job represents a single media download task.
type Job struct {
	ID       int
	URL      string
	DestPath string
}

// Result represents the outcome of a job.
type Result struct {
	JobID int
	URL   string
	Err   error
}

// ProgressWriter tracks download completion.
type ProgressWriter struct {
	Total      uint64
	Downloaded uint64
	JobID      int
	FileName   string
	lastReport uint64
}

// Write captures the stream and prints progress.
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += uint64(n)

	percentage := float64(pw.Downloaded) / float64(pw.Total) * 100

	if pw.Downloaded-pw.lastReport >= pw.Total/10 || pw.Downloaded == pw.Total {
		fmt.Printf("   [Job %d] %s: %.0f%% complete\n", pw.JobID, pw.FileName, percentage)
		pw.lastReport = pw.Downloaded
	}
	return n, nil
}

// worker processes downloads and returns results to the resultsChan.
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

		fmt.Printf("-> Worker %d starting Job %d\n", id, job.ID)
		err := downloadFile(ctx, job)
		results <- Result{JobID: job.ID, URL: job.URL, Err: err}
	}
}

func downloadFile(ctx context.Context, job Job) error {
	// FEATURE: Temp Files
	// We download to a .tmp file first. If it gets interrupted, we don't end up with corrupted files.
	tempPath := job.DestPath + ".tmp"
	out, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	
	// Ensure we close the file when the function exits
	defer out.Close()

	// Use http.NewRequestWithContext to make the network request cancellable!
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if resp.ContentLength <= 0 {
		_, err = io.Copy(out, resp.Body)
	} else {
		pw := &ProgressWriter{
			Total:    uint64(resp.ContentLength),
			JobID:    job.ID,
			FileName: filepath.Base(job.DestPath),
		}
		reader := io.TeeReader(resp.Body, pw)
		_, err = io.Copy(out, reader)
	}

	// Close the file explicitly before renaming
	out.Close()

	if err != nil {
		// Cleanup the partial temp file if download failed or was cancelled
		os.Remove(tempPath)
		return err
	}

	// FEATURE: Safe renaming
	// Only when download is 100% successful do we rename it from .tmp to its final .zip name
	return os.Rename(tempPath, job.DestPath)
}

func main() {
	// --- FEATURE 1: Command-Line Flags ---
	workersPtr := flag.Int("workers", 3, "Number of concurrent download workers")
	urlsPtr := flag.String("urls", "", "Comma-separated list of URLs to download")
	flag.Parse()

	fmt.Println("===================================================")
	fmt.Println("🚀 Starting Snaptube Clone v2.0")
	fmt.Printf("⚙️  Concurrency Limit: %d workers\n", *workersPtr)
	fmt.Println("===================================================")

	// Create output folder
	outDir := "downloads"
	os.MkdirAll(outDir, os.ModePerm)

	// Build jobs list from CLI or use default defaults
	var jobURLs []string
	if *urlsPtr != "" {
		jobURLs = strings.Split(*urlsPtr, ",")
	} else {
		jobURLs = []string{
			"http://speedtest.tele2.net/10MB.zip",
			"http://speedtest.tele2.net/10MB.zip",
			"http://speedtest.tele2.net/10MB.zip",
			"http://speedtest.tele2.net/10MB.zip",
			"http://speedtest.tele2.net/10MB.zip",
		}
	}

	// --- FEATURE 2: Graceful Shutdown (Context) ---
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for Ctrl+C (Interrupt signal)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n⚠️  Ctrl+C pressed! Safely cancelling active downloads and cleaning up partial files...")
		cancel() // This sends a cancellation signal to ALL network requests automatically!
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

		// FEATURE: Dynamic File Extensions & Categorization
		// Automatically extract the correct file extension from the URL (e.g. .jpg, .png, .mp4)
		parsedURL, err := url.Parse(cleanURL)
		ext := ".data" // fallback extension if the URL doesn't have one
		if err == nil && path.Ext(parsedURL.Path) != "" {
			ext = strings.ToLower(path.Ext(parsedURL.Path))
		}

		subfolder := "others"
		switch ext {
		case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv":
			subfolder = "videos"
		case ".mp3", ".wav", ".aac", ".flac", ".ogg", ".m4a":
			subfolder = "audio"
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".tiff":
			subfolder = "images"
		case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".csv":
			subfolder = "documents"
		case ".zip", ".rar", ".7z", ".tar", ".gz":
			subfolder = "archives"
		}

		targetDir := filepath.Join(outDir, subfolder)
		os.MkdirAll(targetDir, os.ModePerm)

		jobsChan <- Job{
			ID:       i + 1,
			URL:      cleanURL,
			DestPath: filepath.Join(targetDir, fmt.Sprintf("download_%d%s", i+1, ext)),
		}
	}
	close(jobsChan)

	// Run WaitGroup Wait in background so main thread can process results
	go func() {
		wg.Wait()
		close(resultsChan) // Close results once all workers clock out
	}()

	// --- FEATURE 3: Fan-In Results Aggregation ---
	successCount := 0
	failCount := 0

	// This loop blocks and reads results as they come in from the workers
	for res := range resultsChan {
		if res.Err != nil {
			failCount++
			fmt.Printf("❌ Job %d Failed: %v\n", res.JobID, res.Err)
		} else {
			successCount++
			fmt.Printf("✅ Job %d Completed Successfully\n", res.JobID)
		}
	}

	// Print Final Report
	fmt.Println("===================================================")
	fmt.Printf("📊 FINAL SUMMARY: %d Successful | %d Failed\n", successCount, failCount)
	fmt.Println("===================================================")
}
