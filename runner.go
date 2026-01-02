package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	printMessage("Welcome to Go Runner!", 0, 1)
	printMessage(strings.Repeat("*", 40), 0, 1)

	for {
		rootOptions := []string{"Syntax", "Packages"}
		rootPrompt := buildOptionsPrompt(rootOptions, "Choose a category")
		rootChoice, err := getUserChoice(rootPrompt, len(rootOptions))
		if err != nil {
			printMessage("Failed to get valid input. Continuing to next selection.", 1, 1)
			continue
		}
		root := strings.ToLower(rootOptions[rootChoice-1])

		subFolders, err := getSubFolders(root)
		if err != nil {
			printMessage(fmt.Sprintf("Error listing subfolders: %v", err), 1, 1)
			continue
		}
		if len(subFolders) == 0 {
			printMessage("No subfolders found in the selected category.", 1, 1)
			continue
		}
		sort.Strings(subFolders) // Sort for consistent ordering
		subPrompt := buildOptionsPrompt(subFolders, "Choose a subfolder")
		subChoice, err := getUserChoice(subPrompt, len(subFolders))
		if err != nil {
			printMessage("Failed to get valid input. Continuing to next selection.", 1, 1)
			continue
		}
		sub := subFolders[subChoice-1]

		mainFilesRoot := filepath.Join(root, sub)
		mainFiles, err := getMainFiles(mainFilesRoot)
		if err != nil {
			printMessage(fmt.Sprintf("Error listing main files: %v", err), 1, 1)
			continue
		}
		if len(mainFiles) == 0 {
			printMessage("No main.go files found in the selected subfolder.", 1, 1)
			continue
		}
		sort.Strings(mainFiles) // Sort for consistent ordering
		fileOptions := make([]string, len(mainFiles))
		for i, file := range mainFiles {
			relPath, _ := strings.CutPrefix(file, mainFilesRoot+"/")
			fileOptions[i] = relPath
		}
		filePrompt := buildOptionsPrompt(fileOptions, "Choose a file to run")
		fileChoice, err := getUserChoice(filePrompt, len(mainFiles))
		if err != nil {
			printMessage("Failed to get valid input. Continuing to next selection.", 1, 1)
			continue
		}
		file := mainFiles[fileChoice-1]

		printMessage(strings.Repeat("=", 40), 1, 1)
		//printMessage(fmt.Sprintf("Running: %s", file), 0, 1)
		cmd := exec.Command("go", "run", file)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr // Capture stderr for better error visibility
		if err := cmd.Run(); err != nil {
			printMessage(fmt.Sprintf("Error running file: %v", err), 1, 1)
		}
		printMessage(strings.Repeat("=", 40), 0, 1)

		continueOptions := []string{"Yes", "No"}
		continuePrompt := buildOptionsPrompt(continueOptions, "Do you want to continue?")
		continueChoice, err := getUserChoice(continuePrompt, len(continueOptions))
		if err != nil {
			printMessage("Failed to get valid input. Exiting.", 1, 1)
			break
		}
		if continueChoice == 2 {
			printMessage(strings.Repeat("*", 40), 1, 1)
			printMessage("Goodbye!", 0, 1)
			break
		}
		printMessage(strings.Repeat("*", 40), 1, 1)
		printMessage("Let's run another Go app!", 0, 1)
		printMessage(strings.Repeat("*", 40), 0, 2)
	}
}

func printMessage(message string, beforeLines, afterLines int) {
	for i := 0; i < beforeLines; i++ {
		fmt.Println()
	}
	fmt.Print(message)
	for i := 0; i < afterLines; i++ {
		fmt.Println()
	}
}

func buildOptionsPrompt(options []string, title string) string {
	var builder strings.Builder
	builder.WriteString("\n----> " + title + " <----\n")
	for i, opt := range options {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, opt))
	}
	builder.WriteString("Enter your choice: ")
	return builder.String()
}

func getUserChoice(prompt string, maxSize int) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printMessage(prompt, 0, 0)
		if !scanner.Scan() {
			printMessage("Error reading input. Please try again.", 1, 0)
			continue
		}
		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil {
			printMessage("Invalid input: please enter a number. Try again.", 1, 0)
			continue
		}
		if choice < 1 || choice > maxSize {
			printMessage(fmt.Sprintf("Invalid choice: must be between 1 and %d. Try again.", maxSize), 1, 0)
			continue
		}
		return choice, nil
	}
}

func getSubFolders(rootFolder string) ([]string, error) {
	var subFolders []string
	err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != rootFolder && strings.Count(path, string(filepath.Separator)) == 1 {
			subFolders = append(subFolders, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return subFolders, nil
}

func getMainFiles(rootFolder string) ([]string, error) {
	var mainFiles []string
	err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Base(path) == "main.go" {
			mainFiles = append(mainFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mainFiles, nil
}
