#!/bin/bash

git_url=$(git remote get-url origin 2>/dev/null)
if [[ "$git_url" == *"github.com"* ]]; then
  repo_type="github"
elif [[ "$git_url" == *"codeup.aliyun.com"* ]]; then
  repo_type="codeup"
else
  repo_type="unknown"
fi 