
# did [![Build](https://travis-ci.org/Link512/godid.svg?&branch=master)](https://travis-ci.org/Link512/godid) [![Go Report Card](https://goreportcard.com/badge/github.com/Link512/godid)](https://goreportcard.com/report/github.com/Link512/godid) [![codecov](https://codecov.io/gh/Link512/godid/branch/master/graph/badge.svg)](https://codecov.io/gh/Link512/godid)

![logo](https://i.imgur.com/FpcrltN.png)

Simple task tracker written in `go`. Use it to quickly write down tasks that you've completed and then access summaries for daily/weekly standup purposes.

## Install

```bash
go get -u gopkg.in/Link512/godid.v1/did
```

## Usage

```text
Usage:
  did [flags]
  did [command]

Available Commands:
  help        Help about any command
  last        Displays the tasks logged in the last custom day duration
  lastWeek    Displays the tasks logged last week
  thisWeek    Displays the tasks logged this week
  today       Displays the tasks logged today
  yesterday   Displays the tasks logged yesterday

Flags:
  -e, --entry string   Entry to log
  -h, --help           help for did
```

## Examples

### Logging a single entry

![Screen2](https://i.imgur.com/NxiKuv2.png)

### Logging multiple entries

Run `did` with no arguments and write each entry on a new line. Press `Ctrl-d` to exit.

![Screen2](https://i.imgur.com/A7ws0YH.png)

### Getting today's summary

![Screen3](https://i.imgur.com/u9UIqwX.png)

### Getting this week's summary, per day

![Screen4](https://i.imgur.com/386ikhB.png)

### Getting last week's summary, flat

![Screen5](https://i.imgur.com/E1qpXSS.png)

### Getting a custom interval summary

![Screen6](https://i.imgur.com/8tEt6it.png)

## Configuration

After first running the tool, a default config file will be present at `~/.godid/config.yml` (also works on Windows). The config file contains only `store_path` to indicate where the entries are stored. The default for this value is `store_path: ~/.godid/store.db`.

## Notes

This is meant to be a very simple tool to keep track of things you do and present a nice summary of them. Chances are I might add some other features to it, but very minor ones in order to keep it from being bloated.

## TODO

[ ] Export summaries
