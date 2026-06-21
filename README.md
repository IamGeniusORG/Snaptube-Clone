# High-Concurrency Throttled Media Processor 🚀

A production-grade, highly parallel mass-downloader built from scratch in Go. This tool is designed to fetch direct web media (images, PDFs, datasets, videos) across multiple concurrent worker threads while managing explicit hardware and network resource bounds (throttling).

This repository also includes a custom-configured `yt-dlp` environment for handling complex, protected streaming sites (like YouTube) that require Javascript decryption and media merging.

---

## 🛠️ Core Features (Go Application)

- **Worker Pool Throttling:** Capped concurrency using Go routines to prevent DDOSing servers or overwhelming local network adapters.
- **Dynamic File Extensions:** The URL parser automatically reads the link and appends the correct extension (`.jpg`, `.png`, `.mp4`, etc.) to the final downloaded file.
- **Graceful Shutdowns (Ctrl+C Safety):** Uses Go's `context` package. If the user cancels the program, it instantly halts all active HTTP connections and cleans up junk `.tmp` files.
- **Live Terminal Progress:** Uses `io.TeeReader` to intercept the byte stream and print real-time percentage updates to the console.
- **Fan-Out/Fan-In Aggregation:** Uses a bidirectional channel architecture to aggregate success/failure logs and print a final status report table.
- **Safe Disk Writes:** Files are downloaded as `.tmp` files and only renamed to their final extension once the HTTP response successfully hits 100%.

## 🚀 How to Use the Go Downloader

You can run the application directly from the command line using flags to set your speed limit and provide your list of URLs.

```bash
# Example: Download 2 files concurrently using 5 workers
go run main.go -workers 5 -urls "http://example.com/data1.csv, http://example.com/image.png"
```

All successful media will be automatically saved into an auto-generated `downloads/` directory.

---

## 🎥 `yt-dlp` Integration (For Protected Streams)

Standard HTTP clients (like our Go program) cannot download from protected streaming services like YouTube or Netflix because the streams are separated and obscured by JavaScript. For these use cases, this project includes a configuration file (`yt-dlp.conf`) designed to be used alongside the industry-standard `yt-dlp` package.

### Configuration Features:
- **Forced IPv4 (`-4`):** Automatically bypasses `googlevideo.com` connection timeouts caused by faulty ISP IPv6 routing.
- **Native MP4 Merging:** Automatically forces `FFmpeg` to stitch the highest quality audio and video streams into a native `.mp4` container instead of `.webm`.

### Usage:
*(Requires `yt-dlp` and `FFmpeg` to be installed on your system)*

```bash
# The configuration file automatically applies the custom rules
yt-dlp "https://youtube.com/watch?v=example"
```

To override the config and download the raw `.webm` stream or extract an `.mp3`:
```bash
yt-dlp --merge-output-format webm "URL"
yt-dlp -x --audio-format mp3 "URL"
```

---

*Built for educational mastery of Go concurrency primitives.*