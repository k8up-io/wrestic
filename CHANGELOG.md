# Wrestic Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.0.7] - 2018-11-08
### Changed
- Update snapshot webhook after every command that may change the repository
- Create a snapshot for each folder in /data

## [v0.0.6] - 2018-11-01
### Added
- Webhook after S3 restore
- Archive command
### Changed
- Refactoring code, every command has now its own go file

## [v0.0.5] - 2018-09-28
### Added
- Ability to post metrics to an arbitrary HTTP endpoint
### Fixed
- Too small buffer for output parsing

## [v0.0.4] - 2018-09-12
### Fixed
- Huge memory leak for stdin backups
- Metrics are updated more often
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

[Unreleased]: https://git.vshn.net/vshn/wrestic/compare/v0.0.6...master
[v0.0.6]: https://git.vshn.net/vshn/wrestic/compare/v0.0.5...v0.0.6
[v0.0.5]: https://git.vshn.net/vshn/wrestic/compare/v0.0.4...v0.0.5
[v0.0.4]: https://git.vshn.net/vshn/wrestic/compare/v0.0.3...v0.0.4
[v0.0.3]: https://git.vshn.net/vshn/wrestic/compare/v0.0.2...v0.0.3
[v0.0.2]: https://git.vshn.net/vshn/wrestic/compare/v0.0.1...v0.0.2
[v0.0.1]: https://git.vshn.net/vshn/wrestic/tree/v0.0.1
