package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	printMessage("Welcome to Runner!", 0, 1)
	printMessage(strings.Repeat("*", 40), 0, 2)
	for {
		root := "syntax"
		rootUserChoice := getUserInput(`1. Syntax
2. Packages
Choose a number`, 2)
		if rootUserChoice == 2 {
			root = "packages"
		}
		folders, err := getSubFolders(root)
		if err != nil {
			log.Fatalf("Walk Directory Error: %v", err)
		}
		var builderSub strings.Builder
		builderSub.WriteString("\n")
		for i, v := range folders {
			builderSub.WriteString(fmt.Sprintf("%d: '%s'\n", i+1, v))
		}
		builderSub.WriteString("Choose a number")
		subUserChoice := getUserInput(builderSub.String(), len(folders))
		sub := folders[subUserChoice-1]
		mainFilesRoot := fmt.Sprintf("%s/%s", root, sub)
		mainFiles, err := getMainFiles(mainFilesRoot)
		if err != nil {
			log.Fatalf("Walk Directory Error: %v", err)
		}
		var builderFiles strings.Builder
		builderFiles.WriteString("\n")
		for i, v := range mainFiles {
			result, _ := strings.CutPrefix(v, mainFilesRoot+"/")
			builderFiles.WriteString(fmt.Sprintf("%d: '%s'\n", i+1, result))
		}
		builderFiles.WriteString("Choose a number")
		fileUserChoice := getUserInput(builderFiles.String(), len(mainFiles))
		file := mainFiles[fileUserChoice-1]
		printMessage(strings.Repeat("=", 40), 1, 1)
		cmd := exec.Command("go", "run", file)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			printMessage(fmt.Sprintf("Go Run Error: %v", err), 1, 1)
		}
		printMessage(strings.Repeat("=", 40), 0, 1)
		continueUserInput := getUserInput(`1. Yes
2. No
Do you want to continue? `, 2)
		if continueUserInput == 2 {
			printMessage(strings.Repeat("*", 40), 1, 1)
			printMessage("Bye!", 0, 1)
			break
		}
		printMessage(strings.Repeat("*", 40), 1, 1)
		printMessage("Run your Go apps again!", 0, 1)
		printMessage(strings.Repeat("*", 40), 0, 2)
	}
}

func printMessage(message string, beforeCount, afterCount int) {
	beforeNewLine := strings.Builder{}
	afterNewLine := strings.Builder{}
	if beforeCount > 0 {
		for range beforeCount {
			beforeNewLine.WriteString("\n")
		}
	}
	if afterCount > 0 {
		for range afterCount {
			afterNewLine.WriteString("\n")
		}
	}
	fmt.Print(beforeNewLine.String() + message + afterNewLine.String())
}

func getUserInput(prompt string, maxSize int) int {
	var numberUserChoice int
	for {
		printMessage(prompt+": ", 0, 0)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			printMessage(fmt.Sprintf("Reading Input Error: %v", err), 1, 1)
			continue
		}
		numberUserChoice, err = strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			printMessage(fmt.Sprintf("Invalid Input Error: %v", err), 1, 1)
			continue
		}
		if numberUserChoice < 1 || numberUserChoice > maxSize {
			printMessage("Invalid Number Error: the number is out of the range!", 1, 1)
			continue
		}
		break
	}
	return numberUserChoice
}

func getSubFolders(rootFolder string) ([]string, error) {
	var subFolders []string
	if err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Count(path, "/") == 1 {
			subFolders = append(subFolders, filepath.Base(path))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return subFolders, nil
}

func getMainFiles(rootFolder string) ([]string, error) {
	var mainFiles []string
	if err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.go" {
			mainFiles = append(mainFiles, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return mainFiles, nil
}
