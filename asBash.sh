#!/bin/bash

set -euo pipefail

# --- CONFIG INTERACTIVE PROMPTS ---

SOURCE_DIR=$(whiptail --inputbox "Enter the full path to the source directory:" 10 70 "/mnt/data" --title "Source Directory" 3>&1 1>&2 2>&3)
DEST_HOST=$(whiptail --inputbox "Enter the remote rsync destination (e.g. backup@host:/path):" 10 70 "backup@nas:/vault" --title "Remote Destination" 3>&1 1>&2 2>&3)

CHUNK_SIZE_GB=$(whiptail --inputbox "Chunk size (GiB):" 10 60 "1000" --title "Chunk Size" 3>&1 1>&2 2>&3)
MAX_FILES=$(whiptail --inputbox "Maximum files per chunk:" 10 60 "10000" --title "File Count Limit" 3>&1 1>&2 2>&3)
ZSTD_THREADS=$(whiptail --inputbox "Number of threads for zstd compression:" 10 60 "48" --title "Compression Threads" 3>&1 1>&2 2>&3)
ZSTD_LEVEL=$(whiptail --inputbox "zstd compression level (suggested: 7–9):" 10 60 "9" --title "Compression Level" 3>&1 1>&2 2>&3)

CHUNK_DIR="./chunks"
ARCHIVE_DIR="./compressed"
mkdir -p "$CHUNK_DIR" "$ARCHIVE_DIR"

# Convert GB to bytes
MAX_CHUNK_BYTES=$((CHUNK_SIZE_GB * 1024 * 1024 * 1024))
CHUNK_INDEX=0
FILES=()

# --- CHUNKING ---

whiptail --infobox "Scanning and chunking files..." 8 60

while IFS= read -r -d '' file; do
    FILES+=("$file")
done < <(find "$SOURCE_DIR" -type f -print0)

CURRENT_CHUNK_SIZE=0
CURRENT_FILE_COUNT=0
CHUNK_LIST=()

for filepath in "${FILES[@]}"; do
    filesize=$(stat -c%s "$filepath")
    CURRENT_CHUNK_SIZE=$((CURRENT_CHUNK_SIZE + filesize))
    CURRENT_FILE_COUNT=$((CURRENT_FILE_COUNT + 1))
    CHUNK_LIST+=("$filepath")

    if (( CURRENT_CHUNK_SIZE >= MAX_CHUNK_BYTES || CURRENT_FILE_COUNT >= MAX_FILES )); then
        chunk_file="$CHUNK_DIR/chunk_$CHUNK_INDEX.list"
        printf "%s\n" "${CHUNK_LIST[@]}" > "$chunk_file"
        CHUNK_INDEX=$((CHUNK_INDEX + 1))
        CURRENT_CHUNK_SIZE=0
        CURRENT_FILE_COUNT=0
        CHUNK_LIST=()
    fi
done

# Final flush
if (( ${#CHUNK_LIST[@]} > 0 )); then
    chunk_file="$CHUNK_DIR/chunk_$CHUNK_INDEX.list"
    printf "%s\n" "${CHUNK_LIST[@]}" > "$chunk_file"
fi

# --- COMPRESSION LOOP ---

for list_file in "$CHUNK_DIR"/*.list; do
    base_name=$(basename "$list_file" .list)
    output_file="$ARCHIVE_DIR/${base_name}.tar.zst"

    whiptail --infobox "Compressing $base_name..." 8 60
    tar -cf - -T "$list_file" | zstd -$ZSTD_LEVEL -T$ZSTD_THREADS -o "$output_file"
done

# --- RSYNC TRANSFER ---

for file in "$ARCHIVE_DIR"/*.tar.zst; do
    whiptail --infobox "Uploading $(basename "$file") to $DEST_HOST..." 8 60
    rsync -a "$file" "$DEST_HOST"
done

whiptail --msgbox "✅ All chunks compressed and transferred successfully." 10 50
