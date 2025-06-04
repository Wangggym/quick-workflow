#!/bin/bash

script_dir="$(dirname "$0")"
source "$script_dir/env_load.sh"
source "$script_dir/base.sh"
source "$script_dir/repo_type.sh"

pr_title=""

if [[ "$repo_type" == "github" ]]; then
  pr_title=$(gh pr view --json title --jq .title 2>/dev/null)
elif [[ "$repo_type" == "codeup" ]]; then
  current_branch=$(git rev-parse --abbrev-ref HEAD)
  if [[ -z "$CODEUP_PROJECT_ID" || -z "$CODEUP_COOKIE" ]]; then
    echo "[Error] 缺少 Codeup 所需的环境变量 (CODEUP_PROJECT_ID, CODEUP_COOKIE)"
  else
    codeup_api_url="https://codeup.aliyun.com/api/v4/projects/code_reviews/advanced_search_cr?_input_charset=utf-8&page=1&search=&order_by=updated_at&state=opened&project_ids=${CODEUP_PROJECT_ID}&sub_state_list=wip%2Cunder_review&per_page=10"
    pr_json=$(curl -s --location --request GET "$codeup_api_url" \
      --header 'X-Requested-With: XMLHttpRequest' \
      --header "Cookie: $CODEUP_COOKIE" \
      --header 'Content-Type: application/x-www-form-urlencoded')
    pr_title=$(echo "$pr_json" | jq -r ".[] | select(.source_branch==\"$current_branch\") | .title" | head -n 1)
  fi
else
  echo "[Error] 未知的仓库类型：$git_url"
fi

if [ -n "$pr_title" ]; then 
    echo -e "${y} get pr title: $pr_title"
    git add --all && git commit -m "$pr_title" && git push
else
    echo -e "${r} Failed to retrieve PR title."
    git add --all && git commit -m "update" && git push
fi

