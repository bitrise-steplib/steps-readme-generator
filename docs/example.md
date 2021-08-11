### Examples

A workflow to generate and update `README.md`:

```yaml
generate-readme:
    steps:
    - git::https://github.com/bitrise-steplib/steps-readme-generator.git@main:
        title: Generate README
```

Run the workflow locally by `bitrise run generate-readme`, then commit the generated file.

You can also include step-specific sections in the readme. The following example uses two files stored in the step repo:

```yaml
generate-readme:
    steps:
    - git::https://github.com/bitrise-steplib/steps-readme-generator.git@main:
        title: Generate README
        inputs:
        - example_section: docs/examples.md
        - contrib_section: docs/contributing.md
```

This will read the contents of `docs/examples.md` and `docs/contributing.md`, then include their contents in the final readme.