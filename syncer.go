package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

// TransferChunkViaSSH copies a file to a remote server using rsync + sshpass
func TransferChunkViaSSH(localFile, remoteTarget string) {
    // Prompt for SSH password
    password, err := WhiptailPassword("Enter SSH password:")
    if err != nil || password == "" {
        ShowInfoBox("Transfer cancelled: no password entered.")
        return
    }

    userHost := remoteTarget
    remotePath := ""

    if split := filepath.SplitList(remoteTarget); len(split) == 1 {
        // Attempt to extract user@host:/path format
        colonIndex := len(userHost)
        for i, r := range userHost {
            if r == ':' {
                colonIndex = i
                break
            }
        }
        if colonIndex < len(userHost) {
            userHost = remoteTarget[:colonIndex]
            remotePath = remoteTarget[colonIndex+1:]
        }
    }

    cmd := exec.Command("sshpass", "-p", password,
        "rsync", "-avz", localFile, fmt.Sprintf("%s:%s", userHost, remotePath),
    )

    ShowInfoBox("Uploading: " + filepath.Base(localFile))
    if err := cmd.Run(); err != nil {
        ShowInfoBox("Failed: " + err.Error())
        return
    }
    ShowInfoBox("Uploaded: " + filepath.Base(localFile) + " ✅")
}

// CopyChunkLocally copies a chunk from the local archive dir to a given destination path
func CopyChunkLocally(filename, destPath string) {
    sourcePath := filepath.Join("./archives", filename)
    destFile := filepath.Join(destPath, filename)

    input, err := os.ReadFile(sourcePath)
    if err != nil {
        ShowInfoBox("Failed to read chunk: " + err.Error())
        return
    }

    err = os.WriteFile(destFile, input, 0644)
    if err != nil {
        ShowInfoBox("Failed to copy chunk: " + err.Error())
        return
    }

    ShowInfoBox("Copied: " + filename + " ✅")
}
