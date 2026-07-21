package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	URL     string `json:"html_url"`
}

type Updater struct {
	RepoOwner string
	RepoName  string
	CurrentVersion string
	CheckTime time.Time
}

func NewUpdater(owner, repo, version string) *Updater {
	return &Updater{
		RepoOwner: owner,
		RepoName:  repo,
		CurrentVersion: version,
	}
}

func (u *Updater) CheckForUpdates() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.RepoOwner, u.RepoName)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}
	
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	
	u.CheckTime = time.Now()
	return &release, nil
}

func (u *Updater) IsNewVersionAvailable(release *GitHubRelease) bool {
	if release == nil {
		return false
	}
	return release.TagName != u.CurrentVersion
}