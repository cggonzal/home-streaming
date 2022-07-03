## Note
Make sure that the video file is a .mp4 and uses H.264 and AAC codecs for video and audio (respectively). If it is not, run the following command on the file:
```
ffmpeg -i input.avi -c:v libx264 -preset slow -crf 20 -c:a aac -b:a 160k -vf format=yuv420p -movflags +faststart output.mp4
```
Make sure to replace <input.avi> with the name of the input file and <output.mp4> with the name of the output file.


# The below instructions are for doing DASH. Browsers will automatically request files in chunks, so DASH is not needed unless you want to dynamically adapt streaming. Which is unecessary for any home project...

### Create different bitrates of a single video file
Input <in.video> below can be .mp4 or .webm. Change name of input and output files as needed.
```
ffmpeg -i in.video -c:v libvpx-vp9 -keyint_min 150 \
-g 150 -tile-columns 4 -frame-parallel 1 -f webm -dash 1 \
-an -vf scale=160:90 -b:v 250k -dash 1 video_160x90_250k.webm \
-an -vf scale=320:180 -b:v 500k -dash 1 video_320x180_500k.webm \
-an -vf scale=640:360 -b:v 750k -dash 1 video_640x360_750k.webm \
-an -vf scale=640:360 -b:v 1000k -dash 1 video_640x360_1000k.webm \
-an -vf scale=1280:720 -b:v 1500k -dash 1 video_1280x720_1500k.webm
```
from step 1 in: https://developer.mozilla.org/en-US/docs/Web/Media/DASH_Adaptive_Streaming_for_HTML_5_Video#using_dash_-_server_side


### Create audio file
Change name of input <in.video> output as needed.
```
ffmpeg -i in.video -vn -acodec libvorbis -ab 128k -dash 1 my_audio.webm

```
from step 1 in: https://developer.mozilla.org/en-US/docs/Web/Media/DASH_Adaptive_Streaming_for_HTML_5_Video#using_dash_-_server_side

### Create the manifest file
Change name of video files and audio files to match the name of the ones that were created in steps above. Also change name of output .mpd manifest file.
```
ffmpeg \
  -f webm_dash_manifest -i video_160x90_250k.webm \
  -f webm_dash_manifest -i video_320x180_500k.webm \
  -f webm_dash_manifest -i video_640x360_750k.webm \
  -f webm_dash_manifest -i video_1280x720_1500k.webm \
  -f webm_dash_manifest -i my_audio.webm \
  -c copy \
  -map 0 -map 1 -map 2 -map 3 -map 4 \
  -f webm_dash_manifest \
  -adaptation_sets "id=0,streams=0,1,2,3 id=1,streams=4" \
  my_video_manifest.mpd

```

from step 2 in: https://developer.mozilla.org/en-US/docs/Web/Media/DASH_Adaptive_Streaming_for_HTML_5_Video#using_dash_-_server_side


### (Optional) Convert .mp4 to .webm
This step is optional and is NOT needed for DASH or any form of streaming to work correctly. It is just here as a reference...
```
ffmpeg -i input.mp4 -c:v libvpx-vp9 -crf 30 -b:v 0 -b:a 128k -c:a libopus output.webm
```
found from: https://stackoverflow.com/questions/47510489/ffmpeg-convert-mp4-to-webm-poor-results