# Go Backend Clean Architecture with DevSecOps best practices (WIP)

[![Latest Release][release-badge]][release-url]
[![Build Status][github-badge]][github-url]

[![Quality Gate Status][sonarcloud-status-badge]][sonarcloud-url]
[![Security Rating][sonarcloud-security-badge]][sonarcloud-url]
[![Maintainability Rating][sonarcloud-maintainability-badge]][sonarcloud-url]

[![Bugs][sonarcloud-bugs-badge]][sonarcloud-url]
[![Code Smells][sonarcloud-codesmells-badge]][sonarcloud-url]
[![Coverage][sonarcloud-coverage-badge]][sonarcloud-url]
[![Duplicated Lines (%)][sonarcloud-duplicated-badge]][sonarcloud-url]

[![License: Apache 2.0][badge-license]][link-license]

This project is a fork of the Go Backend Clean Architecture developed by [**Amit Shekhar**](https://amitshekhar.me).

The goal is to demonstrate the best practices to maintain automatically a GO project with tools like Renovate (fix dependency vulnerabilities), pre-commit, semantic release (versioning, changelog generation, etc.), GitHub Runner hardening, and other useful DevSecOps tools.

The details of the best practices will be added later.

## Details on the project used to demonstrate the DevSecOps best practices

A Go (Golang) Backend Clean Architecture project with Gin, MongoDB, JWT Authentication Middleware, Test, and Docker.

More details can be found on the following GitHub repository: [go-backend-clean-architecture][project-url]

[project-url]: https://github.com/amitshekhariitbhu/go-backend-clean-architecture
[release-badge]: https://img.shields.io/github/v/release/timoa/secure-go-backend-clean-architecture?logoColor=orange
[release-url]: https://github.com/timoa/secure-go-backend-clean-architecture/releases
[github-badge]: https://github.com/timoa/secure-go-backend-clean-architecture/workflows/Build/badge.svg
[github-url]: https://github.com/timoa/secure-go-backend-clean-architecture/actions?query=workflow%3ABuild
[sonarcloud-url]: https://sonarcloud.io/dashboard?id=timoa_secure-go-backend-clean-architecture
[sonarcloud-status-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=alert_status
[sonarcloud-security-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=security_rating
[sonarcloud-maintainability-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=sqale_rating
[sonarcloud-bugs-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=bugs
[sonarcloud-codesmells-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=code_smells
[sonarcloud-coverage-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=coverage
[sonarcloud-duplicated-badge]: https://sonarcloud.io/api/project_badges/measure?project=timoa_secure-go-backend-clean-architecture&metric=duplicated_lines_density
[badge-license]: https://img.shields.io/badge/License-Apache2-blue.svg
[link-license]: https://raw.githubusercontent.com/timoa/secure-go-backend-clean-architecture/master/LICENSE