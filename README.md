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

**go**:

```sh
go install github.com/dhth/tflens@latest
```

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
tflens compare-modules all-envs
```

```text
module      dev       prod-us    prod-eu    in-sync
module_a    1.0.24    1.0.24     1.0.24     ✓
module_b    0.2.0     0.2.0      -          ✗
module_c    1.1.1     1.1.1      1.1.0      ✗
```
