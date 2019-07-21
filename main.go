package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

var vimDirs = []string{
	"autoload",
	"colors",
	// "compiler", This is unnecessary since shell commands are not available on vim.wasm
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

type cliOptions struct {
	repo     string
	baseUrl  string
	debug    bool
	printUrl bool
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
		return nil, nil, xerrors.Errorf("Coult not fetch /repos/%s/%s/contents for path %q: %v", owner, repo, path, err)
	}
	if res != nil && res.StatusCode == 404 {
		return nil, nil, xerrors.Errorf("File path %q of repository \"%s/%s\" not found", path, owner, repo)
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
				return nil, nil, xerrors.Errorf("Error while fetching %q: %w", path, err)
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

func run(o *cliOptions) error {
	slug := strings.SplitN(o.repo, "/", 2)
	if len(slug) <= 1 {
		return xerrors.Errorf("Repository %q is invalid. Please specify in user/repo format with -repo option", o.repo)
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
		return xerrors.Errorf("Could not fetch file entries in repo recursively: %w", err)
	}

	if len(files) == 0 {
		return xerrors.Errorf("Repository %q contain no Vim script file (filename ends with .vim)", o.repo)
	}

	sortContentsByPath(dirs)
	sortContentsByPath(files)

	u, err := url.Parse(o.baseUrl)
	if err != nil {
		return xerrors.Errorf("URL %q specified with -base is broken: %v", o.baseUrl, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return xerrors.Errorf("Given URL with -base option does not have 'http' or 'https' scheme: %s", u.Scheme)
	}

	params := url.Values{}
	if o.debug {
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

	if o.printUrl {
		fmt.Print(u.String())
		return nil
	}

	if err := browser.OpenURL(u.String()); err != nil {
		return xerrors.Errorf("Could not open URL with browser: %v", err)
	}

	return nil
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
	o := &cliOptions{}
	flag.StringVar(&o.repo, "repo", "", "Slug ('user/repo') of your Vim plugin (required)")
	flag.StringVar(&o.baseUrl, "base", "https://rhysd.github.io/vim.wasm/", "Base URL where vim.wasm is hosted")
	flag.BoolVar(&o.debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&o.printUrl, "url", false, "Print URL to stdout instead of opening it in browser")
	flag.Usage = usage
	flag.Parse()

	if err := run(o); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
}
