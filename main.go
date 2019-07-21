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
	repo       string
	baseURL    string
	debug      bool
	printURL   bool
	rev        string
	persistent bool
}

func isVimDirPath(p string) bool {
	for _, d := range vimDirs {
		if p == d || strings.HasPrefix(p, d+"/") {
			return true
		}
	}
	return false
}

func getContentsRecursive(ctx context.Context, api *github.RepositoriesService, owner, repo, ref, path string) ([]*github.RepositoryContent, []*github.RepositoryContent, error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	file, entries, res, err := api.GetContents(ctx, owner, repo, path, opts)
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

			fs, ds, err := getContentsRecursive(ctx, api, owner, repo, ref, e.GetPath())
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
		return xerrors.Errorf("Repository %q is invalid. Did you forgot giving an argument? Please specify in owner/repo format", o.repo)
	}

	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	client := http.DefaultClient
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(ctx, src)
	}
	api := github.NewClient(client)

	files, dirs, err := getContentsRecursive(ctx, api.Repositories, slug[0], slug[1], o.rev, "")
	if err != nil {
		return xerrors.Errorf("Could not fetch file entries in repo recursively: %w", err)
	}

	if len(files) == 0 {
		return xerrors.Errorf("Repository %q contain no Vim script file (filename ends with .vim)", o.repo)
	}

	sortContentsByPath(dirs)
	sortContentsByPath(files)

	u, err := url.Parse(o.baseURL)
	if err != nil {
		return xerrors.Errorf("URL %q specified with -base is broken: %v", o.baseURL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return xerrors.Errorf("Given URL with -base option does not have 'http' or 'https' scheme: %s", u.Scheme)
	}

	prefix := "/usr/local/share/vim"
	if o.persistent {
		prefix = "/home/web_user/.vim"
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
		if !o.persistent {
			if _, ok := existingDirs[dir.GetPath()]; ok {
				continue
			}
		}
		v := fmt.Sprintf("%s/%s", prefix, p)
		params.Add("dir", v)
	}

	for _, file := range files {
		v := fmt.Sprintf("%s/%s=%s", prefix, file.GetPath(), file.GetDownloadURL())
		params.Add("file", v)
	}
	u.RawQuery = params.Encode()

	if o.printURL {
		fmt.Print(u.String())
		return nil
	}

	if err := browser.OpenURL(u.String()); err != nil {
		return xerrors.Errorf("Could not open URL with browser: %v", err)
	}

	return nil
}

const usageHeader = `Usage: vimwasm-try-plugin [flags] 'owner/repo'

  vimwasm-try-plugin is a URL generator to try Vim plugin hosted on GitHub with
  https://rhysd.github.io/vim.wasm. The Vim was compiled to WebAssembly and runs
  in your browser. All plugin files will be fetched on memory and loaded by Vim.

  You can try Vim plugin without installing it on browser.

Example: Open vim.wasm URL including clever-f.vim plugin

  $ vimwasm-try-plugin 'rhysd/clever-f.vim'

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	var version bool

	o := &cliOptions{}
	flag.StringVar(&o.baseURL, "base", "https://rhysd.github.io/vim.wasm/", "Base URL where vim.wasm is hosted")
	flag.BoolVar(&o.debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&o.printURL, "url", false, "Print URL to stdout instead of opening it in browser")
	flag.StringVar(&o.rev, "revision", "", "Name of commit/branch/tag such as 'master', 'd2f17bb', 'v1.0.0'")
	flag.BoolVar(&o.persistent, "persistent", false, "Use ~/.vim instead of /usr/local/share/vim for persistently installing the plugin")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Println("1.0.0")
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(1)
	}

	o.repo = flag.Arg(0)

	if err := run(o); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
}
