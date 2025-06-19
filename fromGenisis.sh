#!/bin/bash

# Goal: Compress, archive, and then rsync to a backup server, using a multi-threaded processor

# Goal: Iterate over each subdirectory in a base path, compress and archive them fully (no chunking),
# then rsync each archive to a backup server using a persistent SSH connection with multi-threaded compression.
#
# Example:
#   Local path: /shared/<Folder>/
#   Remote path: 10.100.0.60:/WD14TB1/
#   Archive format: archive.YYYY.MMDD.Folder-from-r730.tag.tar.zst
#
# Notes:
#   - Only asks for SSH password once using ControlMaster socket
#   - Parameters match those of rsync and support its options


# --- Configuration -----------------------------------------------------

SRC="/shared"                             # Base directory with many folders
DEST="10.100.0.60:/WD14TB1"               # Remote backup location
HOSTNAME="r730"                           # Label for archive filenames
STAMP=$(date '+%Y.%m%d')                  # Timestamp
THREADS=48                                # zstd thread count
TMPDIR=$(mktemp -d /tmp/archive_chunks.XXXXXX)  # Isolated temp space
SSH_USER="kanan"                          # Remote SSH user



# --- Step 0: Prepare persistent SSH connection for rsync----------------------

SSH_CTRL="/tmp/rsync-ctrl-$$"
SSH_OPTS="-o ControlMaster=auto -o ControlPath=$SSH_CTRL -o ControlPersist=10m"

# Kick off connection early (asks for password just once)
ssh -o ControlMaster=yes -o ControlPath="$SSH_CTRL" -o ControlPersist=10m \
    "${SSH_USER}@${DEST%%:*}" -N &



# --- Step 1: Iterate through each subdirectory ----------------------

for DIR in "$SRC"/*/; do
  base=$(basename "$DIR")

  # Temporary name while compressing
  tmpname="archive.${STAMP}.${base}-from-${HOSTNAME}.tmp.tar.zst"
  archive="/tmp/${tmpname}"

  echo "▶ Compressing: $base → $tmpname"
  tar -cf - "$DIR" | zstd -T${THREADS} -o "$archive"

  # Generate short hash (first 8 chars of SHA256)
  hash=$(sha256sum "$archive" | awk '{print $1}' | cut -c1-8)

  # Final filename with hash
  finalname="archive.${STAMP}.${base}-from-${HOSTNAME}.${hash}.tar.zst"
  finalpath="/tmp/${finalname}"
  mv "$archive" "$finalpath"

  echo "▶ Syncing: $finalname"
  rsync -e "ssh $SSH_OPTS" --progress "$finalpath" "${DEST}/" && rm -f "$finalpath"
done


# --- Step 2: Clean up SSH connection ------
ssh -O exit -o ControlPath="$SSH_CTRL" "${SSH_USER}@${DEST%%:*}"
echo "SSH connection closed."

# --- Step 3: clean up ---------------------
rm -rf "$TMPDIR"
echo "Backup completed successfully."