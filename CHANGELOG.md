# Changelog

All notable changes to pgpackage are documented in this file.

## [Unreleased]

- Improved the website home page responsive layout on mobile screens.

## [0.1.0]


- Added Linux x64 release archives alongside the existing macOS arm64 and Linux arm64 binaries.
- Reworked the website home page to better explain pgpackage as desired-state schema management for PostgreSQL.

**Full Changelog**: https://github.com/MagnusOpera/pgpackage/compare/0.0.2...0.1.0

## [0.0.2]


- Limited release artifacts to Linux arm64 and macOS arm64 only, removing Windows and macOS x64 targets.

**Full Changelog**: https://github.com/MagnusOpera/pgpackage/compare/0.0.1...0.0.2

## [0.0.1]


- Added Docusaurus website scaffolding, docs versioning, and GitHub Pages release deployment.
- Added changelog-gated CI, draft GitHub release automation, and Homebrew tap update workflows.
- Added build-time CLI version reporting and cross-platform release packaging targets.
- Added macOS Developer ID signing and notarization to the release workflow for GitHub Release and Homebrew distribution.
- Fixed tag-release artifact uploads from the hidden `.out/` directory and removed the broken Windows arm64 release leg.

**Full Changelog**: https://github.com/MagnusOpera/pgpackage/compare/b015550b6a7bbd28f781fb3662935e7752cd1532...0.0.1
