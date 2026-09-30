# tflens

`tflens` compares Terraform module attributes across environments and produces
terminal or HTML reports. It parses local `.tf` files; it does not run Terraform
or evaluate interpolations and other dynamic expressions.

This project uses `mise` for tool management and tasks. Always use `mise` to
execute project commands; see `mise.toml` for the available tasks. Use `mise run
all` to format, lint, test, and build.

```text
YAML configuration → HCL parsing → comparison → terminal/HTML report
```

- `internal/cmd` defines the Cobra CLI, loads configuration, and handles output
  files and command errors. `main.go` handles process exit and error reporting.
- `internal/domain` contains configuration parsing and validation, comparison
  models, and module statuses.
- `internal/hcl` parses module blocks and extracts attribute values.
- `internal/services` orchestrates comparisons and optional external diff
  commands.
- `internal/view` renders comparison results; the built-in HTML template lives
  in `internal/view/assets/template.html`.

Keep comparison logic in `services`, HCL details in `hcl`, and presentation in
`view`. Renderers consume domain results rather than reading Terraform files or
computing comparison statuses. Prefer pure transformations and explicit inputs
and outputs while keeping the code idiomatic to Go; isolate filesystem and
subprocess effects from transformation logic. Use local mutation when clearer.

Preserve deterministic output: modules are sorted by name, and source columns
follow configuration order. Source paths are relative to the working directory,
not the configuration file. Out-of-sync modules cause a non-zero exit for
terminal output, but do not prevent successful HTML report generation. External
diff commands run only when `--include-diffs` is enabled.

Prefer integration tests in `tests/cli` for user-visible behavior, following the
existing binary fixture and `testdata` patterns. Use focused package tests for
configuration, comparison, and rendering logic. Snapshot tests use `go-snaps`;
update them with `mise run update-snapshots` and inspect the resulting diffs.
