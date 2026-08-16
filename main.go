// fetchdown — unduh release asset GitHub tanpa browser.
// Usage: fetchdown owner/repo [asset-substring]
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const api = "https://api.github.com/repos/%s/releases/latest"

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fetchdown owner/repo [asset-substring]")
		os.Exit(2)
	}
	repo, filter := os.Args[1], ""
	if len(os.Args) > 2 {
		filter = os.Args[2]
	}

	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", fmt.Sprintf(api, repo), nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "fetchdown")
	resp, err := client.Do(req)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		fatal(fmt.Errorf("HTTP %d: %s", resp.StatusCode, b))
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		fatal(err)
	}

	var targets []asset
	for _, a := range rel.Assets {
		if filter == "" || strings.Contains(a.Name, filter) {
			targets = append(targets, a)
		}
	}
	if len(targets) == 0 {
		fatal(fmt.Errorf("no asset matching %q in %s %s", filter, repo, rel.TagName))
	}

	for _, a := range targets {
		if err := download(client, a); err != nil {
			fatal(err)
		}
		fmt.Printf("✓ %s (%d KB)\n", a.Name, a.Size/1024)
	}
}

func download(client *http.Client, a asset) error {
	out, err := os.Create(a.Name)
	if err != nil {
		return err
	}
	defer out.Close()
	req, _ := http.NewRequest("GET", a.BrowserDownloadURL, nil)
	req.Header.Set("User-Agent", "fetchdown")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d downloading %s", resp.StatusCode, a.Name)
	}
	_, err = io.Copy(out, resp.Body)
	return err
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "fetchdown:", err)
	os.Exit(1)
}
