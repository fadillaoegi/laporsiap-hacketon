#!/usr/bin/env bash
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
build_dir="$repo_root/tmp/video"
output_dir="$repo_root/output/video"
mkdir -p "$build_dir" "$output_dir/audio" "$output_dir/clean"
swiftc -module-cache-path "$build_dir/swift-module-cache" \
  "$repo_root/scripts/video/render-brand-motion.swift" -o "$build_dir/render-brand-motion"
python3 "$repo_root/scripts/video/synthesize-sound.py" "$build_dir"
for clip in opening closing; do
  # Preserve an existing render; delete/rename explicitly if regenerating motion.
  silent="$output_dir/clean/LaporSiap-$clip-silent.mp4"
  if [[ ! -f "$silent" ]]; then
    "$build_dir/render-brand-motion" "$clip" "$silent" "$repo_root/src/laporsiap_FE/public/fonts"
  fi
  ffmpeg -hide_banner -loglevel warning -y -i "$build_dir/$clip-sound-raw.wav" \
    -af 'loudnorm=I=-18:TP=-2:LRA=5' -ar 48000 -c:a pcm_s16le "$output_dir/audio/LaporSiap-$clip.wav"
  ffmpeg -hide_banner -loglevel warning -y -i "$silent" -i "$output_dir/audio/LaporSiap-$clip.wav" \
    -map 0:v:0 -map 1:a:0 -c:v libx264 -preset slow -crf 18 -pix_fmt yuv420p \
    -c:a aac -b:a 256k -ar 48000 -movflags +faststart -shortest \
    -metadata title="LaporSiap — $clip" -metadata comment="Original 3D motion and synthesized sound" \
    "$output_dir/LaporSiap-$clip.mp4"
  cp "$output_dir/clean/$clip-poster.png" "$output_dir/$clip-poster.png"
  ffprobe -v error -show_entries stream=codec_name,width,height,pix_fmt,r_frame_rate,sample_rate,channels,duration,nb_frames \
    -show_entries format=duration,size -of json "$output_dir/LaporSiap-$clip.mp4" \
    > "$output_dir/$clip-metadata.json"
done
ffmpeg -hide_banner -loglevel warning -y \
  -i "$output_dir/LaporSiap-opening.mp4" -i "$output_dir/LaporSiap-closing.mp4" \
  -filter_complex "[0:v]select='eq(n,27)+eq(n,60)+eq(n,120)+eq(n,195)',scale=480:270,tile=4x1[a];[1:v]select='eq(n,24)+eq(n,48)+eq(n,90)+eq(n,195)',scale=480:270,tile=4x1[b];[a][b]vstack=inputs=2" \
  -frames:v 1 -update 1 -q:v 2 "$output_dir/contact-sheet.jpg"
