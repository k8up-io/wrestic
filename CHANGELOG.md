# Go-skelleton Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Fixed
- Huge memory leak for stdin backups
### Changed
- Prune is no longer triggered after the backup. It has to be triggered individually. The baas operator has that implemented.

## [v0.0.3] - 2018-08-10
### Changed
- Ability to do backups via OpenShift stdout
- Warpperscript to correctly pass the arguments to wrestic in docker
- Adjustments to the metric handling
- Restic 0.9.2

## [v0.0.2] - 2018-07-27
### Added
- CI/CD pipeline

## [v0.0.1] - 2018-07-26
### Added
- Initial version
- Improved error detection and various bugfixes
- Timeout for initial snapshot listing, default: 30s

[Unreleased]: https://git.vshn.net/vshn/wrestic/
[v0.0.3]: https://git.vshn.net/vshn/wrestic/compare/v0.0.2...v0.0.3
[v0.0.2]: https://git.vshn.net/vshn/wrestic/compare/v0.0.1...v0.0.2
[v0.0.1]: https://git.vshn.net/vshn/wrestic/tree/v0.0.1