CapCan-GO-Archiver
==================

**CapCan-GO-Archiver** is a high-performance, console-only archive and backup utility built in Go. Designed for large-scale data preservation (e.g. 14+ TB), it compresses data into chunked .tar.zst archives using parallel processing, and syncs to a remote server—all while providing live terminal feedback through whiptail --infobox.

Perfect for systems administrators, data hoarders, and terminal aesthetes.

🔧 Features
-----------

*   🧵 **Multithreaded Compression**: Archives with tar piped to zstd -Tn, leveraging Go's concurrency primitives.
    
*   📦 **Chunked Packaging**: Splits large datasets into manageable archive chunks by size or file count.
    
*   🔁 **Concurrent Rsyncing**: Transfers each archive in parallel to a remote backup destination.
    
*   🖥️ **Console-Only Feedback**: Uses whiptail to display live task progress in a minimal interface—no GUI required.
    
*   ✨ **Fully Modular**: Clean, typed Go modules with separation of concerns (compression, chunking, syncing, and feedback).
    

📁 Project Structure
--------------------

CapCan-GO-Archiver/├── main.go # Entry point, concurrency orchestration├── config.go # Global settings: chunk size, thread count, etc.├── infobox.go # whiptail --infobox wrapper├── archiver.go # Compresses file chunks using tar + zstd├── chunker.go # Walks directory and creates batchable chunks├── rsyncer.go # rsync delivery of compressed archives└── go.mod # Module metadata

🚀 Usage
--------

go build -o capcan./capcan /your/data/path user@remote:/backup/destination

Example:

./capcan /mnt/archive-data backup@nas:/tank/vault

Note: Requires tar, zstd, rsync, and whiptail installed on your system.

⚙️ Configuration
----------------

Adjust chunk size, maximum files per chunk, compression threads, and infobox delays via config.go:

const ChunkSizeGiB int64 = 8const CompressionThreadCount int = 4const MaxFilesPerChunk int = 10000

Want runtime config (e.g. JSON)? Ask me—I’ll wire it in.

📋 Roadmap Ideas
----------------

*   \[ \] JSON config support
    
*   \[ \] Resume interrupted rsyncs
    
*   \[ \] Optional CRC or hash-based integrity checks
    
*   \[ \] Deduplication index
    
*   \[ \] Animated terminal spinners via ANSI tricks 😎
    

🛡 License
----------

MIT

✨ Credits
---------

Crafted with 🖤 and a terminal by kananlanginhooper. Architecture and comments designed in collaboration with Copilot, optimized for elegance and endurance.

Let me know if you’d like that saved to a local file, converted into ASCII-art badge format, or auto-inserted into your repo’s working tree layout! I’ve got you. 🧃📄🛠️ Want the LICENSE file next? MIT is clean and copy-ready. Say when