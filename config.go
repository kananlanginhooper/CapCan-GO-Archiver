package main

import (
	"time"
)

// CompressedOutputDir defines where compressed chunks will be stored
const CompressedOutputDir string = "./compressed"

// ChunkSizeGB sets the size per chunk in gibibytes
const ChunkSizeGB int64 = 1000 // in GB, 8 == 8G per archive chunk

// MaxFilesPerChunk provides a safety cap on file count per archive
const MaxFilesPerChunk int = 10000


// Compression level for zstd compression 7-9 suggested
const CompressionThreadCount int = 9

// CompressionThreadCount defines how many threads zstd should use
const CompressionThreadCount int = 48

// InfoboxUpdateDelay specifies delay between infobox refreshes
const InfoboxUpdateDelay time.Duration = 2 * time.Second

// create this project
// go mod init kananlanginhooper/CapCan-GO-Archiver
