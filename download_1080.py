import os
import sys
import subprocess
from pytubefix import YouTube
from pytubefix.cli import on_progress

if len(sys.argv) < 2:
    print('Usage: python download_1080.py "YOUTUBE_URL"')
    sys.exit(1)

url = sys.argv[1]
print("Bypassing YouTube's 403 blocks with pytubefix...")

# Initialize YouTube object using ANDROID_VR client to bypass 403s and token prompts
yt = YouTube(url, client='ANDROID_VR')

print(f"Title: {yt.title}")

# Get highest resolution video stream (1080p)
video_stream = yt.streams.filter(adaptive=True, type="video", file_extension="mp4").order_by("resolution").desc().first()
audio_stream = yt.streams.get_audio_only("mp4")

if not video_stream or not audio_stream:
    print("Could not find 1080p video or audio streams!")
    exit(1)

print(f"Downloading Video: {video_stream.resolution}")
video_file = video_stream.download(output_path=".", filename="video_temp.mp4")

print("\nDownloading Audio")
audio_file = audio_stream.download(output_path=".", filename="audio_temp.mp4")

# Output directory
output_dir = os.path.join(os.getcwd(), "downloads", "mp4")
os.makedirs(output_dir, exist_ok=True)

# Safe filename
safe_title = "".join([c for c in yt.title if c.isalpha() or c.isdigit() or c==' ']).rstrip()
final_output = os.path.join(output_dir, f"{safe_title} [1920x1080].mp4")

print(f"\nMerging Video and Audio into 1080p MP4 using FFmpeg...")
try:
    subprocess.run([
        "ffmpeg", "-y", "-i", video_file, "-i", audio_file, 
        "-c:v", "copy", "-c:a", "aac", final_output
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    
    # Cleanup temp files
    os.remove(video_file)
    os.remove(audio_file)
    
    print(f"\nSUCCESS! Flawless 1080p video saved to: {final_output}")
except Exception as e:
    print(f"FFmpeg merging failed: {e}")
