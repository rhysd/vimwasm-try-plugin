package main

import (
	"github.com/rhysd/go-fakeio"
	"strings"
	"testing"
)

func TestBuildURL(t *testing.T) {
	for _, tc := range []struct {
		what     string
		repo     string
		baseURL  string
		debug    bool
		rev      string
		expected string
	}{
		{
			what:     "simplest",
			repo:     "rhysd/clever-f.vim",
			baseURL:  "https://rhysd.github.io/vim.wasm/",
			debug:    false,
			expected: `https://rhysd.github.io/vim.wasm/?dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fcompat.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fcompat.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fplugin%2Fclever-f.vim`,
		},
		{
			what:     "base",
			repo:     "rhysd/clever-f.vim",
			baseURL:  "http://localhost:1234",
			debug:    false,
			expected: `http://localhost:1234?dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fcompat.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fcompat.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fplugin%2Fclever-f.vim`,
		},
		{
			what:     "debug",
			repo:     "rhysd/clever-f.vim",
			baseURL:  "https://rhysd.github.io/vim.wasm/",
			debug:    true,
			expected: `https://rhysd.github.io/vim.wasm/?debug=&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fcompat.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fcompat.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fplugin%2Fclever-f.vim`,
		},
		{
			what:     "revision",
			repo:     "rhysd/clever-f.vim",
			baseURL:  "https://rhysd.github.io/vim.wasm/",
			debug:    false,
			rev:      "ver1.1",
			expected: `https://rhysd.github.io/vim.wasm/?file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fver1.1%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fver1.1%2Fplugin%2Fclever-f.vim`,
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			f := fakeio.Stdout()
			defer f.Restore()

			o := &cliOptions{
				repo:     tc.repo,
				baseURL:  tc.baseURL,
				debug:    tc.debug,
				rev:      tc.rev,
				printURL: true,
			}

			if err := run(o); err != nil {
				t.Fatal(err)
			}

			actual, err := f.String()
			if err != nil {
				panic(err)
			}

			if actual != tc.expected {
				t.Fatalf("ACTUAL:\n%q\n\nEXPECTED:\n%q\n", actual, tc.expected)
			}
		})
	}
}

func TestInvalidURL(t *testing.T) {
	for _, tc := range []struct {
		what     string
		repo     string
		baseURL  string
		rev      string
		expected string
	}{
		{
			what:     "broken repo",
			repo:     "foo",
			expected: "Repository \"foo\" is invalid.",
		},
		{
			what:     "empty repo",
			expected: "Repository \"\" is invalid.",
		},
		{
			what:     "broken base url",
			repo:     "rhysd/clever-f.vim",
			baseURL:  ":localhost:1234",
			expected: "URL \":localhost:1234\" specified with -base is broken",
		},
		{
			what:     "invalid scheme in base url",
			repo:     "rhysd/clever-f.vim",
			baseURL:  "file://localhost:1234",
			expected: "Given URL with -base option does not have 'http' or 'https' scheme",
		},
		{
			what:     "not existing repo",
			repo:     "rhysd/not-existing-repository",
			expected: "404 Not Found",
		},
		{
			what:     "not a vim plugin",
			repo:     "rhysd/vimwasm-try-plugin",
			expected: "Repository \"rhysd/vimwasm-try-plugin\" contain no Vim script file (filename ends with .vim)",
		},
		{
			what:     "invalid revision",
			repo:     "rhysd/vimwasm-try-plugin",
			rev:      "this-ref-does-not-exist",
			expected: "404 No commit found for the ref this-ref-does-not-exist",
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			u := "https://rhysd.github.io/vim.wasm/"
			if tc.baseURL != "" {
				u = tc.baseURL
			}

			o := &cliOptions{
				repo:     tc.repo,
				baseURL:  u,
				rev:      tc.rev,
				printURL: true,
			}

			err := run(o)
			if err == nil {
				t.Fatal("Error did not occur")
			}

			msg := err.Error()
			if !strings.Contains(msg, tc.expected) {
				t.Fatalf("%q is not contained in error message %q", tc.expected, msg)
			}
		})
	}
}
