#!/usr/bin/env bash
set -e

cd docs || exit
rm chatgpt-cli*.json || true
export  PS1='\[\033[31;1m\]\$\[\033[m\] '

for file in *-demo.tape; do
  echo "Running ${file}"
  vhs "${file}"
  rm chatgpt-cli*.json || true
  echo
done