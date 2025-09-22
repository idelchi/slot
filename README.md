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

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/slot/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# Save a command with Go template variables
$ slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod --description 'deploy to prod'
```

```sh
# Render a command with variable substitution
$ slot render deploy file=k8s.yml
kubectl apply -f k8s.yml
```

```sh
# List all saved slots
$ slot list
NAME    CMD                         TAGS      DESCRIPTION
deploy  kubectl apply -f {{.file}}  k8s,prod  deploy to prod
```

```sh
# List slots filtered by tags
$ slot list --tags k8s
```

```sh
# Remove a slot
$ slot remove deploy
```

## Data Storage

Slots are stored in YAML format at `~/.config/slot/slots.yaml`. Location can be overridden with the
`--config` flag or `SLOTS_FILE` environment variable.

## Shell Integration

Generate shell integration snippets for command placement:

```sh
# Generate integration
$ slot init <shell>
```

The integration enables the `slot run` command which places rendered output
into your shell prompt for editing before execution.

Use `--yes/-y` to execute the rendered command directly without editing.

Adding the `--fzf` flag enables further integration, binding `Ctrl-X` and `Ctrl-Z` keys to running or searching slots.

## Commands

<details>
<summary><strong>save</strong> — Save a command slot</summary>

- **Usage:** `slot save <name> <command> [flags]`
- **Flags:**
  - `--tags` – Tags for the slot (repeatable)
  - `--description` – Description for the slot
  - `--force` – Overwrite existing slot

</details>

<details>
<summary><strong>render</strong> — Render a saved command slot</summary>

- **Usage:** `slot render <name> [key=value...]`

</details>

<details>
<summary><strong>list/ls</strong> — List saved slots</summary>

- **Usage:** `slot list [flags]`
- **Flags:**
  - `--tags` – Filter by tags (repeatable)
  - `--tsv` – Output in TSV format

  </details>

<details>
<summary><strong>remove/rm</strong> — Delete a saved slot</summary>

- **Usage:** `slot remove <name>`

</details>

<details>
<summary><strong>init</strong> — Generate shell integration snippets</summary>

- **Usage:** `slot init <bash|zsh> [flags]`
- **Flags:**
  - `--fzf` – Enable fzf integration (binds to Ctrl-X and Ctrl-Z keys)

</details>

## Templating

In addition to your own `key=value` arguments, the following variables are always available inside templates:

- **`SLOTS_FILE`** – Full path to the slots YAML file
- **`SLOTS_DIR`** – Directory containing the slots YAML file
- **`CLI_ARGS`** – All arguments after `--`, joined into a single space-delimited string
- **`CLI_ARGS_SPLIT`** – Same arguments as above, but preserved as a list (`[]string`) for iteration

All templates use Go’s [`text/template`](https://pkg.go.dev/text/template) syntax, with extra functions from [slim-sprig](https://go-task.github.io/slim-sprig).

## Demo

![Demo](assets/gifs/slot.gif)
