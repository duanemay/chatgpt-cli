# ChatGPT-cli

A way to communicate with ChatGPT from the command line.

## Installing 

```bash
brew tap duanemay/tap
brew install chatgpt-cli
```

## Running

### API Key
To utilize the ChatGPT CLI, a ChatGPT API key is required. You can obtain one by signing up and creating an API key at [API keys](https://platform.openai.com/account/api-keys).

It is recommended to create a file named: `.chatgpt-cli` in your home directory (or a directory you plan to run the cli). This file should contain your ChatGPT API key, and any other default settings you wish to use.

```env
API_KEY=sk-mysupersecretAPIkey
```

### Flags, Environment Variables

Configuration settings follow the precedence order: 
```
command line flags > environment variables > configuration files > defaults
```

A value specified in a configuration file overrides the default value.
A value specified as an environment variable overrides the value in a configuration file.
A value specified on the command line overrides all other settings. 


All flags can be set in a configuration file, by uppercasing and converting any `-`s to `_`s. For example, `--api-key` can be set as `API_KEY`.  The configuration file is loaded from the current directory, then the user's home directory. The configuration file is named `.chatgpt-cli`.
You can specify a different configuration file with the `--config` flag. Settings in a configuration file are specified as `KEY=VALUE` pairs, one per line. Lines beginning with `#` are treated as comments and ignored. 

All flags can also be set as environment variables by prefixing with `CHATGPT_`, uppercasing, and converting `-` to `_`. For example, `--api-key` can be set as `CHATGPT_API_KEY`. 


|  Long            | Short | Config File Key | Default                                | Description           |
|------------------|-------|-----------------|----------------------------------------|-----------------------|
| `--api-key`      | `-k`  | `API_KEY`       | Error                                  | ChatGPT API Key       |
| `--config`       | `-c`  | `CONFIG`        | ./.chatgpt-cli then $HOME/.chatgpt-cli | Config file to load   |
| `--verbose`      | `-v`  | `VERBOSE`       | `false`                                | Verbose logging       |
| `--eom`          |       | `EOM`           | `\s`                                   | End of message marker |
| `--eos`          |       | `EOS`           | `\q`                                   | End of session marker |
| `--session-file` | `-s`  | `SESSION_FILE`  | Generated                              | Session file          |
| `--model`        | `-m`  | `MODEL`         | `gpt-4`                                | Model to use          |
| `--role`         |       | `ROLE`          | `user`                                 | Role of User          |
| `--temperature`  | `-t`  | `TEMPERATURE`   | `1.0`                                  | Temperature: 0-2      |
| `--max-tokens`   |       | `MAX_TOKENS`    | `0`                                    | Max tokens: 8192      |
| `--top-p`        |       | `TOP_P`         | `1.0`                                  | Top P: 0-1            |

Say you don't like `\s` and `\q` for the end of message and end of session markers. You can set them in your configuration file.

## Usage

```bash
chatgpt-cli is a CLI for ChatGPT

Usage:
  chatgpt-cli [command]

Available Commands:
  chat          Enter a chat session with ChatGPT
  completion    Generate the autocompletion script for the specified shell
  help          Help about any command
  list-models    lists all models available to your account
  replay-session Replay a chat session from saved file
  version       displays version information

Flags:
  -k, --api-key string   ChatGPT apiKey
  -c, --config string    Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)
  -h, --help             help for chatgpt-cli
  -v, --verbose          verbose logging

Use "chatgpt-cli [command] --help" for more information about a command.

```

## Commands

### chat

Enter a chat session with ChatGPT

```bash
chatgpt-cli chat 
```

Help for chat command
```
Enter a chat session with ChatGPT

Usage:
  chatgpt-cli chat [flags]

Flags:
      --eom string            Text to enter to mark the end of a message to send to ChatGPT (default "\s")
      --eos string            Text to enter to end of a session with ChatGPT (default "\q")
  -h, --help                  help for chat
      --max-tokens int        number of tokens to generate = $ (default 9223372036854775807)
  -m, --model string          ChatGPT Model (default "gpt-4")
  -r, --role string           ChatGPT Role (default "user")
  -s, --session-file string   Continue a session from a file
  -t, --temperature float32   temperature, between 0 and 2. Higher values make the output more random (default 1)
      --top-p float32         results of the tokens with top_p probability mass (default 1)

Global Flags:
  -k, --api-key string   ChatGPT apiKey
  -c, --config string    Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)
  -v, --verbose          verbose logging
```

This will enter a chat session with ChatGPT. You will be prompted for a message. Enter a message you can enter multi-line text to send. When
your message is complete enter CTRL+D or `\s` on its own line to send. (You can change the default end of message marker with the `--eom` flag)

ChatGPT will respond with a message. You can continue the conversation by entering another message. 

To exit the chat session, enter CTRL+C or `\q` on its own line. Any text already entered will not be sent. (You can change the default end of session marker with the `--eos` flag)

Your session will be saved in session file. The file name will be displayed when you enter the chat session.
You can resume the session by specifying the file using the `--session-file` flag.

```bash
chatgpt-cli chat --session-file session.json
```

Model, Role, Temperature, Max Tokens, Top P can be set with the `--model`, `--role`, `--temperature`, `--max-tokens`, `--top-p` flags.
These can be changed when resuming a session with the `--session-file` flag.

### replay-session

Replay a chat session from saved file, this is useful for displaying a chat session in a easier to read format, than the raw JSON.

```bash
chatgpt-cli replay-session --session-file session.json
```

Help for replay-session command
```
Usage:
  chatgpt-cli replay-session [flags]

Flags:
  -h, --help                 help for replaySession
  -s, --session-file string   File to replay a Session from

Global Flags:
  -k, --api-key string   ChatGPT apiKey
  -c, --config string    Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)
  -v, --verbose          verbose logging
```

### list-models  

lists all models available to your account

```bash
chatgpt-cli list-models
```

Help for list-models command
```
 lists all models available to your account

Usage:
  chatgpt-cli list-models [flags]

Flags:
  -h, --help   help for list-models

Global Flags:
  -k, --api-key string   ChatGPT apiKey
  -c, --config string    Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)
  -v, --verbose          verbose logging
```

### version

displays version information

```bash
chatgpt-cli version
```

## Build & Release

Do a trial run of the deployment process. This is as easy as running a command with a few flags that will stop you from releasing to GitHub.

```bash
goreleaser --snapshot --skip-publish --clean
```

After adding the GitHub token in `./.github_token`
Run:
```bash
goreleaser --clean
```
