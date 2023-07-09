# ChatGPT Command-Line Interface (CLI)

ChatGPT CLI offers an efficient way to communicate with ChatGPT directly from the command line.

## Installation

Using brew to install:

```bash
brew tap duanemay/tap
brew install chatgpt-cli
```

## Setting Up

You need to set up a ChatGPT API key to use the CLI. Sign up on [OpenAI platform](https://platform.openai.com/account/api-keys) to generate one. 

It is recommended to create a `.chatgpt-cli` file in your home directory (or the directory where you'll run the CLI), and put your ChatGPT API key there, together with any other default settings you wish to use.

```env
API_KEY=sk-mysupersecretAPIkey
```

## Flags and Environment Variables

The configuration settings precedence order is:

```
command line flags > environment variables > configuration files > defaults
```

Every flag can be set in a config file by converting it to uppercase and replacing any `-` with `_`. For example, `--api-key` can be set as `API_KEY`. The config file is loaded from the current directory, then from the user's home directory. The configuration file is named `.chatgpt-cli`.

You can select a different configuration file with the `--config` flag. Settings in a config file are specified as `KEY=VALUE` pairs, one per line. Lines starting with `#` are treated as comments and ignored.

```
The full list of available flags and corresponding environment variables:

|Flag|Short|Config File Key|Default|Description|
|--|--|--|--|--|
|`--api-key`|`-k`|`API_KEY`|Error|ChatGPT API Key|
|`--config`|`-c`|`CONFIG`|./.chatgpt-cli then $HOME/.chatgpt-cli|Config file to load|
|`--verbose`|`-v`|`VERBOSE`|`false`|Verbose logging|
|`--eom`||`EOM`|`s`|End of message marker|
|`--eos`||`EOS`|`q`|End of session marker|
|`--session-file`|`-s`|`SESSION_FILE`|Generated|Session file|
|`--model`|`-m`|`MODEL`|`gpt-4`|Model to use|
|`--role`| |`ROLE`|`user`|Role of User|
|`--temperature`|`-m`|`TEMPERATURE`|`1.0`|Temperature: 0-2|
|`--max-tokens`||`MAX_TOKENS`|`0`|Max tokens: 8192|
|`--top-p`||`TOP_P`|`1.0`|Top P: 0-1|

For instance, if you want to change the end of message and end of session markers, you can modify them in your configuration file.
```

## Usage

The basic `chatgpt-cli` command usage:

```bash
chatgpt-cli [command]
```

Available commands include:

- `chat`: Enter a chat session with ChatGPT
- `completion`: Generate the autocompletion script for the specified shell
- `help`: Help about any command
- `list-models`: Lists all models available to your account
- `replay-session`: Replay a chat session from a saved file
- `version`: Displays version information

### chat  

This command initiates a chat session with ChatGPT. You will be prompted to enter a message. Type in desired texts which could be multi-line to send. When your message is complete, you can either type CTRL+D or `\s` on a separate line to send. You can change the default end of message marker with the `--eom` flag.

ChatGPT will then respond with a message. You can continue the chat by inputting another message. 

To exit the chat session, you can either type in CTRL+C or `\q` on a separate line. Any text inputted before exiting would not be sent. The default end of session marker can also be changed using the `--eos` flag.

The chat session will be saved in a session file. You can resume this session by specifying the file using the `--session-file` flag.

For instance:
```bash
chatgpt-cli chat --session-file session.json
```

You can also set Model, Role, Temperature, Max Tokens, and Top P with the `--model`, `--role`, `--temperature`, `--max-tokens`, and `--top-p` flags respectively. These can be changed when resuming a session with the `--session-file` flag.

### replay-session

This command replays a chat session from the saved file, a great way to revisit a chat session in a more readable format than the raw JSON.

```bash
chatgpt-cli replay-session --session-file session.json
```

### list-models

This command lists all models available for your account.

```bash
chatgpt-cli list-models
```

### version

This command displays version information.

```bash
chatgpt-cli version
```

## Build and Release

For a trial run, do a dry-run of the deployment process. This can be achieved by running a command with few flags that will hinder you from releasing it on GitHub.

```bash
goreleaser --snapshot --skip-publish --clean
```

After adding the GitHub token in `./.github_token`, run:

```bash
goreleaser --clean
```

##### Examples for use

###### Improving README
```bash
echo "Rewrite this README file as a user guide. Make it easy to read and informative. Use a helpful and clear style" \
  | chatgpt-cli chat < README.md > README-new.md
```

###### Editing a directory full of notes
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
