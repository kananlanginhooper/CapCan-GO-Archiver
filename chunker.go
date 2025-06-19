package main

import (
    "errors"
    "os"
    "path/filepath"
    "strconv"
)


// FileChunk represents a group of files to be compressed together
type FileChunk struct {
    Files      []string // List of file paths in the chunk
    OutputName string   // Base output filename for this chunk
}

// BuildChunks scans a root directory and splits files into batches
// Each batch will hold ~targetSizeGiB of data or MaxFilesPerChunk entries
func BuildChunks(rootPath string, targetSizeGiB int64) ([]FileChunk, error) {
    var chunks []FileChunk
    var currentChunk []string
    var currentChunkSize int64
    fileIndex := 0

    info, err := os.Stat(rootPath)
    if err != nil || !info.IsDir() {
        return nil, errors.New("invalid source directory: " + rootPath)
    }

    maxBytes := targetSizeGiB * 1024 * 1024 * 1024

    err = filepath.Walk(rootPath, func(path string, info os.FileInfo, walkErr error) error {
        if walkErr != nil {
            return walkErr
        }
        if info.IsDir() {
            return nil
        }

        currentChunk = append(currentChunk, path)
        currentChunkSize += info.Size()

        if currentChunkSize >= maxBytes || len(currentChunk) >= MaxFilesPerChunk {
            chunks = append(chunks, FileChunk{
                Files:      currentChunk,
                OutputName: "archive_" + strconv.Itoa(fileIndex) + ".tar.zst",
            })
            fileIndex++
            currentChunk = nil
            currentChunkSize = 0
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    if len(currentChunk) > 0 {
        chunks = append(chunks, FileChunk{
            Files:      currentChunk,
            OutputName: "archive_" + strconv.Itoa(fileIndex) + ".tar.zst",
        })
    }

    return chunks, nil
}
