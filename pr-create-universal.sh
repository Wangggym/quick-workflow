#!/bin/bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
source "$script_dir/env_load.sh"
source "$script_dir/repo_type.sh"

# 通用 PR 创建脚本
# 参数：--title "标题" --body "描述" -H 源分支 --base 目标分支

# 自动加载 .env 文件（如果存在）
# (已由 env_load.sh 实现)

# 解析参数
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --title)
      commit_title="$2"
      shift; shift
      ;;
    --body)
      pr_body="$2"
      shift; shift
      ;;
    -H)
      branch_name="$2"
      shift; shift
      ;;
    --base)
      base_branch="$2"
      shift; shift
      ;;
    *)
      shift
      ;;
  esac
done

# 如果未指定 base_branch，默认 develop
if [[ -z "$base_branch" ]]; then
  base_branch="develop"
fi

# 获取远程仓库地址
git_url=$(git remote get-url origin 2>/dev/null)
if [[ -z "$git_url" ]]; then
  echo "[Error] 未检测到 git 仓库或远程地址。"
  exit 1
fi

# 用 repo_type 变量判断仓库类型
if [[ "$repo_type" == "github" ]]; then
  # GitHub 仓库，调用 gh pr create
  pr_url=$(gh pr create --title "${commit_title}" --body "${pr_body}" -H "$branch_name")
  echo "$pr_url"
  exit 0
elif [[ "$repo_type" == "codeup" ]]; then
  # Codeup 仓库，调用 Codeup API
  # 需要环境变量: CODEUP_PROJECT_ID, CODEUP_CSRF_TOKEN, CODEUP_COOKIE
  if [[ -z "$CODEUP_PROJECT_ID" || -z "$CODEUP_CSRF_TOKEN" || -z "$CODEUP_COOKIE" ]]; then
    echo "[Error] 缺少 Codeup 所需的环境变量 (CODEUP_PROJECT_ID, CODEUP_CSRF_TOKEN, CODEUP_COOKIE)"
    exit 2
  fi
  data_raw=$(cat <<EOF
{
  "source_project_id": ${CODEUP_PROJECT_ID},
  "target_project_id": ${CODEUP_PROJECT_ID},
  "source_branch": "${branch_name}",
  "target_branch": "${base_branch}",
  "title": "${commit_title}",
  "description": "${pr_body}",
  "tb_user_ids": [],
  "reviewer_user_ids": [],
  "create_from": "WEB"
}
EOF
  )

  response=$(curl -s --location --request POST "https://codeup.aliyun.com/api/v4/projects/${CODEUP_PROJECT_ID}/code_reviews?_csrf=${CODEUP_CSRF_TOKEN}&_input_charset=utf-8" \
    --header 'X-Requested-With: XMLHttpRequest' \
    --header "Cookie: ${CODEUP_COOKIE}" \
    --header 'Content-Type: application/json' \
    --data-raw "$data_raw")
  pr_url=$(echo "$response" | jq -r '.detail_url // empty')
  if [[ -z "$pr_url" ]]; then
    echo "$response"
    echo "[Error] 未能获取 PR URL，请检查 API 响应。"
    exit 3
  fi
  echo "$pr_url"
  exit 0
else
  echo "[Error] 未知的仓库类型：$git_url"
  exit 4
fi 