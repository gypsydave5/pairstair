---
layout: page
title: Installation
permalink: /installation/
---

# Installation

PairStair can be installed in several ways. Choose the method that works best for your environment.

## Prerequisites

PairStair requires Go version 1.18 or later if building from source.

## Homebrew (Recommended)

The easiest way to install PairStair on macOS or Linux:

```bash
brew install gypsydave5/pairstair/pairstair
```

This will install the latest stable version and keep it updated with `brew upgrade`.

## Go Install

If you have Go installed, you can install directly:

```bash
go install github.com/gypsydave5/pairstair@latest
```

This will install the `pairstair` binary in your `$GOPATH/bin` or `$HOME/go/bin`. Make sure this directory is in your `PATH`.

## Download Binary

Download pre-built binaries from the [GitHub Releases](https://github.com/gypsydave5/pairstair/releases) page:

- **Linux**: `pairstair-linux-amd64` or `pairstair-linux-arm64`
- **macOS**: `pairstair-darwin-amd64` or `pairstair-darwin-arm64`
- **Windows**: `pairstair-windows-amd64`

After downloading:

1. Rename the binary to `pairstair` (or `pairstair.exe` on Windows)
2. Make it executable: `chmod +x pairstair`
3. Move it to a directory in your `PATH`

## Build from Source

Clone and build the project:

```bash
git clone https://github.com/gypsydave5/pairstair.git
cd pairstair
go build
```

This creates a `pairstair` binary in the current directory.

## Verify Installation

Test that PairStair is installed correctly:

```bash
pairstair --help
```

You should see the help output with all available options.

## Next Steps

- [User Guide](guide.html) - Learn how to use PairStair
- [Examples](examples.html) - See real-world usage scenarios
