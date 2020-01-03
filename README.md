# Git Changelog

[![Build Status](https://travis-ci.org/c4s4/git-changelog.svg?branch=master)](https://travis-ci.org/c4s4/git-changelog)
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/git-changelog)](https://goreportcard.com/report/github.com/c4s4/git-changelog)
[![Codecov](https://codecov.io/gh/c4s4/git-changelog/branch/master/graph/badge.svg)](https://codecov.io/gh/c4s4/git-changelog)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This tool generates a Markdown changelog from Git logs.

![](spirale.png)

## Installation

### Unix users (Linux, BSDs and MacOSX)

Unix users may download and install latest *git-changelog* release with command:

```bash
sh -c "$(curl http://sweetohm.net/dist/git-changelog/install)"
```

If *curl* is not installed on you system, you might run:

```bash
sh -c "$(wget -O - http://sweetohm.net/dist/git-changelog/install)"
```

**Note:** Some directories are protected, even as *root*, on **MacOSX** (since *El Capitan* release), thus you can't install *git-changelog* in */usr/bin* for instance.

### Binary package

Otherwise, you can download latest binary archive at <https://github.com/c4s4/git-changelog/releases>. Unzip the archive, put the binary of your platform somewhere in your *PATH* and rename it *git-changelog*.

## Usage

To generate a Markdown changelog from Git logs, go in the repository directory and type:

```bash
$ git changelog
# Changelog

## 1.0.1 (2019-12-02)

- First change
- Second change

## 1.0.0 (2019-12-01)

- First change
- Second change
```

This prints the changelog on the terminal. You can also write changelog in a file with *-file* option:

```bash
$ git changelog -file CHANGELOG.md
```

You can also change regular expression for release tags with *-tag* option. Let's say that your release tags are dates in ISO format, you would generate changelog with following command:

```bash
$ git changelog -tag '^\d{4}-\d{2}-\d{2}$'
```

You can omit dates in release titles with *-nodate* option.

To get help, type `git changelog -help` on command line:

```bash
$ git changelog [-help] [-version] [-file changelog] [-tag regexp] [-nodate]
Print markdown changelog from git logs:
-help           To print this help
-version        To print version
-file changelog To write changelog in given file
-tag regexp     To set regexp for release tags (defaults to "^(v|V)?\d+.*$")
-nodate         To omit dates in releases titles
```

## Note about Git Commands

When you type `git foo`, and *foo* is not a known Git command, Git will look for *git-foo* executable in *PATH*. Thus when you type `git changelog`, executable *git-changelog* will run.

*Enjoy!*
