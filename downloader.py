import os
import sys
import subprocess
try:
    import tkinter as tk
    TK_AVAILABLE = True
except ImportError:
    TK_AVAILABLE = False
from pytubefix import YouTube

# ---------------------------------------------------------------------------
# Helper: Detect primary monitor height (cross-platform via tkinter)
# ---------------------------------------------------------------------------
def get_screen_height():
    """Return the height (in pixels) of the primary display.
    Falls back to 1080 if detection fails or tkinter is unavailable.
    """
    if not TK_AVAILABLE:
        return 1080
    try:
        root = tk.Tk()
        root.withdraw()
        height = root.winfo_screenheight()
        root.destroy()
        return height
    except Exception:
        return 1080

# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------
args = sys.argv[1:]
audio_only = '--audio' in args
urls = [arg for arg in args if arg != '--audio']

if len(urls) < 1:
    print('Usage: python downloader.py "URL" [--audio]')
    sys.exit(1)

url = urls[0]
print('🔍 Detecting native screen resolution...')
max_height = get_screen_height()
print(f'📺 Native display height detected: {max_height}p')

# ---------------------------------------------------------------------------
# YouTube handling – use pytubefix for fine‑grained control
# ---------------------------------------------------------------------------
if 'youtube.com' in url or 'youtu.be' in url:
    print('🚀 Using pytubefix for YouTube download (adaptive to screen)')
    yt = YouTube(url, client='ANDROID_VR')
    print(f'Title: {yt.title}')

    if audio_only:
        print('🎵 Audio-only mode detected!')
        audio_stream = yt.streams.get_audio_only('mp4')
        print('Downloading audio...')
        audio_file = audio_stream.download(output_path='.', filename='audio_temp.mp4')
        
        safe_title = ''.join(c for c in yt.title if c.isalnum() or c in (' ', '-', '_')).rstrip()
        out_dir = os.path.join(os.getcwd(), 'downloads', 'mp3')
        os.makedirs(out_dir, exist_ok=True)
        final_path = os.path.join(out_dir, f"{safe_title}.mp3")
        
        print('Converting to 320kbps MP3 via ffmpeg...')
        subprocess.run([
            'ffmpeg', '-y', '-i', audio_file,
            '-c:a', 'libmp3lame', '-b:a', '320k',
            final_path
        ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=True)
        
        os.remove(audio_file)
        print(f'✅ Done! Audio saved to: {final_path}')
    else:
        # Choose the best video stream that does not exceed the screen height
        video_stream = None
        for stream in yt.streams.filter(adaptive=True, type='video', file_extension='mp4').order_by('resolution').desc():
            try:
                res = int(stream.resolution.replace('p', ''))
            except Exception:
                continue
            if res <= max_height:
                video_stream = stream
                break
        # Fallback to the lowest available resolution
        if video_stream is None:
            video_stream = yt.streams.filter(adaptive=True, type='video', file_extension='mp4').order_by('resolution').first()

        audio_stream = yt.streams.get_audio_only('mp4')

        # Download temporary files
        print(f'Downloading video ({video_stream.resolution})...')
        video_file = video_stream.download(output_path='.', filename='video_temp.mp4')
        print('Downloading audio...')
        audio_file = audio_stream.download(output_path='.', filename='audio_temp.mp4')

        # Build safe output filename
        safe_title = ''.join(c for c in yt.title if c.isalnum() or c in (' ', '-', '_')).rstrip()
        out_dir = os.path.join(os.getcwd(), 'downloads', 'mp4')
        os.makedirs(out_dir, exist_ok=True)
        final_path = os.path.join(out_dir, f"{safe_title} [{video_stream.resolution}].mp4")

        # Calculate 16:9 target width based on detected screen height
        target_width = int((max_height / 9) * 16)
        vf_scale = f"scale={target_width}:{max_height}:force_original_aspect_ratio=decrease,pad={target_width}:{max_height}:(ow-iw)/2:(oh-ih)/2"

        # Merge video and audio via ffmpeg with strict quality controls (30fps, 10Mbps constant bitrate, forced 16:9)
        print(f'Re-encoding video to perfectly fit {target_width}x{max_height} at 30fps with 10M constant bitrate...')
        subprocess.run([
            'ffmpeg', '-y', '-i', video_file, '-i', audio_file,
            '-c:v', 'libx264',
            '-b:v', '10M',
            '-maxrate', '10M',
            '-bufsize', '20M',
            '-r', '30',
            '-vf', vf_scale,
            '-c:a', 'aac',
            '-b:a', '320k',
            final_path
        ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=True)

        # Cleanup temporary files
        os.remove(video_file)
        os.remove(audio_file)
        print(f'✅ Done! File saved to: {final_path}')
else:
    # ---------------------------------------------------------------------------
    # Non-YouTube URLs - delegate to yt-dlp to handle complex sites (Twitter, Reddit, etc.)
    # ---------------------------------------------------------------------------
    print('🌐 Non-YouTube URL - using yt-dlp to extract and download media')
    
    if audio_only:
        print('🎵 Audio-only mode detected!')
        out_template = os.path.join(os.getcwd(), 'downloads', 'mp3', '%(title)s.%(ext)s')
        cmd = [
            'yt-dlp',
            '-o', out_template,
            '-x', '--audio-format', 'mp3',
            '--audio-quality', '320K',
            url
        ]
        try:
            subprocess.run(cmd, check=True)
            print('✅ Done! Audio extracted and saved to downloads/mp3/')
        except Exception as e:
            print(f"Failed to download audio using yt-dlp: {e}")
    else:
        import glob
        
        target_width = int((max_height / 9) * 16)
        vf_scale = f"scale={target_width}:{max_height}:force_original_aspect_ratio=decrease,pad={target_width}:{max_height}:(ow-iw)/2:(oh-ih)/2"
        
        # We output to a temporary folder first so we can find the file easily
        temp_dir = os.path.join(os.getcwd(), 'downloads', 'temp_dl')
        os.makedirs(temp_dir, exist_ok=True)
        out_template = os.path.join(temp_dir, '%(title)s.%(ext)s')
        
        # Define yt-dlp format respecting the screen height
        fmt = f"bestvideo[height<={max_height}]+bestaudio/best"
        
        cmd = [
            'yt-dlp',
            '-o', out_template,
            '-f', fmt,
            url
        ]
        
        try:
            subprocess.run(cmd, check=True)
            print('✅ yt-dlp successfully downloaded the raw media!')
            
            # Find the downloaded file
            dl_files = glob.glob(os.path.join(temp_dir, '*'))
            if dl_files:
                dl_file = dl_files[0]
                filename = os.path.basename(dl_file)
                name, _ = os.path.splitext(filename)
                
                final_out_dir = os.path.join(os.getcwd(), 'downloads', 'mp4')
                os.makedirs(final_out_dir, exist_ok=True)
                final_path = os.path.join(final_out_dir, f"{name} [{target_width}x{max_height}].mp4")
                
                print(f'Re-encoding video to perfectly fit {target_width}x{max_height} at 30fps with 10M constant bitrate...')
                subprocess.run([
                    'ffmpeg', '-y', '-i', dl_file,
                    '-c:v', 'libx264',
                    '-b:v', '10M',
                    '-maxrate', '10M',
                    '-bufsize', '20M',
                    '-r', '30',
                    '-vf', vf_scale,
                    '-c:a', 'aac',
                    '-b:a', '320k',
                    final_path
                ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=True)
                
                # Cleanup
                os.remove(dl_file)
                print(f'✅ Done! File standardized and saved to: {final_path}')
        except Exception as e:
            print(f"Failed to download or encode using yt-dlp: {e}")
