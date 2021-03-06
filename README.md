vimwasm-try-plugin
==================
[![Linux CI Status][travis-badge]][travis-ci]
[![Windows CI status][appveyor-badge]][appveyor]

[`vimwasm-try-plugin`][repo] is a command line tool to try a Vim plugin hosted on GitHub using
[vim.wasm][]. You can instantly try Vim plugin and colorscheme without installing it on browser.
About vim.wasm, please visit [the repository][proj].

![screenshot](https://github.com/rhysd/ss/blob/master/vimwasm-try-plugin/main.gif?raw=true)

## Installation

Download and unzip an executable from [release page](https://github.com/rhysd/vimwasm-try-plugin/releases)
for your platform and put it in some `$PATH` directory.

Or build from source using Go toolchain.

```
go get -u github.com/rhysd/vimwasm-try-plugin
```

## Usage

```
vimwasm-try-plugin [{flags}] 'owner/name' [-- {args}]
```

For `{flags}`, please read `-help` output for more details. `{args}` are passed to command line
arguments of Vim run in Web Worker.

For example,

```
vimwasm-try-plugin 'rhysd/clever-f.vim'
```

This command opens:

https://rhysd.github.io/vim.wasm/?dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fcompat.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fcompat.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Fcp932.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Feucjp.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fautoload%2Fclever_f%2Fmigemo%2Futf8.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fplugin%2Fclever-f.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Frhysd%2Fclever-f.vim%2Fmaster%2Fplugin%2Fclever-f.vim

You can try [clever-f.vim](https://github.com/rhysd/clever-f.vim) in your browser without installing it.

For example,

```
vimwasm-try-plugin 'morhetz/gruvbox' -- -c colorscheme\ gruvbox /home/web_user/tryit.js
```

This command opens:

https://rhysd.github.io/vim.wasm/?arg=-c&arg=colorscheme+gruvbox&arg=%2Fhome%2Fweb_user%2Ftryit.js&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fairline&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fairline%2Fthemes&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Flightline&dir=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Flightline%2Fcolorscheme&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fairline%2Fthemes%2Fgruvbox.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fmorhetz%2Fgruvbox%2Fmaster%2Fautoload%2Fairline%2Fthemes%2Fgruvbox.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Fgruvbox.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fmorhetz%2Fgruvbox%2Fmaster%2Fautoload%2Fgruvbox.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fautoload%2Flightline%2Fcolorscheme%2Fgruvbox.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fmorhetz%2Fgruvbox%2Fmaster%2Fautoload%2Flightline%2Fcolorscheme%2Fgruvbox.vim&file=%2Fusr%2Flocal%2Fshare%2Fvim%2Fcolors%2Fgruvbox.vim%3Dhttps%3A%2F%2Fraw.githubusercontent.com%2Fmorhetz%2Fgruvbox%2Fmaster%2Fcolors%2Fgruvbox.vim

[gruvbox](https://github.com/morhetz/gruvbox) is applied by default and `tryit.js` source is opened.
You can preview colorscheme without installing it.

All files are fetched on memory. So they are cleaned up automatically when a browser tab is closed.
[vim.wasm][proj] is a Vim fork to run it on browser by compiling it to WebAssembly.

## Limitation

[vim.wasm][proj] is a Vim compiled to WebAssembly. So Vim is running on your browser and has some limitation.

- Shell commands are not available. So if the Vim plugin uses `system()` or other stuffs which try
  to execute shell commands, it does not work.
- The Vim is built with 'normal' feature set configuration. Some functionalities enabled in 'big' or 'huge' feature set
  are not available. For example, sign, conceal or profile.
- [vim.wasm][] fetches all plugin files before starting Vim. Fetching many files or a large file may slows Vim start up.

## TODO

- Add `-local` string option to specify local directory instead of using GitHub API
- Consider symlinks

## License

This repository is distributed under [the MIT license](./LICENSE.txt).

[repo]: https://github.com/rhysd/vimwasm-try-plugin
[vim.wasm]: https://rhysd.github.io/vim.wasm
[proj]: https://github.com/rhysd/vim.wasm
[travis-badge]: https://travis-ci.org/rhysd/vimwasm-try-plugin.svg?branch=master
[travis-ci]: https://travis-ci.org/rhysd/vimwasm-try-plugin
[appveyor-badge]: https://ci.appveyor.com/api/projects/status/qc4ghqlv2ki7omra/branch/master?svg=true
[appveyor]: https://ci.appveyor.com/project/rhysd/vimwasm-try-plugin/branch/master
