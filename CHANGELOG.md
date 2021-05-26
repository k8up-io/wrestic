# Wrestic Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.0] 2021-05-26
### Changed
- Refactor Integration tests ([#90])
- Improved logging when wrestic message is not JSON ([#88])

### Added
- Pass multiple 'option' flags to restic ([#89])
- Docker image mirror and push to quay.io ([#92])

## [v0.2.3] 2021-03-30
### Changed
- Enable `go mod tidy` for renovate ([#80])
### Fixed
- Also send webhooks for restores ([#77])
- Reduce Restic progress FPS ([#82])
- Make Restic output parsing more resilient ([#81])

## [v0.2.2] 2021-03-09
### Fixed
- Wrong binary permissions in container ([#73])
### Added
- Nonroot variant of image ([#71])
- Simple version and build date injection ([#74])

## [v0.2.1] 2021-03-04
### Fixed
- Make waiting for locks more stable
- Don't drop k8sExe stderr output ([#27])
- Don't print the percentage when using stdin backups ([#28])
- Github actions
### Changed
- Update to Restic 0.12.0 ([#57])
- Adjust exit code handling to restic 0.12.0

## [v0.2.0] 2020-05-29
### Changed
- Complete rewrite of whole wrestic
- Better logging
- Better code design
- Better error handling
- Better maintainability
- Less deadlock potential with command IOs

## [v0.1.9] 2020-05-04
### Changed
- Use root in docker image to avoid permission issues
### Added
- Ability to pass tags to restic
- Race condition in error parsing

## [v0.1.8] 2020-01-03
### Fixed
- Using own restic fork with hotfix from https://github.com/restic/restic/issues/2319 in dockerbuild
- Race condition that makes wrestic hang on certain OpenShift clusters

## [v0.1.7] 2019-11-18
### Fixed
- Race condition with channels in the output parsing that caused wrestic to deadlock

## [v0.1.6] 2019-11-14
### Changed
- Bring version on par with K8up
- Webhook output now occurs after each PVC with metrics about that specific backup. The list with all the snapshots is sent after all PVCs finished. This should reduce the strain on webhook handling for very large backup sets.
- The PVC paths in the Restic snapshot get trimmed away, so we can seamlessly restore directly to a PVC without having to copy stuff around.
- Pass signals to the restic process
- Use Restic 0.9.5
- Realtime backup stats in container
### Fixed
- Make the short ID usable in for the restore


## [v0.0.10] - 2019-04-05
**Attention:** This release needs a custom version of Restic with the new dump
system. An upstream Pullrequest is currently open, with any luck it should
be merged for Restic 9.4+.

### Changed
- Unlock before and after each restic action
- Don't add timestamps to the stdin backups
- Use 'dump folder to tar' ability of restic for faster restores
- More robust Bucket creation
### Added
- First iteration of integration tests
- TravisCI
- Pod lookup from within Wrestic
### Fixed
- Much faster restore to S3

## [v0.0.9] - 2019-01-30
This change contains a complete redesign of wrestic. While keeping backwards
compatibility with older operator versions. Changes to the design contain:
- Better output handling (Webhook/prometheus/errors)
- No more snapshot listings in order to initialise the repository
- Created an API for the restic commandline

### Fixed
- Remove default unlock as this causes race conditions
- Archives not restoring all PVCs
### Added
- Ability to accept file extension for the stdin backup
### Changed
- Redesigned wrestic
- Removed timeout for snapshot list altogether

## [v0.0.8] - 2018-12-09
### Fixed
- Handle backup command quoting correctly
- Exit code 1 on SIGTERM
### Changed
- Don't run a shell in the Docker container
- Fail the whole backup job, if a single stdin/folder backup failed
- Fail the whole archive job, if a single restore failed

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

[Unreleased]: https://git.vshn.net/vshn/wrestic/compare/v0.3.0...master
[v0.3.0]: https://github.com/vshn/wrestic/releases/tag/v0.3.0
[v0.0.7]: https://git.vshn.net/vshn/wrestic/compare/v0.0.6...v0.0.7
[v0.0.6]: https://git.vshn.net/vshn/wrestic/compare/v0.0.5...v0.0.6
[v0.0.5]: https://git.vshn.net/vshn/wrestic/compare/v0.0.4...v0.0.5
[v0.0.4]: https://git.vshn.net/vshn/wrestic/compare/v0.0.3...v0.0.4
[v0.0.3]: https://git.vshn.net/vshn/wrestic/compare/v0.0.2...v0.0.3
[v0.0.2]: https://git.vshn.net/vshn/wrestic/compare/v0.0.1...v0.0.2
[v0.0.1]: https://git.vshn.net/vshn/wrestic/tree/v0.0.1

[#27]: https://github.com/vshn/wrestic/pull/27
[#28]: https://github.com/vshn/wrestic/pull/28
[#57]: https://github.com/vshn/wrestic/pull/57
[#58]: https://github.com/vshn/wrestic/pull/58
[#71]: https://github.com/vshn/wrestic/pull/71
[#73]: https://github.com/vshn/wrestic/pull/73
[#74]: https://github.com/vshn/wrestic/pull/74
[#77]: https://github.com/vshn/wrestic/pull/77
[#80]: https://github.com/vshn/wrestic/pull/80
[#81]: https://github.com/vshn/wrestic/pull/81
[#82]: https://github.com/vshn/wrestic/pull/82
[#88]: https://github.com/vshn/wrestic/pull/88
[#89]: https://github.com/vshn/wrestic/pull/89
[#90]: https://github.com/vshn/wrestic/pull/90
[#92]: https://github.com/vshn/wrestic/pull/92
