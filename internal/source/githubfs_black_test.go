package source_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/zbpack/internal/source"
)

func getGithubToken(t *testing.T) string {
	token, ok := os.LookupEnv("GITHUB_TOKEN")

	if !ok {
		t.Skip("no token (GITHUB_TOKEN) provided: skipping GitHub tests")
	}

	return token
}

func TestGitHubFsOpen_File(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	f, err := fs.Open("readme.md")
	if err != nil {
		if strings.Contains(err.Error(), "401 Bad credentials") {
			t.Skip("Skip due to 401 error.")
			return
		}

		t.Fatalf("error when opening: %v", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error when reading file: %v", err)
	}

	t.Log(content)
}

func TestGitHubFsOpen_Dir(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	f, err := fs.Open("")
	if err != nil {
		if strings.Contains(err.Error(), "401 Bad credentials") {
			t.Skip("Skip due to 401 error.")
			return
		}

		t.Fatal("error when opening directory:", err)
	}

	fileInfo, err := f.Readdir(-1)
	if err != nil {
		t.Fatal("error when reading directory:", err)
	}

	for _, fi := range fileInfo {
		t.Log(fi.Name())
	}
}

func TestGitHubFsOpenFile_WithWriteFlag(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	_, err := fs.OpenFile("readme.md", os.O_RDWR, 0)
	assert.ErrorIs(t, err, source.ErrReadonly)
}

func TestGitHubFsOpenFile_Ref(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zbpack", token, source.GitHubRef("9da82d05f3123cdb76b25d36c40cd12581e4eb82"))
	f, err := fs.OpenFile("go.mod", os.O_RDONLY, 0)
	if err != nil {
		if strings.Contains(err.Error(), "401 Bad credentials") {
			t.Skip("Skip due to 401 error.")
			return
		}

		t.Fatal("error when opening file:", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error when reading file: %v", err)
	}

	if !strings.Contains(string(content), "go 1.19") {
		t.Fatalf("unexpected content: %s", string(content))
	}

	t.Log(content)
}
