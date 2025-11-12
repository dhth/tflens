<p align="center">
  <h1 align="center">tflens</h1>
  <p align="center">
    <a href="https://github.com/dhth/tflens/actions/workflows/main.yml"><img alt="main" src="https://img.shields.io/github/actions/workflow/status/dhth/tflens/main.yml?style=flat-square"></a>
    <a href="https://github.com/dhth/tflens/actions/workflows/vulncheck.yml"><img alt="vulncheck" src="https://img.shields.io/github/actions/workflow/status/dhth/tflens/vulncheck.yml?style=flat-square&label=vulncheck"></a>
  </p>
</p>

`tflens` offers tiny utilities for terraform/opentofu/terragrunt codebases.

> [!NOTE]
> `tflens` is alpha software. It's behaviour and interface is likely to change
> for a while.

Install
---

**homebrew**:

```sh
brew install dhth/tap/tflens
```

**go**:

```sh
go install github.com/dhth/tflens@latest
```

Or get the binary directly from a [release][1]. Read more about verifying the
authenticity of released artifacts [here](#-verifying-release-artifacts).

Usage
---

Consider a terragrunt codebase with three different deployment environments:
`dev`, `prod-us`, and `prod-eu`. If you want to compare modules across all three
environments, you can define a comparison in `tflens.yml`:

```yaml
compareModules:
  # list of configured comparisons
  comparisons:
    # will be used when specifying the comparison to be run
    - name: apps
      # the attribute to use for comparison
      attributeKey: source
      # where to look for terraform files
      sources:
        - path: environments/dev/virginia/apps/main.tf
          # this label will appear in the comparison output
          label: dev
        - path: environments/prod/virginia/apps/main.tf
          label: prod-us
        - path: environments/prod/frankfurt/apps/main.tf
          # regex to extract the desired string from the attribute value
          # only applies to this source, overrides the global valueRegex
          # optional
          valueRegex: "v?(\\d+\\.\\d+\\.\\d+)"
          label: prod-eu

  # regex to extract the desired string from the attribute value
  # applies to all comparisons
  # optional
  valueRegex: "v?(\\d+\\.\\d+\\.\\d+)"
```

You can then compare the modules as follows.

```bash
tflens compare-modules -h
```

```
Usage:
  tflens compare-modules <COMPARISON> [flags]

Flags:
  -c, --config-path string       path to tflens' configuration file (default "tflens.yml")
  -h, --help                     help for compare-modules
      --html-output string       path where the HTML report should be written (default "tflens-report.html")
      --html-template string     path to a custom HTML template (optional)
      --html-title string        title for the HTML report (default "report")
  -i, --ignore-missing-modules   to not have the absence of a module lead to an out-of-sync status
  -o, --output-format string     output format for results; allowed values: [stdout html] (default "stdout")
```

```bash
tflens compare-modules apps
```

```text
module      dev       prod-us    prod-eu    in-sync
module_a    1.0.24    1.0.24     1.0.24     ✓
module_b    0.2.0     0.2.0      -          ✗
module_c    1.1.1     1.1.1      1.1.0      ✗
```

🔐 Verifying release artifacts
---

In case you get the `tflens` binary directly from a [release][1], you may want to
verify its authenticity. Checksums are applied to all released artifacts, and
the resulting checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace `A.B.C` in the commands listed below with the version
you want):

1. Download the following files from the release:

    - tflens_A.B.C_checksums.txt
    - tflens_A.B.C_checksums.txt.pem
    - tflens_A.B.C_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob tflens_A.B.C_checksums.txt \
       --certificate tflens_A.B.C_checksums.txt.pem \
       --signature tflens_A.B.C_checksums.txt.sig \
       --certificate-identity-regexp 'https://github\.com/dhth/tflens/\.github/workflows/.+' \
       --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/tflens/releases/download/vA.B.C/tflens_A.B.C_linux_amd64.tar.gz
   sha256sum --ignore-missing -c tflens_A.B.C_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf tflens_A.B.C_linux_amd64.tar.gz
   ./tflens -h
   # profit!
   ```

[1]: https://github.com/dhth/tflens/releases
