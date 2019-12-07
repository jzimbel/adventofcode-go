package input

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jzimbel/adventofcode-go/color"
)

const (
	aocURL    string = "https://adventofcode.com"
	userAgent string = "advent_of_code_go_input_downloader_jzimbel"
)

var (
	projectTempDirPath = filepath.Join(os.TempDir(), "adventofcode-go")
	userSessionIDPath  = filepath.Join(projectTempDirPath, ".USER_SESSION_ID")
	inputsDirPath      = filepath.Join(projectTempDirPath, "inputs")
)

func init() {
	err := os.MkdirAll(projectTempDirPath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(inputsDirPath, 0777)
	if err != nil {
		panic(err)
	}
}

// GetUserSessionID reads user's adventofcode.com session id from file, or asks them for it and stores it in a file.
func getUserSessionID() (string, error) {
	f, err := os.Open(userSessionIDPath)
	if err == nil {
		defer f.Close()
		// read id from file
		return readIDFromFile(f)
	}
	if os.IsNotExist(err) {
		// ask user for id
		return readIDFromPrompt()
	}
	return "", err
}

func readIDFromPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your session id? It’s needed for downloading puzzle inputs.")
	fmt.Printf("Your id is the value of the cookie named %s on %s.\n", color.B("session"), color.B(aocURL))
	var id string
	for {
		fmt.Print(color.R("> "))
		input, err := reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		input = bytes.TrimSpace(input)
		matched, err := regexp.Match("^[a-f0-9]+$", input)
		if err != nil {
			return "", err
		}
		if matched {
			fmt.Println("Thanks.")
			id = string(input)
			break
		}
		fmt.Fprintln(os.Stderr, "That’s not a valid id. It should consist only of digits and lowercase letters a-f. Please try again.")
	}
	f, err := os.Create(userSessionIDPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func readIDFromFile(f *os.File) (string, error) {
	idBytes := make([]byte, 192)
	count, err := f.Read(idBytes)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(idBytes[:count])), nil
}

// Get loads the contents of an input file into a string and returns it.
func Get(year int, day int) (string, error) {
	inputFilePath := getInputFilePath(year, day)
	stat, err := os.Stat(inputFilePath)
	switch {
	case err == nil && !stat.IsDir():
		return readInputFile(inputFilePath)
	case err == nil && stat.IsDir():
		return "", fmt.Errorf("Input path is a directory, expected a file: %s", inputFilePath)
	case os.IsNotExist(err):
		return downloadInput(year, day, inputFilePath)
	default:
		return "", err
	}
}

func getInputFilePath(year int, day int) string {
	return filepath.Join(inputsDirPath, fmt.Sprintf("%d-%02d", year, day))
}

func getInputURL(year int, day int) string {
	return fmt.Sprintf("%s/%d/day/%d/input", aocURL, year, day)
}

// Downloads input, saves it to file, and returns the content in a string for convenience
func downloadInput(year int, day int, inputFilePath string) (string, error) {
	fmt.Fprintf(os.Stderr, "Input file %s does not exist.\n", color.R(inputFilePath))
	fmt.Fprintf(os.Stderr, "Attempting to download puzzle input from %s\n", color.B(aocURL))
	sessionID, err := getUserSessionID()
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, getInputURL(year, day), nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	req.Header.Add("User-Agent", userAgent)

	attemptCount := 0
	for {
		attemptCount++
		resp, err := client.Do(req)
		if err != nil {
			if err.(*url.Error).Timeout() {
				if attemptCount == 1 {
					fmt.Fprintln(os.Stderr, "Request timed out. Trying again.")
				}
				if attemptCount >= 2 {
					fmt.Fprintln(os.Stderr, "Failed to download input after multiple retries. Giving up.")
					return "", err
				}
				fmt.Fprintln(os.Stderr, "Trying again.")
				continue
			}
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			f, err := os.Create(inputFilePath)
			if err != nil {
				return "", err
			}
			_, err = f.Write(body)
			if err != nil {
				return "", err
			}
			f.Close()
			fmt.Fprintf(os.Stderr, "%s Input downloaded and saved to %s.\n", color.G("Success."), color.B(inputFilePath))
			return readInputFile(inputFilePath)
		}
		return "", fmt.Errorf("server responded with a non-200 status code: %s", resp.Status)
	}
}

func readInputFile(inputFilePath string) (string, error) {
	b, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}
