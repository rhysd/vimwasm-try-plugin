package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

var vimDirs = []string{
	"autoload",
	"colors",
	// "compiler", This is unnecessary since because shell commands are not available on vim.wasm
	"ftplugin",
	"indent",
	"plugin",
	"syntax",
	"ftdetect",
}

var existingDirs = map[string]struct{}{
	"autoload": struct{}{},
	"colors":   struct{}{},
	"ftplugin": struct{}{},
	"indent":   struct{}{},
	"plugin":   struct{}{},
	"syntax":   struct{}{},
}

func isVimDirPath(p string) bool {
	for _, d := range vimDirs {
		if p == d || strings.HasPrefix(p, d+"/") {
			return true
		}
	}
	return false
}

func getContentsRecursive(ctx context.Context, api *github.RepositoriesService, owner, repo, path string) ([]*github.RepositoryContent, []*github.RepositoryContent, error) {
	// TODO: Consider 'ref' option
	file, entries, res, err := api.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, nil, err
	}
	if res != nil && res.StatusCode == 404 {
		return nil, nil, fmt.Errorf("File path %q of repository \"%s/%s\" not found", path, owner, repo)
	}

	files := []*github.RepositoryContent{}
	dirs := []*github.RepositoryContent{}

	if file != nil && strings.HasSuffix(file.GetName(), ".vim") {
		files = append(files, file)
	}

	for _, e := range entries {
		t := e.GetType()
		if t == "file" {
			if strings.HasSuffix(e.GetName(), ".vim") {
				files = append(files, e)
			}
		} else if t == "dir" {
			if !isVimDirPath(e.GetPath()) {
				continue
			}

			dirs = append(dirs, e)

			fs, ds, err := getContentsRecursive(ctx, api, owner, repo, e.GetPath())
			if err != nil {
				return nil, nil, err
			}
			files = append(files, fs...)
			dirs = append(dirs, ds...)
		}
	}
	return files, dirs, nil
}

type byPath []*github.RepositoryContent

func (a byPath) Len() int {
	return len(a)
}
func (a byPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a byPath) Less(i, j int) bool {
	return strings.Compare(a[i].GetPath(), a[j].GetPath()) < 0
}

func sortContentsByPath(a []*github.RepositoryContent) {
	sort.Sort(byPath(a))
}

func dirContainsVimFile(path string, files []*github.RepositoryContent) bool {
	for _, f := range files {
		if strings.HasPrefix(f.GetPath(), path) {
			return true
		}
	}
	return false
}

func queryParam(key, value string) string {
	return key + "=" + url.QueryEscape(value)
}

func run(repo, baseURL string, debug, printURL bool) error {
	slug := strings.SplitN(repo, "/", 2)
	if len(slug) <= 1 {
		return fmt.Errorf("Repository %q is invalid. Please specify in user/repo format", repo)
	}

	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	client := http.DefaultClient
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(ctx, src)
	}
	api := github.NewClient(client)

	files, dirs, err := getContentsRecursive(ctx, api.Repositories, slug[0], slug[1], "")
	if err != nil {
		return err
	}

	sortContentsByPath(dirs)
	sortContentsByPath(files)

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	params := url.Values{}
	if debug {
		params.Set("debug", "")
	}

	for _, dir := range dirs {
		p := dir.GetPath()
		if !dirContainsVimFile(p, files) {
			continue
		}
		if _, ok := existingDirs[dir.GetPath()]; ok {
			continue
		}
		v := fmt.Sprintf("/usr/local/share/vim/%s", p)
		params.Add("dir", v)
	}

	for _, file := range files {
		v := fmt.Sprintf("/usr/local/share/vim/%s=%s", file.GetPath(), file.GetDownloadURL())
		params.Add("file", v)
	}
	u.RawQuery = params.Encode()

	if printURL {
		fmt.Print(u.String())
		return nil
	}

	return browser.OpenURL(u.String())
}

const usageHeader = `Usage: vimwasm-try-plugin {flags}

  vimwasm-try-plugin is a URL generator to try Vim plugin hosted on GitHub with
  https://rhysd.github.io/vim.wasm. The Vim was compiled to WebAssembly and runs
  in your browser. All plugin files will be fetched on memory and loaded by Vim.

  You can try Vim plugin without installing it on browser.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	repo := flag.String("repo", "", "Slug ('user/repo') of your Vim plugin (required)")
	baseURL := flag.String("base", "https://rhysd.github.io/vim.wasm/", "Base URL where vim.wasm is hosted")
	debug := flag.Bool("debug", false, "Enable debug logging")
	printURL := flag.Bool("url", false, "Print URL to stdout instead of opening it in browser")
	flag.Usage = usage
	flag.Parse()

	if err := run(*repo, *baseURL, *debug, *printURL); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
