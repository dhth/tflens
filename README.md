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
`qa`, `staging`, and `prod`. If you want to compare modules across all three
environments, you can define a comparison in `tflens.yml`:

```yaml
compareModules:
  valueRegex: "v?(\\d+\\.\\d+\\.\\d)"
  comparisons:
    - name: all-envs
      attributeKey: source
      sources:
        - path: environments/qa/apps/main.tf
          label: qa
        - path: environments/staging/apps/main.tf
          label: staging
        - path: environments/prod/apps/main.tf
          label: prod
```

You can then compare the modules as follows.

```bash
tflens compare-modules all-envs
```

```text
module      qa       staging    prod     in-sync
module_a    1.1.1    1.1.1      1.1.1    ✓
module_b    1.0.8    1.0.1      1.0.0    ✗
module_c    1.0.5    -          -        -
```
