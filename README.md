# User Guide - ChatGPT Command-Line Interface (CLI)

The ChatGPT CLI allows you to interact with ChatGPT directly from your command line, offering an efficient platform for real-time communication. This user guide will provide you with simple, straightforward instructions on how to install, set up, and effectively use the CLI.

## Examples

### Simple Interactive Chat

![Simple chat](docs/translation-demo.gif)

### Improving README, non-Interactively

```bash
echo "Rewrite this README file as a user guide. Make it easy to read and informative. Use a helpful and clear style" \
  | chatgpt-cli chat < README.md > README-new.md
```

### Editing a directory full of notes, non-Interactively
```bash
for file in notes/*.md; do
  printf "\nEditing '%s'\n============================\n" "$file"
  {
    printf "Revise my notes. "
    printf "Use Markdown format. "
    printf "Revise the text of the note below to use a clear and informative style. "
    printf "Use newlines to keep line length less than 120 characters.\n\n"
  } >"${file}.req"

  chatgpt-cli chat <"${file}.req" >"${file}.new"
  if [ $? -ne 0 ]; then
    printf "  Request Failed: '%s'\n" "${file}"
    rm "${file}.new"
  else
    mv "${file}.new" "${file}"
  fi
  rm "${file}.req"
done
```

## Installation

To install ChatGPT CLI using Homebrew:

```bash
brew tap duanemay/tap
brew install chatgpt-cli
```

## Initial Set Up

To use the ChatGPT CLI, you'll need a ChatGPT API key. You can generate this key by signing up on the [OpenAI platform](https://platform.openai.com/account/api-keys).

It is recommended to create a `.chatgpt-cli` file in your home directory or the directory where you'll run the CLI. Inside this file, include your API key alongside additional default settings you wish to use.

Example:

```env
API_KEY=sk-mysupersecretAPIkey
```

## Configuration Settings

The priority order for the configuration settings is as follows:

```
Command line flags > Environment Variables > Configuration files > Defaults
```

Each flag can be set in a configuration file, by changing it to uppercase and replacing `-` with `_`.

The full list of available flags and corresponding environment variables:

|Flag| Short |Config File Key| Default                                | Description                         |
|--|-------|--|----------------------------------------|-------------------------------------|
|`--api-key`| `-k`  |`API_KEY`| **Required**                           | ChatGPT API Key                     |
|`--config`| `-c`  |`CONFIG`| ./.chatgpt-cli then $HOME/.chatgpt-cli | Config file to load                 |
|`--verbose`| `-v`  |`VERBOSE`| `false`                                | Verbose logging                     |
|`--eom`|       |`EOM`| `s`                                    | End of message marker               |
|`--eos`|       |`EOS`| `q`                                    | End of session marker               |
|`--session-file`| `-s`  |`SESSION_FILE`| Generated                              | Session file                        |
|`--no-write-session`|       |`NO_WRITE_SESSION`| false                                  | Do not write or update session file |
|`--model`| `-m`  |`MODEL`| `gpt-4`                                | Model to use                        |
|`--role`|       |`ROLE`| `user`                                 | Role of User                        |
|`--temperature`| `-m`  |`TEMPERATURE`| `1.0`                                  | Temperature: 0-2                    |
|`--max-tokens`|       |`MAX_TOKENS`| `0`                                    | Max tokens: 8192                    |
|`--top-p`|       |`TOP_P`| `1.0`                                  | Top P: 0-1                          |

For instance, if you want to change the end of the message and session markers, modify them in your configuration file.

You can select a different configuration file using `--config` flag. Each config file should specify settings as `KEY=VALUE` pairs, with each pair on a separate line. Lines commencing with `#` are considered comments and ignored.

## Usage

The basic command to use `chatgpt-cli` is as follows:

```bash
chatgpt-cli [command]
```

The available commands are as follows:

1. `chat`: Start a chat session with ChatGPT.
2. `completion`: Generate the autocomplete script for your chosen shell.
3. `help`: Seek help regarding any command.
4. `list-models`: Retrieve a list of all models available to your account.
5. `replay-session`: Replay a chat session from a previously saved file.
6. `version`: Get version information.

### Chatting

Initiate a chat session with ChatGPT using the `chat` command:

```bash
chatgpt-cli chat
```

You'll be prompted to input your message, which can span multiple lines. Send your message with CTRL+D or by entering `\s` as a separate line. The default end of the message marker can be changed with the `--eom` flag.

Continue your conversation with ChatGPT by inputting a new message once you receive a response.

Exiting the chat is made possible by inputting CTRL+C or `\q`. Note that any text entered before exiting won't be sent. The default end of session marker can be changed with the `--eos` flag.

All chat sessions are saved in a session file, for which the `--session-file` flag can specify the file of your choice:

```bash
chatgpt-cli chat --session-file session.json
```

The `--model`, `--role`, `--temperature`, `--max-tokens`, and `--top-p` flags allow for individual configuration of the Model, Role, Temperature, Max Tokens, and Top P respectively.

### Replaying a Session

Replaying a chat session lets you revisit a previous chat in a more readable format than the raw JSON. Use the `replay-session` command:

```bash
chatgpt-cli replay-session --session-file session.json
```

### Listing Models

`list-models` retrieves a list of all available models related to your account:

```bash
chatgpt-cli list-models
```

### Checking the Version

Fetch information about the CLI's version using `version`:

```bash
chatgpt-cli version
```

## Build and Release

For a dry-run of the deployment process, you can run the following:

```bash
goreleaser --snapshot --skip-publish --clean
```

For a full deployment, you'll need to add a GitHub token to the `./.github_token` file, then run the following:

```bash
goreleaser --clean
brew upgrade
```
