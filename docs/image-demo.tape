Output image-demo.gif
Set FontSize 32
Set Width 2200
Set Height 600

Sleep 2s

Type@100ms "chatgpt-cli image -o maltipoo"
Sleep 2s
Enter
# Admire the output for a bit.
Sleep 8s

Type@250ms "white maltipoo dressed in a shark costume, on a dock by a lake"
Sleep 2s
Enter
Tab

Sleep 8s
Ctrl+C

Type@250ms "ls maltipoo*"
Sleep 2s
Enter

# now let's see what got generated
Type@100ms "chatgpt-cli vision --file maltipoo-01.png"
Sleep 2s
Enter
Sleep 5s

Type@300ms "Describe what is in this image."
Sleep 2s
Enter
Tab

# Admire the output for a bit.
Sleep 20s