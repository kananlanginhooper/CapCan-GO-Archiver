package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

// RunCompression creates a tarball of a file chunk and compresses it using zstd
func RunCompression(chunk FileChunk, outputDir string) {
    // Ensure output directory exists before writing
    err := os.MkdirAll(outputDir, os.ModePerm)
    if err != nil {
        ShowInfoBox("Failed to create output directory: " + err.Error())
        return
    }

    // Fully-qualified output path for this archive
    outputPath := filepath.Join(outputDir, chunk.OutputName)

    // Prepare arguments to feed tar with a fixed list of input files
    args := append([]string{"-cf", "-"}, chunk.Files...)

    // Define `tar` command for packing data to stdout
    tarCmd := exec.Command("tar", args...)

    // Define `zstd` command to read from stdin and compress
    zstdCmd := exec.Command("zstd", "-T", fmt.Sprintf("%d", CompressionThreadCount), "-o", outputPath)

    // Setup pipe from tar stdout to zstd stdin
    tarStdout, err := tarCmd.StdoutPipe()
    if err != nil {
        ShowInfoBox("Pipe error: " + err.Error())
        return
    }
    zstdCmd.Stdin = tarStdout

    // Start both processes in the correct order
    err = tarCmd.Start()
    if err != nil {
        ShowInfoBox("Tar error: " + err.Error())
        return
    }

    err = zstdCmd.Start()
    if err != nil {
        ShowInfoBox("Zstd error: " + err.Error())
        return
    }

    // Show archive progress while the commands are running
    ShowInfoBox("Compressing " + chunk.OutputName)

    // Wait for both commands to finish
    tarCmd.Wait()
    zstdCmd.Wait()

    // Final status update after completion
    ShowInfoBox("Finished: " + chunk.OutputName + " ✅")

    // Optional: add sleep to improve UX pacing
    time.Sleep(InfoboxUpdateDelay)
}
