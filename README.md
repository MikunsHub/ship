# Ship CLI

A lightweight CLI tool that streamlines your Git and GitHub workflows. Ship automates feature branch creation and pull request management with AI-powered descriptions using Google's Gemini API.

## Features

- **Feature Branch Creation** - Quickly start new features from the main branch
- **Automated PR Creation** - Create pull requests to multiple branches (main, stage, dev)
- **AI-Powered Descriptions** - Generate professional PR descriptions using Google Gemini
- **Secure API Key Management** - Store API keys securely in your system keyring
- **Fast and Reliable** - Built in Go for performance and cross-platform compatibility

## Installation

### Quick Install (Recommended)

With Go 1.24+:

```bash
go install github.com/MikunsHub/ship@latest
```

### Specific Version

```bash
go install github.com/MikunsHub/ship@v0.1.0
```

### From Source

```bash
git clone https://github.com/MikunsHub/ship.git
cd ship
make install
```

### Verify Installation

```bash
ship --version
ship -h
```

## Usage

### Start a New Feature

```bash
ship -f <branch-name>
```

Creates a new feature branch from main with the specified name:

```bash
ship -f mikun/authentication
```

The tool will:
1. Checkout the main branch
2. Pull the latest changes
3. Create and switch to your feature branch
4. Push to origin with upstream tracking

### Create Pull Requests

Create PRs to all branches (main, stage, dev):

```bash
ship prs
```

Create PR from current branch to a specific branch:

```bash
ship prs -s stage
```

Create PR from a specific branch:

```bash
ship prs <branch-name>
```

The tool will:
1. Fetch commits between the branches
2. Generate an AI-powered PR description (if API key is configured)
3. Show you the description and ask for confirmation
4. Allow you to accept, reject, or edit the description
5. Create the PR on GitHub

### Configure Your API Key

Store your Google Gemini API key securely:

```bash
ship config set-key
```

Check configuration status:

```bash
ship config status
```

Remove your stored API key:

```bash
ship config remove-key
```

### Get Help

```bash
ship -h
ship help
```

## Configuration

### API Key Setup

To use AI-powered PR descriptions, you'll need a Google Gemini API key:

1. Get your free API key at [Google AI Studio](https://aistudio.google.com/apikey)
2. Run `ship config set-key` and enter your key securely
3. The key is stored in your system keyring (Keychain on macOS)

Alternatively, set the environment variable:

```bash
export GEMINI_API_KEY=your-api-key-here
```

## Development

### Build Locally

```bash
make build      # Creates bin/ship
make install    # Builds and installs to $GOPATH/bin
make clean      # Remove build artifacts
make version    # Show current version
```

### Project Structure

```
ship/
├── main.go              # Core CLI logic and commands
├── config.go            # Configuration and API key management
├── keyring.go           # Secure credential storage
├── pr_body_prompt.go    # AI prompt template
├── go.mod               # Go module definition
├── Makefile             # Build automation
└── .github/workflows/   # GitHub Actions CI/CD
```

## Releasing a New Version

Ship uses semantic versioning (v0.1.0, v0.2.0, etc.) and GitHub Actions for automated releases.

### Create a Release

1. Ensure all changes are committed and pushed:

```bash
git add .
git commit -m "feat: add new feature"
git push
```

2. Create and push a semantic version tag:

```bash
# For new features (v0.1.0 → v0.2.0)
git tag -a v0.2.0 -m "Release v0.2.0: Add new features"
git push --tags

# For bug fixes (v0.1.0 → v0.1.1)
git tag -a v0.1.1 -m "Release v0.1.1: Fix issue"
git push --tags

# For major changes (v0.1.0 → v1.0.0)
git tag -a v1.0.0 -m "Release v1.0.0: Stable API"
git push --tags
```

3. GitHub Actions automatically:
   - Detects the tag
   - Builds the binary
   - Creates a GitHub Release with the binary

4. Users can then install the new version:

```bash
go install github.com/MikunsHub/ship@v0.2.0
```

## Versioning

Ship follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1.0.0) - Breaking changes
- **MINOR** (v0.2.0) - New features (backward compatible)
- **PATCH** (v0.1.1) - Bug fixes (backward compatible)

## Requirements

- Go 1.24+
- Git
- GitHub CLI (`gh`)
- Google Gemini API key (optional, for AI descriptions)

## License

MIT

## Contributing

Contributions are welcome! Feel free to open issues and pull requests.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
