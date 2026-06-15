# AI Review Guidelines — `cli-terraform`

## Meta Information for AI Agents

**Review scope**: AI agents must always review the **entire pull request in full**,
not just the incremental diff. This means:
- Read and consider all files in the repository that are relevant to the changes,
  not only the lines added or modified.
- Evaluate the impact of changes in the context of the whole codebase.
- Ensure consistency, correctness, and adherence to project conventions across the
  full scope of the PR, including unchanged surrounding code and its dependencies.

---

This file is the source of truth used by AI review agents when reviewing pull requests
in the `cli-terraform` repository. Every rule below MUST be checked for the changed
code (additions and modifications). Do NOT flag pre-existing issues in code that is
not part of the change unless the PR directly touches it.

Each rule has a stable `id` so review tools can reference it.

---

## 1. General practices

- **GEN-01**: If newly added functionality is not yet Globally Available (GA), the
  code must contain a comment stating when it is expected to reach GA.
- **GEN-02**: When the change is part of a multi-repo change (e.g. requires a
  matching edgegrid-golang or terraform-provider-akamai change), branch names
  across all the repos must be identical.
- **GEN-03**: If the change depends on a newly added edgegrid-golang method or a
  newly added/changed Terraform resource/data source, this PR must be released
  together with the related PRs. State this clearly in the PR description.

## 2. Common coding practices (apply to every change)

- **COD-01**: Any change affecting customers must be reflected in `CHANGELOG.md`.
- **COD-02**: No secrets, real contract IDs, real Akamai employee IDs, or any
  other non-public data may appear anywhere in the code (especially unit tests
  and fixtures).
- **COD-03**: Never skip error checking. Every returned `error` must be handled.
- **COD-04**: Do not add files that are never used.
- **COD-05**: Prefer `any` over `interface{}`.
- **COD-06**: Acronyms in exported identifiers (e.g. `IP`, `DNS`, `URL`, `ID`,
  `API`) must be in all uppercase. Unexported identifiers are exempt.
- **COD-07**: Descriptions and comments that are full sentences must end with a
  period.
- **COD-08**: Unit tests must cover all corner cases — especially presence and
  absence of fields. Cases may be aggregated where it makes sense.

## 3. Changelog rules

- **CHG-01**: Entries describe what changed from the customer perspective of
  THIS project (cli-terraform).
- **CHG-02**: Fixes/changes based on GitHub issues must include a link to the
  issue, e.g.
  `([#436](https://github.com/akamai/terraform-provider-akamai/issues/436))`.
- **CHG-03**: Use past tense.
- **CHG-04**: Use backticks (`` ` ``) around proper names.
- **CHG-05**: Place entries at a random line position within the correct section
  and release to mitigate merge conflicts.
- **CHG-06**: Do NOT delete empty lines in the changelog.

## 4. PR hygiene

- **PRH-01**: Keep WIP / not-ready PRs in Draft.
- **PRH-02**: Branch must be rebased on top of its target branch before review.
- **PRH-03**: Commits squashed where it makes sense.
- **PRH-04**: Every commit message must include the appropriate JIRA number and
  the story title or a description of the change.

---

## 5. Cli-terraform specific rules

- **CLI-01**: If a string field can contain multiline characters (typical for
  description fields set from the UI, but not limited to them), the exported
  Terraform configuration must handle multi-line values correctly using the
  `multiline.tmpl` template (see e.g.
  `cli-terraform/pkg/providers/appsec/templates/multiline.tmpl`).
- **CLI-02**: In unit tests, always provide the FULL expected request when
  mocking API calls — never `mock.Anything`, `mock.MatchedBy`, or
  `mock.AnythingOfType`. The single allowed exception is the `context.Context`
  argument, for which `mock.Anything` is permitted.
- **CLI-03**: Integration tests must ensure that for every correct
  server-side configuration, the sequence `export` → `import` →
  `terraform plan` produces the correct state AND an empty plan.
- **CLI-04**: Use distinct CLI exit codes to communicate the nature of a
  failure:
    - exit code `1` → general / unexpected errors,
    - exit code `2` → operations that are explicitly not supported (e.g.
      exporting an `ALIAS` zone).
  This allows callers/scripts to distinguish user errors from unsupported
  scenarios.
- **CLI-05**: Export commands MUST follow the standard structure:
    1. The CLI entry point is a function `CmdCreateXxx(c *cli.Context) error`.
    2. The core logic lives in a separate unexported function
       `createXxx(ctx, params)`.
    3. All dependencies of the core function are bundled in a struct named
       `createXxxParams`.
  This keeps the logic independently testable without a CLI context.
- **CLI-06**: When exporting resources that have per-network state (e.g.
  Staging and Production activations), generate ONE separate entry per network
  — never a single combined entry.
- **CLI-07**: When exporting activation resources for an entity that has never
  been activated on any network, still emit the activation resource in the
  generated configuration but keep it commented out. This prevents accidental
  activation while leaving a ready-to-use template the user can uncomment.
- **CLI-08**: When handling errors returned from edgegrid-golang, prefer using
  `errors.Is(err, <package>.ErrXxx)` with exported sentinel errors over
  inspecting HTTP status codes or parsing error message strings. If a needed
  sentinel error does not exist in edgegrid-golang yet, request it from the
  library maintainers (or add it yourself in a companion PR).
- **CLI-09**: Each command and each parameter must have whole-flow unit tests
  that, for a given command and set of parameters, verify that the content of
  the generated files is correct.
- **CLI-10**: If a change to commands, flags, or arguments forces end-users to
  manually adjust how they invoke the export command, it is a breaking change.
  It must target the `feature/sp-breaking-changes` branch and is released only
  once a quarter or less often. When suspecting breaking changes in the PR, let
  someone from the DEVEXP team know to help identify the correct
  `sp-breaking-changes` branch.

## 6. General Go coding practices applicable here

- **GO-01**: Apply all of section 2 (Common coding practices) to every file.
- **GO-02**: Acronyms rule (`COD-06`) applies equally to template-driven code
  that produces exported identifiers in Terraform configurations.
- **GO-03**: When tests assert error scenarios, they must check the full error
  message (callstack may be omitted) — consistent with the broader DEVEXP
  testing approach.

