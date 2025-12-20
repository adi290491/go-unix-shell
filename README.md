# Go-Unix-Shell

A fully-featured Unix-like shell implemented in Go.

## Features
- Builtin commands: cd, pwd, history, type, exit
- Quoting (single, double, backslash rules)
- Redirection (stdout, stderr, append)
- Pipelines (multi-command, builtins + external)
- Autocompletion
- Persistent history (HISTFILE support)

## Motivation
Built as part of the CodeCrafters “Build Your Own Shell” challenge,
extended to complete **all stages and extensions**.

## How to Run
```bash
go build -o shell app/*.go
./shell
