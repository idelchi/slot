# slot

Save and render named shell commands with Go template substitution

---

[![GitHub release](https://img.shields.io/github/v/release/idelchi/slot)](https://github.com/idelchi/slot/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/slot.svg)](https://pkg.go.dev/github.com/idelchi/slot)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/slot)](https://goreportcard.com/report/github.com/idelchi/slot)
[![Build Status](https://github.com/idelchi/slot/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/slot/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`slot` is a command-line tool for managing and executing shell commands with ease.

- Save commands with Go template variables and tags for organization in a simple YAML file
- Render commands with variable substitution using `--with KEY=VAL`
- Shell integration places rendered commands into your prompt for execution

## Installation

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/slot/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# Save a command with Go template variables
$ slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod
```

```sh
# Render a command with variable substitution
$ slot render deploy --with file=k8s.yml
kubectl apply -f k8s.yml
```

```sh
# List all saved slots
$ slot ls
NAME    TAGS     CMD
deploy  k8s,prod kubectl apply -f {{.file}}
```

```sh
# List slots filtered by tag
$ slot ls --tag k8s
```

```sh
# Remove a slot
$ slot rm deploy
```

## Data Storage

Slots are stored in YAML format at `~/.config/slot/slots.yaml`.

## Shell Integration

Generate shell completion snippets for command placement:

```sh
# Generate bash integration
$ slot completions bash

# Generate zsh integration
$ slot completions zsh
```

The integration allows `slot render` commands to place rendered output into your shell prompt for editing before execution.

Use `--yes` to execute the rendered command directly without editing.

## Commands

<details>
<summary><strong>save</strong> — Save a command slot</summary>

- **Usage:** `slot save <name> <command> [flags]`
- **Flags:**
  - `--tags` – Tags for the slot (repeatable)
  - `--force` – Overwrite existing slot

</details>

<details>
<summary><strong>render</strong> — Render a saved command slot</summary>

- **Usage:** `slot render <name> [flags]`
- **Flags:**
  - `--with KEY=VAL` – Variable substitution (repeatable)

</details>

<details>
<summary><strong>ls</strong> — List saved slots</summary>

- **Usage:** `slot ls [flags]`
- **Flags:**
  - `--tag` – Filter by tag (repeatable)

</details>

<details>
<summary><strong>rm</strong> — Delete a saved slot</summary>

- **Usage:** `slot rm <name>`

</details>

<details>
<summary><strong>path</strong> — Show the path to the slots file</summary>

- **Usage:** `slot path`

</details>

<details>
<summary><strong>completions</strong> — Generate shell completion snippets</summary>

- **Usage:** `slot completions <shell>`
- **Supported shells:** bash, zsh

</details>

## Shell insertion

Run `slot completions <shell>` to generate shell integration snippets.

You can then use `slot run <name> [--yes/-y]` to either place the rendered command
into your prompt for editing or execute it directly.

## Multiline commands

For multiline commands, use either a `$`:

```sh
slot save ls $'if [ "{{ .INPUT }}" = "true" ]; then\n  echo "{{ .OUTPUT }}";\nfi'
```

or edit the slot directly in the YAML file.
