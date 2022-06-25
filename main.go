package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"yversion/logging"
	"yversion/version"
)

const PORT = 40055

func main() {
	logging.Infof("Launching yVersion server...\n")

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/lts/stable", ltsVersion)

	logging.Successf("Serving on port %d\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	if err != nil {
		logging.Errorf("ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	execLocation := os.Args[0]
	currentDirectory, err := GetParentDirectory(execLocation)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(fmt.Sprintf("%v%s%v%s%v", currentDirectory, string(os.PathSeparator), "static", string(os.PathSeparator), "index.html"))
	if err != nil {
		logging.Errorf("Open: %v\n", err)
		return
	}

	http.ServeContent(w, r, "index.html", time.Now(), file)
}

// This is a very manual & temporary solution
func ltsVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	versionSpecString, err := json.Marshal(version.VersionSpec{
		Version:       "v1.4.0",
		Channel:       "stable",
		Description:   "Added the ability to check for updates. This version does not include the ability to fetch and install updates.",
		ReleasedAt:    time.Unix(1656126000, 0),
		CanAutoUpdate: false,
		URL:           "https://github.com/yyjlincoln-unsw/lts/releases/tag/v1.4.0",
	})
	if err != nil {
		http.Error(w, "Internal Server Error: Could not decode VersionSpec.", http.StatusInternalServerError)
		return
	}
	w.Write(versionSpecString)
}

func GetParentDirectory(path string) (string, error) {
	// Remove trailing "/"
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	sep := strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
	if len(sep) == 1 {
		return "", fmt.Errorf("no parent directory is available")
	}
	newPath := ""
	sepLast := len(sep) - 1
	for i, v := range sep {
		if i != sepLast {
			newPath += v
			if i != sepLast-1 {
				newPath += string(os.PathSeparator)
			}
		}
	}
	if newPath == "" {
		newPath = "/"
	}
	return newPath, nil
}
