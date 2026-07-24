<h1 align="center">Snaptube Clone 🚀</h1>

<p align="center">
  <strong>Next-Generation Universal Media Extraction & Concurrency Engine</strong><br>
  <i>Powered by Python, FFmpeg, Go, and Adaptive Display Intelligence.</i>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Python-3.14-blue.svg?style=flat-square&logo=python" alt="Python">
  <img src="https://img.shields.io/badge/Go-1.22-cyan.svg?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/FFmpeg-Ready-red.svg?style=flat-square&logo=ffmpeg" alt="FFmpeg">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS-lightgrey.svg?style=flat-square" alt="Platform">
</p>

---

## 📋 Table of Contents
- [📖 Overview](#-overview)
- [✨ Features](#-features)
- [⚙️ Installation](#️-installation)
  - [Windows](#windows)
  - [macOS](#macos)
- [🚀 How to Use](#-how-to-use)
  - [Batch Downloading (Go)](#batch-downloading-go-orchestrator)
  - [Single Downloads (Python)](#single-downloads-python-direct)
- [📂 File Sorting](#-file-sorting)
- [📚 Technical Architecture](#-technical-architecture)

---

## 📖 Overview

**Snaptube Clone v3.0** is a next-generation command-line tool built to autonomously extract and standardize media from anywhere on the internet. 

Moving beyond legacy single-threaded downloaders, this system utilizes a cutting-edge **Hybrid Concurrency Engine**. A high-performance Go orchestrator natively manages a scalable pool of Python extraction bots. It automatically analyzes your hardware, scales video resolution, and surgically extracts the absolute highest-quality stream available before engaging a localized FFmpeg pipeline to aggressively upscale, format, and re-encode every frame to perfection.

---

## ✨ Features

- **⚡ High-Concurrency Batch Orchestrator:** Use the `main.go` engine to securely run up to 15 Python instances simultaneously. The Go script acts as a master traffic controller, safely managing your CPU load without freezing your computer.
- **📺 Hardware-Adaptive Resolution:** Utilizing cross-platform GUI detection, the engine actively reads your exact display height (e.g., 1080p on your current laptop). It guarantees the downloaded media flawlessly synchronizes to your specific monitor's limits. If you run this on a 1440p or 4K MacBook, the script automatically scales and formats the extraction target to perfectly match that device!
- **💡 Universal Extraction Core:** A single, intelligent Python interface natively handles YouTube, Twitter, Reddit, standard file links, and virtually any media platform.
- **🎵 Audio-Only Mode:** Easily rip pristine MP3s from any video (music, podcasts, lectures) without wasting storage on the video track.
- **🛡️ UUID Sandbox Isolation:** Every single active worker generates a cryptographically unique `temp_dl_UUID` directory for raw downloads, ensuring files never clash, even under massive concurrency loads. The root directory remains spotless.
- **🔄 Auto-Fallback Anti-Bot Pipeline:** If YouTube aggressively IP-bans the primary `pytubefix` engine during a massive batch download, the script instantly catches the error and flawlessly routes the URL into a robust `yt-dlp` fallback engine to finish the job.
- **🎥 Studio-Grade Re-encoding:** Every video downloaded is aggressively passed through an optimized FFmpeg normalization pipeline. We enforce a **perfect 16:9 widescreen ratio**, lock the framerate to a buttery smooth **30fps**, and lock the bitrate at **10,000 kbps (10M)** using hardware-safe thread limits (`-threads 2`) and the `ultrafast` algorithm.

---

## ⚙️ Installation

### Windows
1. **Install System Dependencies:**
   Ensure you have **FFmpeg**, **yt-dlp**, and **Go** installed. You can quickly install them via an admin PowerShell:
   ```powershell
   choco install ffmpeg yt-dlp go
   ```
2. **Install Python Libraries:**
   ```bash
   pip install -r requirements.txt
   ```

### macOS
1. **Install System Dependencies:**
   Use Homebrew to install the required extraction tools:
   ```bash
   brew install ffmpeg yt-dlp go
   ```
2. **Install Python Libraries:**
   *(macOS strictly uses python3)*
   ```bash
   pip3 install pytubefix
   ```

3. **Download the Project:**
   ```bash
   git clone https://github.com/IamGeniusORG/Snaptube-Clone.git
   cd Snaptube-Clone
   ```

---

## 🚀 How to Use

### Batch Downloading (Go Orchestrator)
The absolute best way to use Snaptube Clone is through the Go wrapper. The Go script dynamically detects your operating system (natively using `python3` on Mac) and securely handles all worker processes.

**Mass-Download Videos:**
```bash
go run main.go -workers 3 -urls "link1, link2, link3, link4"
```

**Mass-Download Audio Only (MP3):**
*(Audio encoding is incredibly fast, so you can safely crank up the workers!)*
```bash
go run main.go -audio -workers 10 -urls "link1, link2, link3"
```

### Single Downloads (Python Direct)
To test a single link without the Go Orchestrator, pass your URL in quotes directly to the Python engine:
```bash
# Windows
python downloader.py "https://youtube.com/watch?v=example"

# macOS
python3 downloader.py "https://youtube.com/watch?v=example"
```
*(Add `--audio` to the end of the Python command for an MP3 rip).*

---

## 📂 File Sorting

Your standardized media will be autonomously compiled and saved in the generated `downloads/mp4/` or `downloads/mp3/` folders. The engine natively labels the exact resolution directly in the output filename (e.g., `MyVideo [1920x1080].mp4` if downloaded on your current laptop, or `MyVideo [2560x1440].mp4` if downloaded on a 1440p display). Every single device gets its own perfectly synchronized resolution.

---

## 📚 Technical Architecture
- **Automatic Mac Detection:** The Go orchestrator dynamically parses `runtime.GOOS` to natively route commands to `python3` on macOS/Linux while preserving standard `python` routes for Windows.
- **Hybrid Network Design:** YouTube heavily flags raw `yt-dlp` connections. All YouTube links route through our `pytubefix` engine via Android VR spoofing. If blocked by Bot Detection, it gracefully transitions back into the `yt-dlp` engine. 
- **Safe macOS Parsing:** The python script actively checks for the presence of macOS UI toolkits (`python-tk`). If missing, it dynamically bypasses GUI detection and defaults to standard 1080p, preventing any traceback crashes.
