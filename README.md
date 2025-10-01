# Step readme generator

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-readme-generator?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-readme-generator/releases)

Generate a `README.md` based on `step.yml` data

<details>
<summary>Description</summary>

This step can be used to generate a standard README.md based on step.yml data. The resulting readme will list the step's most important details, such as the description, a table of inputs and outputs, and links for further information.

The exact structure is defined in `README.md.gotemplate`

Additional step-specific content can be included from local Markdown files, such as example usage and information for contributors.
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

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

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `example_section` | Path to a Markdown file containing step-specific examples. If specified, the contents will be included in the Get started section. |  |  |
| `contrib_section` | Path to a Markdown file about step-specific information for contributors. If specified, the contents will be included in the Contributing section. |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-readme-generator/pulls) and [issues](https://github.com/bitrise-steplib/steps-readme-generator/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
