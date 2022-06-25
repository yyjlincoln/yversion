package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	file, err := os.Open("static/index.html")
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
