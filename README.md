# priori

A command-line task manager with priority setting via flags, built for CLI-based workflows.

### Quick Install (Linux & macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/robkokochak/priori/main/install.sh | bash
```

### Manual Installation
Download the appropriate binary for your platform from the [releases page](https://github.com/robkokochak/priori/releases) and place the executable in your path.

## Usage
Run `priori --help` at any time to see the available commands and flags.

Priori writes tasks to a markdown file. This makes it ideal for use with a markdown viewer like Obsidian, which can be kept open alongside your terminal.

### Adding Tasks
Add tasks with different priority levels using flags:
```bash
priori Complete project documentation -m     # Medium priority
priori Review pull requests -h               # High priority
priori Respond to planning epic comments -h  # High priority
priori Update dependencies -l                # Low priority
priori Check emails                          # No priority (will be added to a list with heading '~')
```

### Listing Tasks
To print the tasks to the terminal, run:
```bash
priori list
```
```
### High Priority
- Respond to planning epic comments
- Review pull requests
### Medium Priority
- Complete project documentation
### Low Priority
- Update dependencies
### ~
- Check emails
```

### Deleting Tasks
Delete tasks by their index within their priority section:
```bash
priori delete 0 -h  # Delete 1st high-priority task
priori delete 2 -m  # Delete 3rd medium-priority task
priori delete 1 -l  # Delete 2nd low-priority task
```

### Configuration
Configure the path to your tasks file:
```bash
priori config set-path /path/to/directory  # Set tasks file location
priori config get-path                     # Show current tasks file location
```

### Flags
- `-h, --high`: Set high priority
- `-m, --medium`: Set medium priority
- `-l, --low`: Set low priority
- `--help`: Show help information

## Contributing
Contributions are welcome! Please feel free to reach out, submit a PR, or open an issue.
