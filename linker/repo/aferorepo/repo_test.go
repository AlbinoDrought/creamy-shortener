package aferorepo

import (
	"net/url"
	"testing"

	"github.com/spf13/afero"
)

func TestSetGet(t *testing.T) {
	fs := afero.NewMemMapFs()

	repo, err := Make(fs)
	if err != nil {
		t.Errorf("Make() error = %v, wantErr false", err)
	}

	if _, err := repo.Get("foo"); err == nil {
		t.Error("repo.Get(\"foo\") error = nil, wantErr true")
	}

	testURL := "https://example.com/foo?bar"
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		t.Errorf("unable to parse test url: %v", err)
	}

	err = repo.Set("foo", parsedURL)
	if err != nil {
		t.Errorf("repo.Set() error = %v, wantErr false", err)
	}

	retrievedURL, err := repo.Get("foo")
	if err != nil {
		t.Errorf("repo.Get(\"foo\") error = %v, wantErr false", err)
	}

	if retrievedURL.String() != testURL {
		t.Errorf("repo.Get() = %v, want %v", retrievedURL.String(), testURL)
	}

	// _should_ persist
	repo, err = Make(fs)
	if err != nil {
		t.Errorf("Make() error = %v, wantErr false", err)
	}

	retrievedURL, err = repo.Get("foo")
	if err != nil {
		t.Errorf("repo.Get(\"foo\") error = %v, wantErr false", err)
	}

	if retrievedURL.String() != testURL {
		t.Errorf("repo.Get() = %v, want %v", retrievedURL.String(), testURL)
	}
}
