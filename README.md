<h1 align="center">Snaptube Clone 🚀</h1>

<p align="center">
  <i>A fast, universal, adaptive media downloader powered by Python!</i>
</p>

---

## 📖 Overview

**Snaptube Clone** is a powerful command-line tool that lets you download media from anywhere on the internet. 

Instead of relying on heavy third-party tools that get easily blocked, this project uses a single, highly optimized **Python script** (`downloader.py`). It automatically detects your native screen resolution and perfectly fetches the highest-quality version of the media that fits your display—whether it's a YouTube video, a direct image link, or anything else!

---

## ✨ Features

- **💡 Universal Downloader:** A single Python script handles ANY media URL (YouTube, images, PDFs, archives).
- **📺 Adaptive Resolution:** Automatically detects your display height (e.g., 1080p, 1440p, 4k) and guarantees the downloaded media perfectly matches your screen!
- **▶️ Ultimate YouTube Bypass:** Uses a specialized Android VR spoofing technique inside Python to completely bypass YouTube's 403 and 429 bot blockers. No more errors!
- **📁 Smart File Sorting:** Automatically organizes your downloaded files into a `downloads/mp4/` directory with clean, resolution-tagged filenames.
- **🛑 Safe & Clean:** Cleans up all temporary video and audio tracks automatically after securely merging them with FFmpeg.

---

## ⚙️ Installation

1. **Install System Dependencies:**
   Make sure you have **FFmpeg** and **yt-dlp** installed on your system.
   - On Windows, you can install them via Chocolatey: `choco install ffmpeg yt-dlp`

2. **Install Python Dependencies:**
   Install the required libraries (specifically `pytubefix` for our YouTube blocker-bypassing engine):
   ```bash
   pip install -r requirements.txt
   ```

3. **Download the Project:**
   ```bash
   git clone https://github.com/IamGeniusORG/Snaptube-Clone.git
   cd Snaptube-Clone
   ```

---

## 🚀 How to Use

To download absolutely any media from the internet, simply run the Python script and pass the URL in quotes:

**Download a YouTube Video (Auto-Adapts to your Screen):**
```bash
python downloader.py "https://youtube.com/watch?v=example"
```
*(The script will read your monitor's resolution, bypass YouTube's blockers using Python, and grab the highest quality possible up to your native display limits!)*

**Download from Twitter, Reddit, or other streaming sites:**
```bash
python downloader.py "https://twitter.com/..."
```
*(The Python script will automatically detect that it's not a YouTube link and perfectly hand the URL over to `yt-dlp` to extract the video!)*

### 📁 Where do the files go?
Your perfectly downloaded media will be automatically saved in the auto-generated `downloads/mp4/` folder (or `downloads/others/` for non-YouTube files) with the resolution clearly labeled in the filename!

---

## 📚 Technical Notes
- **Hybrid System:** This project uses a hybrid downloading system. Because YouTube heavily flags `yt-dlp` IP addresses, all YouTube links are securely routed through a custom Python `pytubefix` engine that mimics an Android VR headset to completely bypass 403 and 429 blockers. For all other websites (Twitter, Reddit, Vimeo), the Python script automatically routes the traffic back to `yt-dlp` since they don't have the same strict IP bans!
- **Go Batch Downloader (Optional):** We have left the original `main.go` file in the repo in case you ever want to run mass-concurrent batch downloads for direct links. You can run it via `go run main.go -urls "link1,link2"`.
