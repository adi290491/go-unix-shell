# GoShell - A Feature-Rich Unix Shell in Go

A fully-featured, POSIX-compliant Unix shell implementation built from scratch in Go. This project demonstrates advanced system programming concepts including process management, I/O redirection, pipeline implementation, and interactive command-line features.

## 🎯 Features

### Core Shell Functionality
- **Interactive REPL** - Read-Eval-Print Loop with custom prompt
- **Command Execution** - Execute both built-in and external commands
- **Error Handling** - Graceful handling of invalid commands and edge cases

### Built-in Commands
- `pwd` - Print working directory
- `cd` - Change directory (supports absolute paths, relative paths, and `~` for home)
- `echo` - Display text with argument expansion
- `type` - Identify command types (builtin vs external)
- `exit` - Exit the shell with optional exit code
- `history` - View and manage command history

### Advanced Input Handling
- **Quote Processing**
  - Single quotes (`'...'`) - Literal string interpretation
  - Double quotes (`"..."`) - Variable expansion with escape sequences
  - Backslash escaping both inside and outside quotes
  
- **Tab Completion**
  - Built-in command completion
  - External executable completion from `$PATH`
  - Argument completion
  - Partial match completion
  - Multiple completion suggestions

### I/O Redirection
- `>` - Redirect stdout to file (overwrite)
- `>>` - Append stdout to file
- `2>` - Redirect stderr to file (overwrite)
- `2>>` - Append stderr to file

### Pipelines
- Multi-command pipelines with `|` operator
- Support for built-in commands in pipelines
- Proper process synchronization and error handling
- Example: `ls /tmp | grep test | wc -l`

### Command History
- **Interactive Navigation**
  - ↑ / ↓ arrows to navigate history
  - Ctrl+R for reverse search
  
- **History Management**
  - `history` - Display all commands
  - `history N` - Display last N commands
  - `history -r <file>` - Read history from file
  - `history -w <file>` - Write history to file
  - `history -a <file>` - Append new commands to file
  
- **Persistent History**
  - Automatic loading from `$HISTFILE` on startup
  - Automatic saving on exit
  - Append-only mode to preserve history across sessions

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher
- Unix-like operating system (Linux, macOS, WSL)

### Build from Source
```bash