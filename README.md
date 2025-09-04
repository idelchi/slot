# slot

Save and render named shell commands with Go template substitution

---

[![GitHub release](https://img.shields.io/github/v/release/idelchi/slot)](https://github.com/idelchi/slot/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/slot.svg)](https://pkg.go.dev/github.com/idelchi/slot)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/slot)](https://goreportcard.com/report/github.com/idelchi/slot)
[![Build Status](https://github.com/idelchi/slot/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/slot/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`slot` is a command-line tool for saving, organizing, and executing templated shell commands.

- Save commands with Go template variables and tags for organization in a simple YAML file
- Render commands with variable substitution using `KEY=VAL`
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
$ slot render deploy file=k8s.yml
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

Generate shell integration snippets for command placement:

```sh
# Generate integration
$ slot init <shell>
```

The integration enables the `slot run` command which places rendered output
into your shell prompt for editing before execution.

Use `--yes/-y` to execute the rendered command directly without editing.

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

- **Usage:** `slot render <name> [key=value...]`

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
<summary><strong>init</strong> — Generate shell integration snippets</summary>

- **Usage:** `slot init <bash|zsh> [flags]`
- **Flags:**
  - `--fzf` – Enable fzf integration (zsh only, binds to Ctrl-X and Ctrl-Z keys)

</details>

## Templating

Supports basic `text/template` syntax as well as the functions provided by [slim-sprig](https://go-task.github.io/slim-sprig).

## Multiline commands

For multiline commands, use either a `$`:

```sh
# Save 'ls' as a multiline expression
$ slot save ls $'if [ "{{ .INPUT }}" = "true" ]; then\n  echo "{{ .OUTPUT }}";\nfi'
```

or edit the slot directly in the YAML file.
