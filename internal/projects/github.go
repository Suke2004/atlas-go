package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RepoStats struct {
	Owner          string    `json:"owner"`
	Repo           string    `json:"repo"`
	Stars          int64     `json:"stars"`
	Forks          int64     `json:"forks"`
	OpenIssues     int64     `json:"open_issues"`
	PrimaryLanguage string    `json:"primary_language"`
	LastPushedAt   time.Time `json:"last_pushed_at"`
	TechStack      []string  `json:"tech_stack"`
}

type GitHubIssue struct {
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type GitHubClient struct {
	httpClient *http.Client
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func ParseGitHubURL(rawURL string) (owner string, repo string, err error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", "", fmt.Errorf("empty GitHub URL")
	}

	// Handle shorthand "owner/repo"
	if !strings.Contains(rawURL, "://") && strings.Count(rawURL, "/") == 1 {
		parts := strings.Split(rawURL, "/")
		return parts[0], parts[1], nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("URL path must contain owner and repo")
	}

	owner = pathParts[0]
	repo = strings.TrimSuffix(pathParts[1], ".git")
	return owner, repo, nil
}

func (c *GitHubClient) FetchRepoStats(ctx context.Context, owner, repo string) (*RepoStats, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Atlas-Personal-OS")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github repo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var data struct {
		StargazersCount int64  `json:"stargazers_count"`
		ForksCount      int64  `json:"forks_count"`
		OpenIssuesCount int64  `json:"open_issues_count"`
		Language        string `json:"language"`
		PushedAt        string `json:"pushed_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode repo stats: %w", err)
	}

	lastPushed, _ := time.Parse(time.RFC3339, data.PushedAt)

	// Fetch top languages for tech stack
	languages, _ := c.FetchLanguages(ctx, owner, repo)

	return &RepoStats{
		Owner:           owner,
		Repo:            repo,
		Stars:           data.StargazersCount,
		Forks:           data.ForksCount,
		OpenIssues:      data.OpenIssuesCount,
		PrimaryLanguage: data.Language,
		LastPushedAt:    lastPushed,
		TechStack:       languages,
	}, nil
}

func (c *GitHubClient) FetchLanguages(ctx context.Context, owner, repo string) ([]string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Atlas-Personal-OS")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github languages api returned status %d", resp.StatusCode)
	}

	var langMap map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&langMap); err != nil {
		return nil, err
	}

	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}

	return languages, nil
}

func (c *GitHubClient) FetchOpenIssues(ctx context.Context, owner, repo string) ([]GitHubIssue, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=open&per_page=10", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Atlas-Personal-OS")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github issues api returned status %d", resp.StatusCode)
	}

	var issues []GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}
