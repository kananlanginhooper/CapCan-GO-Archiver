package main

import (
	"fmt"
	"path/filepath"
)


func main() {
    // Step 1: Interactive directory input
    sourceDir := PromptForSourceDir()
    if sourceDir == "" {
        ShowInfoBox("Cancelled: no source directory selected.")
        return
    }

    // Step 2: Custom archive base name + auto timestamp
    archiveBase := PromptForArchiveName()
    if archiveBase == "" {
        ShowInfoBox("Cancelled: no archive name provided.")
        return
    }

    // Step 3: Destination selection (local or remote)
    destinationMode := PromptForDestination()
    if destinationMode == "" {
        ShowInfoBox("Cancelled: no destination selected.")
        return
    }

    // Step 4: Prompt for destination path
    destinationPath := PromptForDestinationPath(destinationMode)
    if destinationPath == "" {
        ShowInfoBox("Cancelled: no destination path provided.")
        return
    }

    // Step 5: Confirm access or permissions
    if !VerifyWriteAccess(destinationMode, destinationPath) {
        ShowInfoBox("Error: Cannot write to selected destination.")
        return
    }

    // Step 6: Scan and chunk the directory
    ShowInfoBox("Scanning directory and preparing chunks...")
    chunks, err := BuildChunks(sourceDir, ChunkSizeGB)
    if err != nil {
        ShowInfoBox("Error during chunking: " + err.Error())
        return
    }

    // Step 7: Update output names for each chunk
    for i := range chunks {
        chunks[i].OutputName = fmt.Sprintf("%s.part%03d.tar.zst", archiveBase, i)
    }

    // Step 8: Compress each chunk
    for _, chunk := range chunks {
        RunCompression(chunk, "./archives") // Save locally first
    }

    // Step 9: Transfer chunks to destination
    if destinationMode == "remote" {
        for _, chunk := range chunks {
            localFile := filepath.Join("./archives", chunk.OutputName)
            TransferChunkViaSSH(localFile, destinationPath)
        }
    } else {
        for _, chunk := range chunks {
            CopyChunkLocally(chunk.OutputName, destinationPath)
        }
    }

    ShowInfoBox("Backup completed successfully. 🚀")
}
