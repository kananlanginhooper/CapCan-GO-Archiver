package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

// ShowInfoBox displays a brief status message using whiptail
func ShowInfoBox(message string) {
    cmd := exec.Command("whiptail", "--infobox", message, "8", "50")
    _ = cmd.Run()
}

// WhiptailInput prompts the user for input with a text box
func WhiptailInput(prompt string) (string, error) {
    cmd := exec.Command("whiptail", "--inputbox", prompt, "10", "60", "3>&1", "1>&2", "2>&3")
    output, err := cmd.Output()
    return strings.TrimSpace(string(output)), err
}

// WhiptailPassword prompts the user for a password securely
func WhiptailPassword(prompt string) (string, error) {
    cmd := exec.Command("whiptail", "--passwordbox", prompt, "10", "60", "3>&1", "1>&2", "2>&3")
    output, err := cmd.Output()
    return strings.TrimSpace(string(output)), err
}

// PromptForSourceDir prompts for the directory to back up
func PromptForSourceDir() string {
    dir, _ := WhiptailInput("Enter the full path of the directory to back up:")
    return dir
}

// PromptForArchiveName prompts for base archive name, appends date
func PromptForArchiveName() string {
    base, _ := WhiptailInput("Enter archive base name (no extension):")
    if base == "" {
        return ""
    }
    date := time.Now().Format("2006.0102") // YYYY.MMDD
    return base + "_" + date
}

// PromptForDestination prompts to choose between local and remote
func PromptForDestination() string {
    cmd := exec.Command("whiptail", "--title", "Destination", "--menu", "Select where to send your backup:", "15", "50", "2",
        "local", "Save to local drive",
        "remote", "Send via SSH to a remote system",
    )
    output, _ := cmd.Output()
    return strings.TrimSpace(string(output))
}

// PromptForDestinationPath asks for the local dir or remote path (user@host:/path)
func PromptForDestinationPath(mode string) string {
    label := "Enter path to store backup:"
    if mode == "remote" {
        label = "Enter remote target (user@host:/path):"
    }
    dest, _ := WhiptailInput(label)
    return dest
}

// VerifyWriteAccess tests if we can write to the path, locally or over SSH
func VerifyWriteAccess(mode, target string) bool {
    if mode == "local" {
        testFile := target + "/.write_test"
        err := os.WriteFile(testFile, []byte("ok"), 0644)
        defer os.Remove(testFile)
        return err == nil
    }

    // For remote SSH: try creating and removing a temp file
    testCmd := fmt.Sprintf(`touch %s/.write_test && rm %s/.write_test`, target, target)
    userHost := strings.SplitN(target, ":", 2)[0]
    cmd := exec.Command("ssh", userHost, testCmd)
    err := cmd.Run()
    return err == nil
}
