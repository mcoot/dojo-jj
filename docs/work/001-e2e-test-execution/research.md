---
id: 001-e2e-test-execution
name: Dojo CLI E2E Test Execution
description: Research how to run normal Go E2E tests against the compiled dojo CLI with isolated Jujutsu repositories and disposable workspaces.
status: complete
---

## Research question

How should dojo-jj structure end-to-end tests so they remain ordinary Go tests runnable from `go test`, an IDE, or a task target, while exercising the compiled `dojo` binary against isolated Jujutsu repositories and disposable workspace directories?

## Findings summary

The project is already positioned for a Go-native E2E harness: it has an `internal/e2e` placeholder package, uses `testify/suite` in unit tests, and builds the executable from `./cmd/dojo/main.go`. The selected package structure is to move E2E tests to a top-level `test/e2e` package so the artifact-level black-box boundary is obvious. The existing E2E package is only a placeholder today and is not yet registered with a `Test...` function, so it will not execute as a Go test suite until moved/registered.

The selected binary strategy is for E2E tests to require `DOJO_BIN`, pointing at the compiled `dojo` binary under test. The Go test harness should not directly invoke production Go code or build the binary implicitly. Tests should invoke `DOJO_BIN` with `os/exec`, set `cmd.Dir` to the active test workspace, and pass a controlled environment that disables user config leakage where possible.

For JJ repository isolation, the chosen default is a colocated template repository copied per test. This gives the harness reusable fixture state without sharing mutable state across tests, while matching the common JJ/Git workspace mode that local `jj 0.42.0` uses by default. Each test should copy the template into `t.TempDir`, perform all JJ and dojo actions inside that copy or derived workspace directories, and rely on Go's temp cleanup to prevent persistence between runs.

E2E execution should be opt-in. Because a top-level `test/e2e` package would still be discovered by `go test ./...`, the implementation should use an `e2e` build tag on E2E test files and provide a `task e2e` target that builds `out/dojo`, sets `DOJO_BIN`, and runs `go test -tags=e2e ./test/e2e`.

## Detailed findings

### Current project shape

- CLI entrypoint: `cmd/dojo/main.go` builds the app with `factory.BuildApp()`, calls `cmd.BuildCli(app)`, then executes Cobra.
- CLI construction: `internal/cmd/cli.go` creates root command `dojo` and currently mounts only `get`.
- Current behavior: `dojo get` calls `DojoService.GetWorkspace()`.
- Current service behavior: `GetWorkspace()` checks whether `jj` is available on `PATH`; otherwise it returns a `DojoError` with code `JJ_NOT_ON_PATH`.
- JJ dependency: `internal/dependencies/jj_client.go` uses `exec.LookPath("jj")` for availability. `AddWorkspace` exists in the interface but is not implemented yet.
- Existing unit tests use `testify/suite`.
- Existing E2E files:
  - `internal/e2e/workspace_pool_test.go` defines `WorkspacePoolE2ESuite`.
  - `internal/e2e/world/world.go` defines an empty `E2EWorld`.
  - There is no `Test...` function calling `suite.Run`, so the suite is not executable yet.
- Selected E2E location: move this harness to `test/e2e` to make the compiled-binary black-box boundary clearer.

### Constraints from the request

- E2E tests should be normal Go tests, runnable from command line or IDE.
- Tests should execute the compiled `dojo` CLI binary.
- Tests should not directly invoke internal Go code.
- Tests should run against isolated JJ repositories.
- Tests should operate against workspace environments that are trivial to clean up.
- Tests should not persist state between runs.
- Harness complexity should stay low.

### JJ facts relevant to test isolation

The local `jj` version inspected during research is `jj 0.42.0`.

`jj git init` creates a Git-backed repository. In this installed version, colocation is the default unless disabled. The selected default is to keep that default colocated mode for the template repository, because it better represents common user repositories. `jj git init --no-colocate <destination>` remains available for future compatibility tests that specifically need non-colocated behavior.

`jj workspace add <destination>` can create additional JJ workspaces. It supports:

- `--name <NAME>` to choose the workspace name.
- `-r, --revision <REVSETS>` to choose parents for the new working-copy commit.
- `-m, --message <MESSAGE>` for the new change description.
- `--sparse-patterns` with `copy`, `full`, or `empty`.

JJ user config can leak into tests. The help text says user config locations can be overridden with the `JJ_CONFIG` environment variable, and `JJ_CONFIG=` ignores user config files. The selected approach is to provide a test-owned JJ config path/file with explicit values, so subprocess behavior is deterministic without relying on the developer's global config.

### Go harness options

#### Option A: Auto-build inside Go test setup

Build `dojo` in E2E setup with `go build -o <temp>/dojo ./cmd/dojo/main.go`. Each test receives a fresh world rooted in `t.TempDir()`, creates a JJ repo there, and calls the compiled binary using `exec.CommandContext`.

Pros:

- Stays a normal Go test suite.
- Works from IDEs if build tags and package targets are configured.
- Exercises the actual compiled command.
- Does not require external scripts to prepare the binary.
- Cleanup is natural through `t.TempDir`.
- Enables helpful per-test helpers for stdout, stderr, exit code, working dir, and environment.

Cons:

- First test run pays a build cost.
- Package-level binary build requires careful concurrency handling if tests are parallelized.
- Tests must avoid importing production internals for behavior, while still using helper code for harness setup.

This was not selected because the E2E suite should always run against an explicit artifact supplied by `DOJO_BIN`.

#### Option B: Taskfile builds the binary, E2E tests consume `DOJO_BIN`

The Taskfile builds `out/dojo`, then runs `go test` with `DOJO_BIN=out/dojo`. E2E tests fail fast if `DOJO_BIN` is missing.

Pros:

- The binary under test is explicit.
- Useful for testing release artifacts or cross-compiled binaries.
- Avoids building from inside tests.

Cons:

- Less IDE-friendly unless IDE run configs provide `DOJO_BIN`.
- Adds a two-step workflow for local development.
- Easier to accidentally test a stale binary.

This is the selected binary strategy. It keeps the artifact under test explicit and avoids accidentally building a different binary from inside the test process. The tradeoff is that command-line and IDE runs need `DOJO_BIN` configured, likely through a `task e2e` helper and documented IDE run configuration.

#### Option C: `TestMain` builds once, all tests share package-level binary path

Use `TestMain(m *testing.M)` in `test/e2e` to build the binary once before running tests.

Pros:

- Efficient for many E2E tests.
- Keeps build logic centralized.

Cons:

- `TestMain` has less direct access to `testing.T` helpers.
- Cleanup and error reporting are clunkier than a suite setup or helper with `t.TempDir`.
- If the package grows multiple suites, shared mutable package state can become awkward.

This can work, but a suite-level or package helper is easier to evolve at this stage.

#### Option D: Shell/script-driven E2E wrapper

Use a shell script or Taskfile command to create temp repos, run commands, and assert behavior.

Pros:

- Simple for a small smoke test.
- Mirrors manual CLI usage.

Cons:

- Not as IDE-friendly.
- Assertions and fixtures become ad hoc quickly.
- Harder to compose with Go unit tests and package filtering.

This does not meet the "normal Go E2E suite" goal as well as the selected Taskfile plus `DOJO_BIN` path.

### Recommended harness shape

Create a small E2E world under `test/e2e/world` with responsibilities like:

- Resolve the required `DOJO_BIN` path once.
- Create a per-test root with `t.TempDir`.
- Create isolated JJ repositories under the test root.
- Run `jj` commands needed to prepare fixtures.
- Run `dojo` commands against a chosen workspace directory.
- Capture stdout, stderr, exit code, duration, and errors.
- Provide assertions/helpers for filesystem state and JJ state.

The production CLI should still be treated as a black box. E2E helper code can call `jj` and the compiled `dojo` binary, but should not instantiate `factory.App`, `cmd.BuildCli`, or service types. Binary compilation should happen outside the test process, for example in `task e2e`.

### Recommended isolation model

Chosen default: create or maintain a colocated template JJ repository, copy it into a per-test temp directory, and run each test against that copy. The template should be treated as read-only by convention. Any test-specific commits, workspace additions, files, or config should happen only inside the copied test root.

Default per-test layout:

```text
<test-temp>/
  bin/ or external DOJO_BIN
  jj-config/
  repos/
    primary/
      .jj/
  workspaces/
    ...
```

The exact layout can be simpler if using one `t.TempDir` per test and a package-scoped binary path. The important boundary is that repository data, workspaces, config, and command working directories all live under test-owned temp directories.

For command environments:

- Set `PATH` so the compiled `dojo` and installed `jj` are discoverable.
- Set `JJ_CONFIG` to a test-owned config path/file to avoid user config leakage while keeping explicit settings available.
- Set `NO_COLOR=1` or pass `--color=never` where output assertions need stability.
- Consider setting `HOME` to a test-owned directory if future behavior reads home-relative files.
- Use context timeouts around subprocesses to prevent hung tests.

The test-owned JJ config should include at least stable user identity and output-related settings, and it should explicitly preserve the selected colocated default:

```toml
user.name = "Dojo E2E"
user.email = "dojo-e2e@example.invalid"
ui.color = "never"
git.colocate = true

[signing]
behavior = "drop"
backend = "none"
```

### Binary build strategy

The selected binary strategy is:

- Require `DOJO_BIN` to be set.
- Fail fast with a clear test setup error if `DOJO_BIN` is empty or does not point to an executable file.
- Provide a `task e2e` target that builds `out/dojo` first and then runs `go test -tags=e2e ./test/e2e` with `DOJO_BIN` set.

This keeps E2E execution black-box and artifact-oriented. It does mean IDE runs need an environment variable configured, or developers should run the Taskfile target from the command line.

### Test package location and registration

E2E test files should use an `e2e` build tag so they do not run as part of ordinary `go test ./...`:

```go
//go:build e2e
```

The existing E2E suite should add an entry point like:

```go
func TestWorkspacePoolE2ESuite(t *testing.T) {
    suite.Run(t, &WorkspacePoolE2ESuite{})
}
```

The selected location is top-level `test/e2e`, with test helpers under `test/e2e/world` if that remains useful. Since these tests should not invoke production internals, this package should not import `internal/cmd`, `internal/factory`, or service packages. It should interact with the app through `DOJO_BIN` and external commands only.

### Task targets and developer workflow

Add task targets similar to:

```yaml
tasks:
  test:
    cmds:
      - go test ./...

  e2e:
    deps:
      - build
    env:
      DOJO_BIN: "{{.TASKFILE_DIR}}/out/dojo"
    cmds:
      - go test -tags=e2e ./test/e2e
```

This keeps everyday unit tests fast and dependency-light while making the E2E path explicit. IDE users can still run the suite by adding the `e2e` build tag and `DOJO_BIN` environment variable to their run configuration.

## Key insights and clarifications

- A normal Go E2E suite and a compiled-binary black-box test are compatible: the Go test process can own setup and assertions while subprocesses perform real CLI actions.
- The lowest-complexity cleanup boundary is `t.TempDir` per test.
- The selected repository lifecycle is a template repo copied per test: this gives better fixture reuse than creating a fresh repo from scratch while preserving per-test cleanup.
- Colocated JJ repos are the selected default because they match common JJ/Git workspace setups and the local `jj 0.42.0` default.
- `DOJO_BIN` is required for E2E tests, so the binary under test is always explicit.
- Config isolation matters early. The selected default is a test-owned JJ config with explicit identity, color, signing, and colocation settings.
- The current E2E suite is not executable yet because it lacks a `Test...` entry point.
- E2E tests should move to top-level `test/e2e` to make the black-box boundary explicit.
- E2E tests should be opt-in via an `e2e` build tag and `task e2e`; ordinary `go test ./...` should not require `DOJO_BIN` or JJ fixture setup.
- Existing `Taskfile.yaml` has no test or e2e task yet.

## Key questions for implementation

1. Resolved: the default harness should use a template JJ repository copied per test.
2. Resolved: test repositories should use the current JJ default colocated mode.
3. Resolved: E2E tests should require `DOJO_BIN` to point at the compiled binary under test.
4. Resolved: the E2E harness should isolate JJ config with a test-owned config path/file containing explicit settings.
5. Resolved: E2E tests should live under top-level `test/e2e`.
6. Resolved: E2E tests should be opt-in through an `e2e` build tag and `task e2e`.

## References

- `cmd/dojo/main.go`
- `internal/cmd/cli.go`
- `internal/cmd/get.go`
- `internal/factory/factory.go`
- `internal/service/dojo_service.go`
- `internal/dependencies/jj_client.go`
- `internal/e2e/workspace_pool_test.go`
- `internal/e2e/world/world.go`
- `Taskfile.yaml`
- Local command output from `jj --version`: `jj 0.42.0`
- Local command help from `jj help git init`
- Local command help from `jj help workspace add`
- Local command help from `jj help -k config`
