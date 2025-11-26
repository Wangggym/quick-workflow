package github

import (
	"testing"
)

func TestParsePRFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantPR    int
		wantErr   bool
	}{
		{
			name:      "HTTPS URL",
			url:       "https://github.com/brain/planning-api/pull/2001",
			wantOwner: "brain",
			wantRepo:  "planning-api",
			wantPR:    2001,
			wantErr:   false,
		},
		{
			name:      "HTTP URL",
			url:       "http://github.com/owner/repo/pull/123",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    123,
			wantErr:   false,
		},
		{
			name:      "URL without protocol",
			url:       "github.com/owner/repo/pull/456",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    456,
			wantErr:   false,
		},
		{
			name:      "URL with query parameters",
			url:       "https://github.com/owner/repo/pull/789?comments=all",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    789,
			wantErr:   false,
		},
		{
			name:      "URL with fragment",
			url:       "https://github.com/owner/repo/pull/999#discussion_r123456",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    999,
			wantErr:   false,
		},
		{
			name:      "URL with /files suffix",
			url:       "https://github.com/brain/planning-api/pull/2001/files",
			wantOwner: "brain",
			wantRepo:  "planning-api",
			wantPR:    2001,
			wantErr:   false,
		},
		{
			name:      "URL with /commits suffix",
			url:       "https://github.com/owner/repo/pull/123/commits",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    123,
			wantErr:   false,
		},
		{
			name:      "URL with /checks suffix",
			url:       "https://github.com/owner/repo/pull/456/checks",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    456,
			wantErr:   false,
		},
		{
			name:      "URL with /files and query params",
			url:       "https://github.com/owner/repo/pull/789/files?file-filters%5B%5D=.js",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantPR:    789,
			wantErr:   false,
		},
		{
			name:      "Invalid URL - missing pull segment",
			url:       "https://github.com/owner/repo/123",
			wantOwner: "",
			wantRepo:  "",
			wantPR:    0,
			wantErr:   true,
		},
		{
			name:      "Invalid URL - not a number",
			url:       "https://github.com/owner/repo/pull/abc",
			wantOwner: "",
			wantRepo:  "",
			wantPR:    0,
			wantErr:   true,
		},
		{
			name:      "Invalid URL - too short",
			url:       "github.com/owner",
			wantOwner: "",
			wantRepo:  "",
			wantPR:    0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, pr, err := ParsePRFromURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePRFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("ParsePRFromURL() owner = %v, want %v", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("ParsePRFromURL() repo = %v, want %v", repo, tt.wantRepo)
				}
				if pr != tt.wantPR {
					t.Errorf("ParsePRFromURL() pr = %v, want %v", pr, tt.wantPR)
				}
			}
		})
	}
}

func TestIsPRURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid HTTPS PR URL",
			url:  "https://github.com/owner/repo/pull/123",
			want: true,
		},
		{
			name: "Valid HTTP PR URL",
			url:  "http://github.com/owner/repo/pull/123",
			want: true,
		},
		{
			name: "Valid PR URL without protocol",
			url:  "github.com/owner/repo/pull/123",
			want: true,
		},
		{
			name: "Not a PR URL - just repo",
			url:  "https://github.com/owner/repo",
			want: false,
		},
		{
			name: "Not a PR URL - issue",
			url:  "https://github.com/owner/repo/issues/123",
			want: false,
		},
		{
			name: "Not a GitHub URL",
			url:  "https://gitlab.com/owner/repo/pull/123",
			want: false,
		},
		{
			name: "Just a number",
			url:  "123",
			want: false,
		},
		{
			name: "Empty string",
			url:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPRURL(tt.url); got != tt.want {
				t.Errorf("IsPRURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRepositoryFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "HTTPS URL",
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "SSH URL",
			url:       "git@github.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Short format",
			url:       "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "Invalid format - too many parts",
			url:       "owner/repo/extra",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name:      "Invalid format - single part",
			url:       "owner",
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRepositoryFromURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRepositoryFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("ParseRepositoryFromURL() owner = %v, want %v", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("ParseRepositoryFromURL() repo = %v, want %v", repo, tt.wantRepo)
				}
			}
		})
	}
}

