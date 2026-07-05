<!--
SPDX-FileCopyrightText: 2025 James Pond <james@cipher.host>

SPDX-License-Identifier: CC0-1.0
-->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-07-05

### Added

- Added service tests for task creation, updates, completion, deletion, collection behavior, statistics, validation, and bulk fixtures.
- Added repository tests for task, reward, collection, and statistics persistence behavior.
- Added API handler tests for task, collection, statistics, health, and error responses.
- Added API mapper tests for task, reward, collection, and statistics response DTOs.
- Added CORS middleware tests.
- Added database migration tests.
- Added test helpers and fixture data for backend tests.

### Changed

- Updated API responses to use DTOs instead of exposing internal task models directly.
- Updated route registration to support creating dedicated routers for tests and server setup.
- Updated collection persistence to distinguish normal and shiny versions of the same Pokémon.
- Improved SQLite migration and test database setup.
- Updated README to better describe Taskemon as a self-hosted gamified task manager backend.

### Fixed

- Fixed repository error handling paths that wrapped nil errors.
- Fixed task reward reveal behavior and missing reward handling.
- Fixed collection update behavior for duplicate Pokémon.
- Fixed statistics update behavior after task lifecycle operations.
- Fixed API error mapping for validation, not found, conflict, timeout, and internal errors.
- Fixed formatting and staticcheck issues.

## [0.1.0] - 2025-07-24

- Initial release.
