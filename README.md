<h1 align="center">Snaptube Clone 🚀</h1>

<p align="center">
  <i>A fast, easy-to-use mass-downloader for internet media!</i>
</p>

---

## 📖 Overview

**Snaptube Clone** is a powerful command-line tool that lets you download multiple videos, images, or documents at the same time without crashing your internet. 

Whether you want to download a batch of files directly or download videos from sites like YouTube, this tool handles it all quickly and safely. It also features a built-in smart sorter that automatically organizes your downloaded files into perfectly categorized folders!

---

## ✨ Features

- **🚀 Fast Downloads:** Downloads multiple files at the same time.
- **🚦 Safe Limits:** Automatically limits downloads so it doesn't slow down your home network.
- **📁 Smart File Sorting:** Automatically sorts downloaded files into specific folders (`downloads/videos/`, `downloads/audio/`, `downloads/mp4/`, etc.) based on their formats!
- **▶️ Ultimate YouTube Bypass:** Uses a special `yt-dlp.conf` file to pretend to be an Android phone, completely bypassing YouTube's bot blockers and DRM locks. 
- **🛑 Safe Cancel:** If you press Ctrl+C, it stops cleanly without leaving broken/junk files on your computer.
- **📊 Live Progress:** Shows you exactly how much of each file has been downloaded.

---

## ⚙️ Installation

Make sure you have installed the required software listed in `requirements.txt` (Go, yt-dlp, and FFmpeg).

Then, download the project:
```bash
git clone https://github.com/IamGeniusORG/Snaptube-Clone.git
cd Snaptube-Clone
```

---

## 🚀 How to Use (Direct Links)
*Use this for direct links to images, PDFs, or standard video files.*

**1. Basic Usage (Test files):**
```bash
go run main.go
```

**2. Download Your Own Links:**
```bash
go run main.go -workers 5 -urls "http://example.com/image.png, http://example.com/video.mp4"
```
*Your files will be automatically sorted into `downloads/images/`, `downloads/videos/`, etc., based on their extension!*

---

## 🎥 How to Use (YouTube & Streaming Sites)
*Use this for YouTube, Twitter, Reddit, etc.*

Because this project includes a special `yt-dlp.conf` file, you don't need to type any complicated flags. It will automatically bypass YouTube's blockers, merge high-quality audio and video using FFmpeg, and sort the final file into `downloads/mp4/` or `downloads/mp3/`.

**Download a Video (Standard 360p - Highest Stability):**
```bash
yt-dlp "https://youtube.com/watch?v=example"
```

**Download a Video (Flawless 1080p - Bypass Mode):**
If your IP is flagged by YouTube (403/429 errors), use our custom built-in Python script which bypasses the blocks to fetch crystal-clear 1080p video:
```bash
python download_1080.py "https://youtube.com/watch?v=example"
```

**Download Audio Only (MP3):**
```bash
yt-dlp -x --audio-format mp3 "https://youtube.com/watch?v=example"
```
