---
title: Installation
---

## Homebrew

Published releases are distributed through `magnusopera/homebrew-tap`.

```bash
brew tap magnusopera/tap
brew install pgpackage
```

## GitHub Releases

Each tagged release publishes zip archives for:

- macOS x64
- macOS arm64
- Linux x64
- Linux arm64
- Windows x64
- Windows arm64

Download the archive for your platform from [GitHub Releases](https://github.com/MagnusOpera/pgpackage/releases), extract it, and place `pgpackage` on your `PATH`.

## Verify the install

```bash
pgpackage --version
```
