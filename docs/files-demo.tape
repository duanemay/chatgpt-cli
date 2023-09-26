Output files-demo.gif
Set FontSize 32
Set Width 2200
Set Height 600

Sleep 2s

Type@50ms '{ echo "What is the largest file in this directory?"; ls -l } | chatgpt-cli chat'
Sleep 1s
Enter
# Wait for the response.
Sleep 4s
Type@0ms 'The largest file in this directory is "translation-demo.gif" with a size of 853 KB.'

# Admire the output for a bit.
Sleep 10s
