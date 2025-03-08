# (Heist Tech Solutions) HTS Agent

A Go application that interacts with various LLM providers.

## Supported LLM Providers

- **Ollama**: Local LLM provider
- **Grok**: X.AI's Grok API

## Setup

### Prerequisites

- Go 1.24 or higher
- Ollama (for local LLM usage)
- Grok API key (for using Grok)

### Installation

1. Clone the repository
2. Install dependencies:

```bash
go get
```

### Configuration

The application uses environment variables for configuration. You can set these in a `.env` file in the root directory.

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Edit the `.env` file and add your API keys:

```
# API Keys for LLM Providers
GROK_API_KEY=your_grok_api_key_here
```

## Usage

### Running the Application

```bash
go run main.go
```

By default, the application uses the Ollama provider. To use the Grok provider, uncomment the Grok code in `main.go` and ensure you have set your Grok API key in the `.env` file.

### Testing

Run the tests:

```bash
go test ./...
```

## Development

### Adding a New Provider

To add a new LLM provider:

1. Create a new file in the `llm_provider` directory
2. Implement the provider's API client
3. Add tests for the new provider
4. Update the `main.go` file to use the new provider
5. Add any required API keys to the `.env.example` file and update the `config` package

## License

[MIT](LICENSE)
