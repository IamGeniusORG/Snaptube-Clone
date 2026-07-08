<h1 align="center">Snaptube Clone 🚀</h1>

<p align="center">
  <i>A production-grade, highly parallel mass-downloader built from scratch in Go.</i>
</p>

---

## 📖 Overview

If you try to download 500 files at once in a standard web browser, your computer will crawl to a halt and your network router might crash. **High-Concurrency Throttled Media Processor** solves this problem. 

Built as an educational deep-dive into Go's concurrency primitives, this command-line tool fetches raw web media (images, PDFs, datasets, videos) across multiple concurrent worker threads while managing strict hardware limits (throttling). 

This repository also includes a custom environment for **`yt-dlp`**, perfectly configured to handle complex, protected streaming sites (like YouTube) that require JavaScript decryption.

---

## ✨ Features

- **🚦 Worker Pool Throttling:** Strict concurrency caps prevent you from DDOSing servers or overwhelming your local network adapter.
- **🧠 Dynamic File Extensions:** The URL parser automatically analyzes links and appends the correct extension (`.jpg`, `.png`, `.mp4`, etc.) to the final downloaded file.
- **🛑 Graceful Shutdowns (Ctrl+C Safety):** Uses Go's `context` package. If you abort the program, it instantly halts all active HTTP connections and cleans up junk `.tmp` files.
- **📊 Live Terminal Progress:** Uses `io.TeeReader` to intercept the byte stream and print real-time percentage updates to the console.
- **📈 Fan-Out/Fan-In Aggregation:** Uses a bidirectional channel architecture to aggregate success/failure logs and print a final status report table.

---

## ⚙️ Installation

Please see the [requirements.txt](requirements.txt) file for the exact copy/paste terminal commands to install **Go**, **yt-dlp**, and **FFmpeg** on Windows, macOS, or Linux.

Once installed, clone this repository:
```bash
git clone https://github.com/IamGeniusORG/High-Concurrency-Throttled-Media-Processor.git
cd High-Concurrency-Throttled-Media-Processor
```

---

## 🚀 Usage Guide: The Go Downloader
*Use this for direct, public file links (Datasets, Images, Documents, direct MP4s).*

You can run the Go application directly from the command line using interactive flags:

**1. Basic Usage (Defaults to 3 workers and test files):**
```bash
go run main.go
```

**2. Custom Usage (Set your speed limit and provide your links):**
```bash
go run main.go -workers 5 -urls "http://example.com/data.csv, http://example.com/image.png"
```
*All successful media will be automatically saved into an auto-generated `downloads/` directory.*

---

## 🎥 Usage Guide: `yt-dlp` Integration
*Use this for protected streaming services (YouTube, Twitter, Reddit) where streams are separated and obscured.*

This repository contains a `yt-dlp.conf` file. Whenever you run `yt-dlp` inside this project folder, it automatically applies two custom rules:
1. **Forced IPv4 (`-4`):** Bypasses standard `googlevideo.com` connection timeouts.
2. **Native MP4 Merging:** Forces `FFmpeg` to stitch high-quality audio/video into `.mp4` containers instead of `.webm`.

**Download a Video:**
```bash
yt-dlp "https://youtube.com/watch?v=example"
```

**Override Default to extract only MP3 Audio:**
```bash
yt-dlp -x --audio-format mp3 "https://youtube.com/watch?v=example"
```

---

## 🧠 Under the Hood (For Go Developers)

This tool was built to showcase real-world Go concurrency patterns:
- **`make(chan Job)`**: Acts as the conveyor belt. The main thread loads URLs onto the belt, while worker goroutines safely pull them off one by one.
- **`sync.WaitGroup`**: Ensures the main program patiently waits until every single worker has safely finished writing their data to disk before exiting.
- **`context.WithCancel`**: Tied to `os.Interrupt`, allowing the application to gracefully collapse and run defer cleanups instead of violently crashing.
