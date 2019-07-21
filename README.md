vimwasm-try-plugin
==================

[`vimwasm-try-plugin`][repo] is a command line tool to open try a Vim plugin hosted on GitHub using
[vim.wasm][]. You can instantly try Vim plugin without installing it on browser. About vim.wasm,
please visit [the repository][proj].

![screenshot](https://github.com/rhysd/ss/blob/master/vimwasm-try-plugin/main.gif?raw=true)

## Installation

Please build from source at this moment.

```
go get -u github.com/rhysd/vimwasm-try-plugin
```

## Usage

```
vimwasm-try-plugin -repo 'owner/name'
```

For example,

```
vimwasm-try-plugin -repo 'rhysd/clever-f.vim'
```


This command opens:

https://rhysd.github.io/vim.wasm/?dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fcompat.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fcompat.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fplugin%2Fclever-f.vim

You can try [clever-f.vim](https://github.com/rhysd/clever-f.vim) in your browser without installing it.

All files are fetched on memory. So they are cleaned up automatically when a browser tab is closed.
[vim.wasm][proj] is a Vim fork for Web

## Limitation

[vim.wasm][proj] is a Vim compiled to WebAssembly. So entire Vim is running on your browser and has some limitation.

- Shell commands are not available. So if the Vim plugin uses `system()` or other stuffs which try
  to execute shell commands, it does not work.
- The Vim is built with 'normal' feature set configuration. Some functionalities enabled in 'big' or 'huge' feature set
  are not available. For example, sign, conceal or profile.
- [vim.wasm][] fetches all plugin files before starting Vim. Fetching many files or a large file may slows Vim start up.

## TODO

- Run CI
- Make release
- Add `-revision` string option to specify revision of the repository. `/repos/:owner/:repo/contents` has query parameter `ref` for it
- Add `-persistent` bool option to copy files to `~/.vim` instead of `/usr/local/share/vim`
- Add `-local` string option to specify local directory instead of using GitHub API

## License

This repository is distributed under [the MIT license](./LICENSE.txt).

[repo]: https://github.com/rhysd/vimwasm-try-plugin
[vim.wasm]: https://rhysd.github.io/vim.wasm
[proj]: https://github.com/rhysd/vim.wasm
