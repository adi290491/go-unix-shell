[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)]()

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
# Clone the repository
git clone https://github.com/yourusername/goshell.git
cd goshell

# Build the binary
go build -o goshell app/*.go

# Run the shell
./goshell
```

### Quick Start with Script
```bash
# Make the run script executable
chmod +x run.sh

# Run the shell
./run.sh
```

## 📖 Usage

### Basic Commands
```bash
$ pwd
/home/user

$ cd /tmp
$ pwd
/tmp

$ cd ~
$ pwd
/home/user

$ echo "Hello, World!"
Hello, World!

$ type echo
echo is a shell builtin

$ type ls
ls is /usr/bin/ls
```

### Quoting Examples
```bash
# Single quotes - literal interpretation
$ echo 'hello   world'
hello   world

# Double quotes - preserves spacing
$ echo "hello   world"
hello   world

# Escaping
$ echo "It's a \"beautiful\" day"
It's a "beautiful" day

$ echo 'backslash: \n remains literal'
backslash: \n remains literal
```

### Redirection
```bash
# Redirect stdout
$ echo "log entry" > output.txt

# Append to file
$ echo "another entry" >> output.txt

# Redirect stderr
$ ls nonexistent 2> errors.log

# Append stderr
$ ls another_nonexistent 2>> errors.log
```

### Pipelines
```bash
# Simple pipeline
$ ls /usr/bin | wc -l

# Multi-stage pipeline
$ cat access.log | grep "ERROR" | wc -l

# Pipeline with built-ins
$ echo -e "banana\napple\ncherry" | cat | grep apple
```

### History
```bash
# View all history
$ history
    1  pwd
    2  cd /tmp
    3  ls -la
    4  history

# View last 5 commands
$ history 5

# Navigate with arrows
$ ↑    # Previous command
$ ↓    # Next command

# Save history to file
$ history -w backup.txt

# Load history from file
$ history -r backup.txt

# Append new commands to file
$ history -a backup.txt
```

### Tab Completion
```bash
# Complete built-in commands
$ ec<TAB>    # Completes to "echo"

# Complete executables
$ gre<TAB>   # Suggests "grep"

# Multiple matches
$ py<TAB>    # Shows: python python3 pydoc
```

## 🏗️ Project Structure
```
goshell/
├── app/
│   ├── main.go              # Entry point and REPL loop
│   ├── builtins.go          # Built-in command implementations
│   ├── external.go          # External command execution
│   ├── pipeline.go          # Pipeline parsing and execution
│   ├── history.go           # History management and persistence
│   ├── completion.go        # Tab completion logic
│   ├── redirect.go          # I/O redirection handling
│   ├── parser.go            # Command parsing and quote handling
│   └── utils.go             # Utility functions
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── run.sh                   # Convenience script to build and run
├── README.md                # This file
└── .gitignore              # Git ignore rules
```

## 🔧 Configuration

### Environment Variables

- `HISTFILE` - Path to history file (default: `/tmp/readline.tmp`)
- `PATH` - Executable search paths

### Example Configuration
```bash
# Set custom history file
export HISTFILE=~/.goshell_history

# Run the shell
./goshell
```

## 🧪 Testing
```bash
# Run all tests
go test ./app/...

# Run with verbose output
go test -v ./app/...

# Run specific test
go test -run TestPipeline ./app/
```

## 🎓 Technical Highlights

### Advanced Concepts Implemented

1. **Process Management**
   - Fork/exec model for external commands
   - Proper parent-child process synchronization
   - Signal handling and cleanup

2. **I/O Multiplexing**
   - Pipe creation and management
   - File descriptor handling
   - Buffered I/O for performance

3. **Parser Design**
   - Lexical analysis for quote handling
   - Command tokenization
   - Argument parsing with escape sequences

4. **Concurrent Execution**
   - Pipeline command parallelization
   - Goroutines for async I/O
   - Channel-based error propagation

5. **Interactive Features**
   - GNU Readline-like functionality
   - History management with deduplication
   - Smart tab completion with prefix matching

## 🐛 Known Limitations

- No job control (`&`, `fg`, `bg`)
- No command substitution (`$(...)` or `` `...` ``)
- No variable expansion (`$VAR`)
- No globbing (`*.txt`)
- No aliases
- Single-threaded command execution (no `&` backgrounding)

## 🤝 Contributing

This is a learning project, but suggestions and improvements are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/improvement`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/improvement`)
5. Open a Pull Request

## 📚 Learning Resources

This project was built as part of the [CodeCrafters](https://codecrafters.io) "Build Your Own Shell" challenge. Great resources for understanding shell internals:

- [Advanced Programming in the UNIX Environment](https://www.apuebook.com/)
- [The Linux Programming Interface](https://man7.org/tlpi/)
- [GNU Bash Manual](https://www.gnu.org/software/bash/manual/)
- [POSIX Shell Command Language](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html)

## 📄 License

MIT License - feel free to use this code for learning purposes.

## 🙏 Acknowledgments

- Built as part of [CodeCrafters](https://codecrafters.io) challenge
- Inspired by Bash, Zsh, and other Unix shells
- Uses [readline](https://github.com/chzyer/readline) library for line editing

## 👤 Author

**Your Name**
- GitHub: [@yourusername](https://github.com/yourusername)
- LinkedIn: [your-profile](https://linkedin.com/in/yourprofile)

---

⭐ Star this repo if you found it helpful!