#!/usr/bin/env bash
set -e

filter=${1:-*}

cd docs || exit
echo "Using filter=${filter}"
ls -l "${filter}-demo.tape";

rm chatgpt-cli*.json dall-e-*.png || true
export  PS1='\[\033[31;1m\]\$\[\033[m\] '

for file in ${filter}-demo.tape; do
  echo "Running ${file}"
  vhs "${file}"
  rm chatgpt-cli*.json dall-e-*.png || true
  echo
done