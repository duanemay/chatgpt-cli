# User Guide - ChatGPT Command-Line Interface (CLI)

The ChatGPT CLI allows you to interact with ChatGPT directly from your command line, offering an efficient platform for real-time communication. This user guide will provide you with simple, straightforward instructions on how to install, set up, and effectively use the CLI.

## Examples

### Simple Interactive Chat

![Translation Demo](docs/translation-demo.gif)

### Ask a Question, non-Interactively

```bash
{ echo "What is the largest file in this directory?"; ls -l } | chatgpt-cli chat
```

### Improving README, non-Interactively

```bash
echo "Rewrite this README file as a user guide. Make it easy to read and informative. Use a helpful and clear style" \
  | chatgpt-cli chat < README.md > README-new.md
mv README-new.md README.md
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
    cat "${file}"
  } | chatgpt-cli chat > "${file}.new"

  if [ $? -ne 0 ]; then
    printf "  Request Failed: '%s'\n" "${file}"
    rm "${file}.new"
  else
    mv "${file}.new" "${file}"
  fi
done
```

### Image Generation and Vision Chat

![Image Demo](docs/image-demo.gif)

Which gives us the resulting image.

![Result for requesting an image “white maltipoo dressed in a shark costume”](docs/maltipoo-01.png)

### Generate an Image, non-Interactively

```bash
echo "Monkey in a banana costume, on a high wire above a circus" | chatgpt-cli image -o monkey
```

In non-Interactive mode only the name of the output files are sent to stdout.
In this case, monkey-01.png, shown here.

![Result for requesting an image “Monkey in a banana costume, on a high wire above a circus”](docs/monkey-01.png)


```bash
echo "What is in this image?" | chatgpt-cli vision --file monkey-01.png
```

This will upload the image to ChatGPT for use in the chat, and use it to answer the question.
You might get a response like this:

```
The image depicts a whimsical circus scene with a monkey in a banana suit balancing on a tightrope. The atmosphere is lively, with a large audience in tents watching the performance. The setting has bright lights and colorful decorations, capturing the festive spirit of a circus.                                                                                                                                              
```

### Generate an Audio file from Text (Text to Speech)

![Speech Demo](docs/speech-demo.gif)

Which gives us the resulting audio file.


https://github.com/duanemay/chatgpt-cli/assets/991170/bb629736-3dc1-4ee1-99f8-1a93646c60da

And similarly, you can do this non-Interactively

```bash
echo "Hello, I am a robot. I am here to help you with your tasks." | chatgpt-cli speech
```

## Installation

To install ChatGPT CLI using Homebrew:

```bash
brew tap duanemay/tap
brew install chatgpt-cli
```

### Generating Cover Images for a directory full of notes, non-Interactively
```bash
for file in notes/*.md; do
  dir=$(dirname "${file}")
  filename=$(basename "${file}")
  base_filename=${filename%.*}

  printf "\nCreating Images '%s'\n============================\n" "$file"
  {
    printf "Describe the contents of an image that would make a good cover image for the blog post below.\n\n"
    cat "${file}"
  } | chatgpt-cli chat --system-message "As an expert creator of blog posts" > "${dir}/${base_filename}-img-description.txt"

  if [ $? -ne 0 ]; then
    printf "  Request Failed: '%s'\n" "${file}"
    rm "${dir}/${base_filename}-img-description.txt"
    continue
  fi

  chatgpt_file=$( ls chatgpt-cli* | tail -1 )
  grep -A 1 user "${chatgpt_file}" | grep '"content":' | cut -d':' -f 2 | sed 's/"//g'
  cat "${dir}/${base_filename}-img-description.txt" | chatgpt-cli image -o "${dir}/${base_filename}-img"
done
```

## Initial Set Up

To use the ChatGPT CLI, you'll need a ChatGPT API key. You can generate this key by signing up on the [OpenAI platform](https://platform.openai.com/account/api-keys). ChatGPT CLI now supports Project API keys which have replaced user API keys. In addition, you will need to add funds to your account.

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
Environment variables are the same as the configuration file flag, with the addition of the prefix `CHATGPT_`.

As an example, the `--api-key sk` command line flag can be set in a configuration file as `API_KEY=sk`, and as an environment variable as `CHATGPT_API_KEY=sk`.

The full list of available flags and corresponding config file variables:

*Common Flags:*

| Flag         | Short | Config File Key  | Default                                | Description          |
|--------------|-------|------------------|----------------------------------------|----------------------|
| `--api-key`  | `-k`  | `API_KEY`        | **Required**                           | ChatGPT API Key      |
| `--config`   | `-c`  | `CONFIG`         | ./.chatgpt-cli then $HOME/.chatgpt-cli | Config file to load  |
| `--verbose`  | `-v`  | `VERBOSE`        | `false`                                | Verbose logging      |

*Chat Flags:*

| Flag                   | Short | Config File Key      | Default       | Description                            |
|------------------------|-------|----------------------|---------------|----------------------------------------|
| `--system-message`     |       |                      | ``            | Initial System message sent to ChatGPT |
| `--session-file`       | `-s`  | `SESSION_FILE`       | Generated     | Session file                           |
| `--skip-write-session` |       | `SKIP_WRITE_SESSION` | false         | Do not write or update session file    |
| `--model`              | `-m`  | `MODEL`              | `gpt-4o-mini` | Model to use (default will change)     |
| `--role`               |       | `ROLE`               | `user`        | Role of User                           |
| `--temperature`        | `-m`  | `TEMPERATURE`        | `1.0`         | Temperature: 0-2                       |
| `--max-tokens`         |       | `MAX_TOKENS`         | `0`           | Max tokens: 8192                       |
| `--top-p`              |       | `TOP_P`              | `1.0`         | Top P: 0-1                             |

*Image Flags:*

| Flag              | Short | Config File Key | Default       | Description                  |
|-------------------|-------|-----------------|---------------|------------------------------|
| `--model`         | `-m`  | `MODEL`         | `gpt-image-1` | Model to use                 |
| `--number`        | `-n`  |                 | `1`           | Number of images to generate |
| `--quality`       |       | `QUALITY`       | `high`        | Image Quality                |
| `--size`          | `-s`  | `SIZE`          | 1024x1024     | Image Size                   |
| `--style`         |       | `STYLE`         | `vivid`       | Image Style                  |
| `--output-prefix` | `-o`  | `OUTPUT_PREFIX` | Generated     | File Name Prefix             |

*Speech Flags:*

| Flag              | Short | Config File Key | Default   | Description                 |
|-------------------|-------|-----------------|-----------|-----------------------------|
| `--model`         | `-m`  | `MODEL`         | `tts-1`   | Text to Speech Model to use |
| `--speed`         | `-s`  |                 | `1`       | Speed of Audio              |
| `--voice`         |       | `VOICE`         | `alloy`   | Voice Used                  |
| `--output-prefix` | `-o`  | `OUTPUT_PREFIX` | Generated | File Name Prefix            |

*Vision Flags:*

| Flag                   | Short | Config File Key      | Default       | Description                            |
|------------------------|-------|----------------------|---------------|----------------------------------------|
| `--file`               | `-f`  | `FILE`               | required      | Input image files                      |
| `--model`              | `-m`  | `MODEL`              | `gpt-4o-mini` | Model to use (default will change)     |
| `--system-message`     |       |                      | ``            | Initial System message sent to ChatGPT |
| `--session-file`       | `-s`  | `SESSION_FILE`       | Generated     | Session file                           |
| `--skip-write-session` |       | `SKIP_WRITE_SESSION` | false         | Do not write or update session file    |
| `--role`               |       | `ROLE`               | `user`        | Role of User                           |

*Transcription Flags:*

| Flag               | Short | Config File Key | Default     | Description                    |
|--------------------|-------|-----------------|-------------|--------------------------------|
| `--file`           | `-f`  | `FILE`          | required    | Input audio files              |
| `--model`          | `-m`  | `MODEL`         | `whisper-1` | Transcription Model to use     |
| `--language`       | `-l`  | `LANGUAGE`      | detected    | Language of the input audio    |
| `--system-message` |       |                 | ``          | Initial Prompt sent to ChatGPT |

For instance, if you want to change the end of the message and session markers, modify them in your configuration file.

You can select a different configuration file using `--config` flag. Each config file should specify settings as `KEY=VALUE` pairs, with each pair on a separate line. Lines beginning with `#` are considered comments and ignored.

## Usage

The basic command to use `chatgpt-cli` is as follows:

```bash
chatgpt-cli [command]
```

The available commands are as follows:

1. `chat`: Start a chat session with ChatGPT.
2. `vision`: Upload an image to ChatGPT for use in chat.
3. `image`: Generate an image
4. `speech`: Generate speech using ChatGPT
4. `transcribe`: Transcribe audio to text using ChatGPT
5. `completion`: Generate the autocomplete script for your chosen shell.
6. `help`: Seek help regarding any command.
7. `list-models`: Retrieve a list of all models available to your account.
8. `replay-session`: Replay a chat session from a previously saved file.
7`version`: Get version information.

### Chatting

Initiate a chat session with ChatGPT using the `chat` command:

```bash
chatgpt-cli chat
```

You'll be prompted to input your message, which can span multiple lines. Send your message with TAB or CTRL+C.

Continue your conversation with ChatGPT by inputting a new message once you receive a response.

Exiting the chat is made possible by inputting CTRL+C or TAB with no message. 

All chat sessions are saved in a session file, for which the `--session-file` flag can specify the file of your choice:

```bash
chatgpt-cli chat --session-file session.json
```

The `--model`, `--role`, `--temperature`, `--max-tokens`, and `--top-p` flags allow for individual configuration of the Model, Role, Temperature, Max Tokens, and Top P respectively. The default `--model` used will be updated overtime as new models are released.

A system prompt can be set by using the `--system-message` flag:

```bash
chatgpt-cli chat --system-message "You are a captivating storyteller who brings history to life by narrating the events, people, and cultures of the past."
```

### Replaying a Session

Replaying a chat session lets you revisit a previous chat in a more readable format than the raw JSON. Use the `replay-session` command:

```bash
chatgpt-cli replay-session --session-file session.json
```

### Refer to an image in a Chat

Initiate a chat with images uploaded to ChatGPT using the `vision` command:

```bash
chatgpt-cli vision --file image.jpg --file image2.jpg
```

This will upload the image to ChatGPT for use in the chat.
It also supports the session file, you can start a chat with `chat`, use `vision` to upload a file and then continue to ask more questions with `chat` by linking with a `--session-file` between commands.

You'll be prompted to input your message, which can span multiple lines. Send your message with TAB or CTRL+C.

### Generating Images

Generate an image with GPT-Image-1 or DALL-E using the `image` command:

```bash
chatgpt-cli image
```

You'll be prompted to input your description of an image, which can span multiple lines. Send your description with TAB or CTRL+C.

Exiting the chat is made possible by inputting CTRL+C or TAB with no description.

All images are saved with a prefix in the form `image-DATE-TIME-nn.png` where DATE-TIME is the timestamp when the session started, and nn for the image number from the session. You can override the ``--output-prefix`` or `-o` flags.

You can control how many variants of the requested images to generate with the `--number` or `-n` flag. The Number of Images must be between 1 and 10, inclusive.
 
You can control the size of the requested images with the `--size` or `-s` flag. The allowed sizes vary based on the model used.

### Generating Text to Speech

Generate an audio file, reading some text using the `speech` command:

```bash
chatgpt-cli speech
```

You'll be prompted to input the text you would like to read, which can span multiple lines. Send your description with TAB or CTRL+C.

Exiting the chat is made possible by inputting CTRL+C or TAB with no description.

Audio files are saved with a prefix in the form `tts-DATE-TIME-nn.mp3` where DATE-TIME is the timestamp when the session started, and nn for the image number from the session. You can override the ``--output-prefix`` or `-o` flags.

You can control the speed of the audio with the `--speed` or `-s` flag. The speed must be between 0.25 and 4.0, inclusive.
You can control the voice of the audio with the `--voice` or `-v` flag. The voice must be one of alloy, echo, fable, onyx, nova, or shimmer.
You can control the TTS model with the `--model` or `-m` flag. The model must be one of tts-1, tts-1-hd, or canary-tts.

### Transcribing Audio to Text

Transcribe text from an audio file using the `transcribe` command:

```bash
chatgpt-cli transcribe --file audio.mp3
```

Audio files are accepted in the following formats: flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.

You can indicate the language of the input audio with the `--language` or `-l` flag. Supplying the input language in ISO-639-1 format will improve accuracy and latency.
You can control the speech to text model with the `--model` or `-m` flag. Currently only whisper-1 is allowed.
Optional prompt to guide the model's style or continue a previous audio segment, can be set by using the `--system-message` flag. The prompt should match the audio language.

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
goreleaser --snapshot --skip=publish --clean
```

For a full deployment, you'll need to add a GitHub token to the `./.github_token` file, then run the following:

```bash
git tag <vNewVersion>
goreleaser --clean
brew upgrade

docs/generate-demos.sh
```
