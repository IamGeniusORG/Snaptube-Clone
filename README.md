<h1 align="center">Snaptube Clone 🚀</h1>

<p align="center">
  <strong>Next-Generation Universal Media Extraction Engine</strong><br>
  <i>Powered by Python, FFmpeg, and Adaptive Display Intelligence.</i>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Python-3.14-blue.svg?style=flat-square&logo=python" alt="Python">
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
- [📂 File Sorting](#-file-sorting)
- [📚 Technical Architecture](#-technical-architecture)

---

## 📖 Overview

**Snaptube Clone** is a next-generation command-line tool built to autonomously extract and standardize media from anywhere on the internet. 

Moving beyond legacy downloaders that get IP-banned, this engine uses a cutting-edge **Hybrid System**. It automatically analyzes your hardware's native screen resolution and surgically extracts the absolute highest-quality stream available. From there, it engages a localized FFmpeg pipeline to aggressively upscale, format, and re-encode every frame to perfection. 

---

## ✨ Features

- **💡 Universal Extraction Core:** A single, intelligent Python interface (`downloader.py`) natively handles YouTube, Twitter, Reddit, standard file links, and virtually any media platform.
- **🎵 Audio-Only Mode:** Easily rip pristine MP3s from any video (music, podcasts, lectures) without wasting storage on the video track.
- **📺 Hardware-Adaptive Resolution:** Utilizing cross-platform GUI detection, the engine reads your exact display height (e.g., 1080p, 1440p) and guarantees the downloaded media flawlessly matches your monitor's limits.
- **▶️ Android VR Spoofing:** To bypass YouTube's aggressive 403/429 bot blockers, the engine securely routes connections through a simulated Android VR environment via `pytubefix`.
- **🎥 Studio-Grade Re-encoding:** Every video downloaded is aggressively passed through an FFmpeg normalization pipeline. We enforce a **perfect 16:9 widescreen ratio**, lock the framerate to a buttery smooth **30fps**, and lock the bitrate at **10,000 kbps (10M)**.
- **📁 Autonomous Sorting:** Drops all final files into a unified `downloads/mp4/` directory, while securely wiping raw files from hidden temporary folders.

---

## ⚙️ Installation

### Windows
1. **Install System Dependencies:**
   Ensure you have **FFmpeg** and **yt-dlp** installed. You can quickly install them via an admin PowerShell:
   ```powershell
   choco install ffmpeg yt-dlp
   ```
2. **Install Python Libraries:**
   ```bash
   pip install -r requirements.txt
   ```

### macOS
1. **Install System Dependencies:**
   Use Homebrew to install the required extraction tools:
   ```bash
   brew install ffmpeg yt-dlp
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

To download absolutely any media, invoke the engine and pass your URL in quotes:

**Download a YouTube Video:**
```bash
# Windows
python downloader.py "https://youtube.com/watch?v=example"

# macOS
python3 downloader.py "https://youtube.com/watch?v=example"
```

**Download from Social Media (Twitter, Reddit, etc.):**
```bash
# Windows
python downloader.py "https://twitter.com/..."

# macOS
python3 downloader.py "https://twitter.com/..."
```
*(The engine will instantly detect non-YouTube platforms and hand the URL to the yt-dlp extraction core before applying the FFmpeg studio formatting!)*

**Audio-Only Mode (Rip MP3s):**
To instantly download and extract just the audio from ANY video (YouTube, Reddit, etc.), simply add `--audio` to the end of your command:
```bash
python downloader.py "URL_HERE" --audio
```
*(The engine will download the high-quality audio track and forcefully convert it to a pristine MP3 file via FFmpeg!)*

---

## 📂 File Sorting

Your standardized media will be autonomously compiled and saved in the generated `downloads/mp4/` folder. The engine natively labels the exact resolution directly in the output filename (e.g., `MyVideo [1920x1080].mp4`), so you know you are getting pristine quality.

---

## 📚 Technical Architecture
- **Hybrid Network Design:** YouTube heavily flags raw `yt-dlp` connections. To solve this, all YouTube links are routed through our Python `pytubefix` engine. All other domains (Twitter, Reddit, Vimeo) are automatically routed back to the `yt-dlp` extractor since they lack YouTube's strict IP bans.
- **Safe macOS Parsing:** The script actively checks for the presence of macOS UI toolkits (`python-tk`). If missing, it dynamically bypasses GUI detection and defaults to standard 1080p, preventing any traceback crashes.
- **Go Batch Downloader (Legacy):** The original `main.go` file remains in the repository architecture for high-concurrency batch processing of direct file links. Execute via `go run main.go -urls "link1,link2"`.
