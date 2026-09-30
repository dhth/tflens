<p align="center">
  <h1 align="center">tflens</h1>
  <p align="center">
    <a href="https://github.com/dhth/tflens/actions/workflows/main.yml"><img alt="Build Status" src="https://img.shields.io/github/actions/workflow/status/dhth/tflens/main.yml?style=flat-square"></a>
    <a href="https://github.com/dhth/tflens/actions/workflows/vulncheck.yml"><img alt="Vulnerability Check" src="https://img.shields.io/github/actions/workflow/status/dhth/tflens/vulncheck.yml?style=flat-square&label=vulncheck"></a>
  </p>
</p>

`tflens` lets you compare Terraform modules across environments.

Compare module attributes like `source` or `version` to see which environments are in sync and which have drifted. View the results in your terminal or generate an HTML report to share.

> [!NOTE]
> `tflens` is alpha software. Its behaviour and interface are likely to change for a while.

💾 Installation
---

### Pre-built binaries

Download a pre-built binary from the [latest release](https://github.com/dhth/tflens/releases/latest). See [Verifying release artifacts](#-verifying-release-artifacts) for instructions on verifying your download.

### Install from source

You can also install from source using the `go` toolchain:

```sh
go install github.com/dhth/tflens@latest
```

⚡️ Quick start
---

Create a configuration file if you do not already have one:

```sh
tflens config sample > tflens.yml
```

Edit the generated file with the paths to your `.tf` files and labels for your environments. The sample defines a comparison named `apps` and extracts version numbers from each module's `source` attribute.

Validate the configuration, then run the comparison:

```sh
tflens config validate
tflens compare-modules apps
```

For three environments named `dev`, `prod-us`, and `prod-eu`, the output might look like this:

```text
 module       dev        prod-us     prod-eu     in-sync

 module_a     1.0.24     1.0.24      1.0.24      ✓
 module_b     0.2.0      0.2.0       -           ✗
 module_c     1.1.1      1.1.1       1.1.0       ✗
```

`-` means the module is missing from that environment or does not define the selected attribute. By default, missing values count as out of sync; use `--ignore-missing-modules` to ignore their absence. Terminal comparisons exit with a non-zero status when modules are out of sync, so you can also use them in CI.

`>_` Commands
---

| Command                               | What it does                           |
|---------------------------------------|----------------------------------------|
| `tflens config sample`                | Print a sample configuration           |
| `tflens config validate`              | Validate the configuration             |
| `tflens compare-modules <COMPARISON>` | Compare modules for a named comparison |
| `tflens help`                         | Show all commands and flags            |

Run `tflens <command> --help` for details about a particular command.

⚙️ Configuration
---

`tflens` reads configuration from `tflens.yml` in the current directory by default. Use `--config-path` (or `-c`) with `config validate` or `compare-modules` to read a different file. Source paths are relative to the current directory, not the configuration file.

Each comparison has a name, an attribute to compare, and at least two sources. The generated sample configuration looks like this:

```yaml
# tflens.yml

compareModules:
  # Define one or more named comparisons.
  comparisons:
    - name: apps
      # The module attribute to compare, such as source or version.
      attributeKey: source
      # Compare at least two .tf files. Paths are relative to the current directory.
      sources:
        - path: environments/dev/virginia/apps/main.tf
          # This label appears in the comparison output.
          label: dev
        - path: environments/prod/virginia/apps/main.tf
          label: prod-us
        - path: environments/prod/frankfurt/apps/main.tf
          label: prod-eu

  # Optional. Extract version numbers instead of comparing the full attribute.
  # Applies to all comparisons unless overridden by a comparison.
  # Uses the first capture group; falls back to the full value if there is no match.
  valueRegex: 'v?(\d+\.\d+\.\d+)'
```

Add more entries to `comparisons` to compare other groups of modules. For Terraform Registry modules, use `attributeKey: version` to compare the configured `version` attributes rather than the `source` addresses.

### Extracting values

Without `valueRegex`, `tflens` compares the full attribute value. The regex above extracts `1.3.0` from a source such as `git@github.com:owner/repo//modules/module_a?ref=module-a-v1.3.0`.

Set `valueRegex` under `compareModules` to apply it to all comparisons, or under an individual comparison to override it for that comparison.

`tflens` uses the first capture group. If the regex does not match or has no capture group, it compares the original attribute value.

### Ignoring modules

To exclude specific modules, add `ignoreModules` to the comparison:

```yaml
ignoreModules:
  - module_x
  - module_y
```

📊 HTML reports
---

Generate a report you can open in a browser:

```sh
tflens compare-modules apps --output-format html
```

The report is written to `tflens-report.html` by default. Set `--html-output` to choose another path, `--html-title` to change the title, or `--html-template` to use a custom template.

![tflens HTML comparison report](https://tools.dhruvs.space/images/tflens/v0-1-0/html-report.png)

### Including diffs

If the compared values are version tags, you can include diffs in the report. Add `diffConfig` to the comparison, choosing the base and head environment labels and a command that generates the diff:

```yaml
diffConfig:
  baseLabel: prod-us
  headLabel: dev
  cmd: ["./scripts/generate-diff.sh", "apps"]
```

Provide your own script or command. `tflens` passes it these environment variables:

| Variable                  | What it contains                         |
|---------------------------|------------------------------------------|
| `TFLENS_DIFF_BASE_REF`    | Compared value from the base environment |
| `TFLENS_DIFF_HEAD_REF`    | Compared value from the head environment |
| `TFLENS_DIFF_MODULE_NAME` | Name of the module being compared        |

The command's stdout becomes the diff shown in the report. Enable diff generation with `--include-diffs`:

```sh
tflens compare-modules apps --output-format html --include-diffs
```

🔐 Verifying release artifacts
---

Each release includes checksums for all artifacts. The checksum file is signed using [cosign](https://docs.sigstore.dev/cosign/installation/) (version `3.1.3`).

Replace `x.y.z` below with the release version you want to verify.

1. Get the checksum and cosign signature bundle from the release:

    ```shell
    curl -sSLO https://github.com/dhth/tflens/releases/download/vx.y.z/tflens_x.y.z_checksums.txt
    curl -sSLO https://github.com/dhth/tflens/releases/download/vx.y.z/tflens_x.y.z_checksums.txt.sigstore.json
    ```

2. Verify the checksum file's signature:

    ```shell
    cosign verify-blob \
        --bundle tflens_x.y.z_checksums.txt.sigstore.json \
        --certificate-identity-regexp 'https://github\.com/dhth/tflens/\.github/workflows/.+' \
        --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
        tflens_x.y.z_checksums.txt
    ```

3. Download the archive for your platform and validate its checksum. For example, for Linux x86-64:

    ```shell
    curl -sSLO https://github.com/dhth/tflens/releases/download/vx.y.z/tflens_x.y.z_linux_amd64.tar.gz
    sha256sum --ignore-missing -c tflens_x.y.z_checksums.txt
    ```

4. Once both checks pass, extract the archive:

    ```shell
    tar -xzf tflens_x.y.z_linux_amd64.tar.gz
    ./tflens -h
    ```
