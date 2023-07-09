#!/usr/bin/env bash

cd docs || exit
rm chatgpt-cli*.json
export  PS1='\[\033[31;1m\]\$\[\033[m\] '

for file in *-demo.tape; do
  echo "Running $file"
  vhs $file
  rm chatgpt-cli*.json
  echo
done