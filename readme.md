# modupwiz

Modules Update Wizard (`modupwiz`) is a helper to manage dependency updates.

The goal of this tool is to analyze the dependencies to detect the updates and displays information that allow to check what is inside each update.

## CLI

```
NAME:
   modupwiz - Modules Update Wizard

USAGE:
   modupwiz [global options] command [command options] 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --direct       Only direct modules (default: true)
   --explicit     Only explicit indirect modules (default: false)
   --indirect     All indirect modules (default: false)
   --compare      Display compare links from GitHub if possible (default: true)
   --versions     Display versions of modules (default: false)
   --pipe         Allow to pipe the command (default: false)
   --path value   File path to write the output. (Default: os.Stdout)
   --help, -h     show help
   --version, -v  print the version

```

## Examples

```console
$ modupwiz
|              MODULE              |                                COMPARE                                |
|----------------------------------|-----------------------------------------------------------------------|
| github.com/cpuguy83/go-md2man/v2 | https://github.com/cpuguy83/go-md2man/compare/v2.0.2...v2.0.3         |
| github.com/mattn/go-runewidth    | https://github.com/mattn/go-runewidth/compare/v0.0.9...v0.0.15        |
| github.com/xrash/smetrics        | https://github.com/xrash/smetrics/compare/039620a65673...1d8dd44e695e |
```

## Install

```bash
go install github.com/ldez/modupwiz@latest
```
