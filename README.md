# priori

A command-line task manager with priority setting via flags, built for CLI-based workflows.

### Quick Install (Linux & macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/robkokochak/priori/main/install.sh | bash
```

### Manual Installation
Download the appropriate binary for your platform from the [releases page](https://github.com/robkokochak/priori/releases) and place the executable in your path.

## Usage

Priori writes tasks to a markdown file. This makes it ideal for use with a markdown viewer like Obsidian, which can be kept open alongside your terminal.

### Adding Tasks
Add tasks with different priority levels using flags:
```bash
priori Complete project documentation -h  # High priority
priori Review pull requests -m           # Medium priority
priori Update dependencies -l            # Low priority
priori Check emails                      # No priority (will be added to a list with heading '~')
```

### Listing Tasks
View all tasks organized by priority:
```bash
priori list
```

### Deleting Tasks
Delete tasks by their index within their priority section:
```bash
priori delete 0 -h  # Delete 1st high-priority task
priori delete 2 -m  # Delete 3rd medium-priority task
priori delete 1 -l  # Delete 2nd low-priority task
```

### Configuration
Configure the storage location for your tasks:
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
Contributions are welcome! Please feel free to submit a PR or open an issue.
