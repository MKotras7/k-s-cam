#!/bin/bash

# Directory containing the images
IMAGE_DIR="./manual_captures"

# Output video filename
OUTPUT_VIDEO="./output.mp4"

# Frame rate for the video
FRAME_RATE=30

# Check if the image directory exists
if [ ! -d "$IMAGE_DIR" ]; then
  echo "Directory $IMAGE_DIR does not exist."
  exit 1
fi

# Check if FFmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
  echo "FFmpeg is not installed. Please install it and try again."
  exit 1
fi

# Change to the image directory
cd "$IMAGE_DIR" || exit

# Run FFmpeg to convert images to video
ffmpeg -framerate "$FRAME_RATE" -pattern_type glob -i 'frame_*.jpg' -c:v libx264 -preset slow -crf 22 -pix_fmt yuv420p "./$OUTPUT_VIDEO"

echo "Video has been created: $OUTPUT_VIDEO"