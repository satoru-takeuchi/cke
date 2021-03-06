# Change Log

All notable changes to this project will be documented in this file.
This project employs a versioning scheme described in [RELEASE.md](RELEASE.md#versioning).

## [Unreleased]

## [1.17.10] - 2020-06-16

### Added
- Default gatherer to prometheus metrics (#353)

### Changed
- Update Kubernetes image to resolve k/k #92019 (#354)

## [1.17.9] - 2020-06-04

### Fixed
- Fix a sabakan template validation bug (#348)

## [1.17.8] - 2020-06-03

### Changed
- Do not mix new/old API servers (#340)

### Fixed
- Fix a sabakan integration bug and enhance cluster validation (#346)

## [1.17.7] - 2020-06-01

### Fixed
- Do not restart kubelet when cgroup_driver is changed (#342)
- Fix race bug (#343)
- Remove redundant test (#343, #344)
- Stabilize sonobuoy test (#341)

## [1.17.6] - 2020-05-29

### Added
- Add predicates and priorities support for kube-scheduler policy (#329)
- Add cgroup driver option to KubeletParams (#337)

### Changed
- Update Kubernetes to 1.17.6 (#335)
- Make machine scheduling of sabakan-integration better (#327, #331)
- Update etcd to 3.3.22 (#332)
- Fix a crash bug (#321)
- Use flannel instead of Calico in examples (#328, #338)

### Removed
- Remove `pod_subnet` field from cluster.yml (#334)

## [1.17.5] - 2020-05-12

Nothing changed.

## [1.17.4] - 2020-05-12

### Changed
- Update vault api version (#317)

### Fixed
- Fix node label bug (#316)

## [1.17.3] - 2020-04-21

### Changed
- Update to k8s 1.17.5 (#314)

## [1.17.2] - 2020-04-10

### Changed
- Fix resource application bug (#311, #312)

## [1.17.1] - 2020-04-07

### Changed
- Add run-on-vagrant target to sonobuoy/Makefile (#309)

## [1.17.0] - 2020-04-01

No change from v1.17.0-rc.1.

## [1.17.0-rc.1] - 2020-03-31

### Changed
- Add new op for upgrading Kubelet without draining nodes (#304)
- Update etcd: v3.3.19.1 (#303)
- Update images for Kubernetes 1.17 (#302)
- Add label for each role (#300)
- Server Side Apply (#299)
    - Kubernetes 1.17.4
    - CNI plugins 0.8.5
    - CoreDNS 1.6.7
    - Unbound 1.10.0

## Ancient changes

- See [release-1.16/CHANGELOG.md](https://github.com/cybozu-go/cke/blob/release-1.16/CHANGELOG.md) for changes in CKE 1.16.
- See [release-1.15/CHANGELOG.md](https://github.com/cybozu-go/cke/blob/release-1.15/CHANGELOG.md) for changes in CKE 1.15.
- See [release-1.14/CHANGELOG.md](https://github.com/cybozu-go/cke/blob/release-1.14/CHANGELOG.md) for changes in CKE 1.14.
- See [release-1.13/CHANGELOG.md](https://github.com/cybozu-go/cke/blob/release-1.13/CHANGELOG.md) for changes in CKE 1.13.
- See [release-1.12/CHANGELOG.md](https://github.com/cybozu-go/cke/blob/release-1.12/CHANGELOG.md) for changes in CKE 1.12.

[Unreleased]: https://github.com/cybozu-go/cke/compare/v1.17.10...HEAD
[1.17.10]: https://github.com/cybozu-go/cke/compare/v1.17.9...v1.17.10
[1.17.9]: https://github.com/cybozu-go/cke/compare/v1.17.8...v1.17.9
[1.17.8]: https://github.com/cybozu-go/cke/compare/v1.17.7...v1.17.8
[1.17.7]: https://github.com/cybozu-go/cke/compare/v1.17.6...v1.17.7
[1.17.6]: https://github.com/cybozu-go/cke/compare/v1.17.5...v1.17.6
[1.17.5]: https://github.com/cybozu-go/cke/compare/v1.17.4...v1.17.5
[1.17.4]: https://github.com/cybozu-go/cke/compare/v1.17.3...v1.17.4
[1.17.3]: https://github.com/cybozu-go/cke/compare/v1.17.2...v1.17.3
[1.17.2]: https://github.com/cybozu-go/cke/compare/v1.17.1...v1.17.2
[1.17.1]: https://github.com/cybozu-go/cke/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/cybozu-go/cke/compare/v1.17.0-rc.1...v1.17.0
[1.17.0-rc.1]: https://github.com/cybozu-go/cke/compare/v1.16.4...v1.17.0-rc.1
