#!/bin/bash

# Source the base.sh file
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "$SCRIPT_DIR/base.sh"

# Make all scripts executable
chmod +x "$SCRIPT_DIR"/*.sh

# Add environment variables if they don't exist
ENV_VARS=(
    "EMAIL=your.email@example.com"
    "JIRA_API_TOKEN=your_jira_api_token"
    "JIRA_SERVICE_ADDRESS=your_jira_service_address"
    "GH_BRANCH_PREFIX=your_branch_prefix"
)

# Function to add environment variable if it doesn't exist
add_env_var_if_not_exists() {
    local env_var=$1
    local var_name="${env_var%%=*}"
    if ! grep -q "export $var_name=" "$RC_FILE"; then
        echo -e "${y} Adding environment variable $var_name to $RC_FILE"
        echo "export $env_var" >> "$RC_FILE"
    else
        echo -e "${y} Environment variable $var_name already exists in $RC_FILE. Skipping..."
    fi
}

# Define the aliases
ALIAS_NAME="qkupdate"
SCRIPT_PATH="$SCRIPT_DIR/get-pr-titile.sh"
ALIAS_COMMAND="alias $ALIAS_NAME=\"$SCRIPT_PATH\""

PROXY_ALIAS_NAME="proxy"
PROXY_SCRIPT_PATH="$SCRIPT_DIR/check-set-proxy.sh"
PROXY_ALIAS_COMMAND="alias $PROXY_ALIAS_NAME=\"$PROXY_SCRIPT_PATH\""

# Define the new alias
QKLOGS_ALIAS_NAME="qklogs"
QKLOGS_SCRIPT_PATH="$SCRIPT_DIR/qklogs/qklogs.sh"
QKLOGS_ALIAS_COMMAND="alias $QKLOGS_ALIAS_NAME=\"$QKLOGS_SCRIPT_PATH\""

# Define the qkfind alias
QKFIND_ALIAS_NAME="qkfind"
QKFIND_SCRIPT_PATH="$SCRIPT_DIR/qkfind.sh"
QKFIND_ALIAS_COMMAND="alias $QKFIND_ALIAS_NAME=\"$QKFIND_SCRIPT_PATH\""

# Define the qksearch alias
QKSEARCH_ALIAS_NAME="qksearch"
QKSEARCH_SCRIPT_PATH="$SCRIPT_DIR/qksearch.sh"
QKSEARCH_ALIAS_COMMAND="alias $QKSEARCH_ALIAS_NAME=\"$QKSEARCH_SCRIPT_PATH\""

# Define the qk alias
QK_ALIAS_NAME="qk"
QK_SCRIPT_PATH="$SCRIPT_DIR/qk.sh"
QK_ALIAS_COMMAND="alias $QK_ALIAS_NAME=\"$QK_SCRIPT_PATH\""

# Define the pr-create alias
PR_CREATE_ALIAS_NAME="pr-create"
PR_CREATE_SCRIPT_PATH="$SCRIPT_DIR/pr-create.sh"
PR_CREATE_ALIAS_COMMAND="alias $PR_CREATE_ALIAS_NAME=\"$PR_CREATE_SCRIPT_PATH\""

# Define the pr-merge alias
PR_MERGE_ALIAS_NAME="pr-merge"
PR_MERGE_SCRIPT_PATH="$SCRIPT_DIR/pr-merge.sh"
PR_MERGE_ALIAS_COMMAND="alias $PR_MERGE_ALIAS_NAME=\"$PR_MERGE_SCRIPT_PATH\""

# Determine the user's default shell
USER_SHELL=$(basename "$SHELL")

# Determine the correct rc file
case "$USER_SHELL" in
    zsh)
        RC_FILE="$HOME/.zshrc"
        echo -e "${y} Detected Zsh as the default shell"
        ;;
    bash)
        RC_FILE="$HOME/.bashrc"
        echo -e "${y} Detected Bash as the default shell"
        ;;
    *)
        echo -e "${n} Unsupported shell: $USER_SHELL. Please add the aliases manually to your shell's rc file."
        exit 1
        ;;
esac

# Function to add an alias only if it doesn't exist
add_alias_if_not_exists() {
    local alias_name=$1
    local alias_command=$2
    if ! grep -q "alias $alias_name=" "$RC_FILE"; then
        echo -e "${y} Adding alias $alias_name to $RC_FILE"
        echo "$alias_command" >> "$RC_FILE"
    else
        echo -e "${y} Alias $alias_name already exists in $RC_FILE. Skipping..."
    fi
}

# Add the qkupdate alias if it doesn't exist
add_alias_if_not_exists "$ALIAS_NAME" "$ALIAS_COMMAND"

# Add the proxy alias if it doesn't exist
add_alias_if_not_exists "$PROXY_ALIAS_NAME" "$PROXY_ALIAS_COMMAND"

# Add the qklogs alias if it doesn't exist
add_alias_if_not_exists "$QKLOGS_ALIAS_NAME" "$QKLOGS_ALIAS_COMMAND"

# Add the qkfind alias if it doesn't exist
add_alias_if_not_exists "$QKFIND_ALIAS_NAME" "$QKFIND_ALIAS_COMMAND"

# Add the qksearch alias if it doesn't exist
add_alias_if_not_exists "$QKSEARCH_ALIAS_NAME" "$QKSEARCH_ALIAS_COMMAND"

# Add the qk alias if it doesn't exist
add_alias_if_not_exists "$QK_ALIAS_NAME" "$QK_ALIAS_COMMAND"

# Add the pr-create alias if it doesn't exist
add_alias_if_not_exists "$PR_CREATE_ALIAS_NAME" "$PR_CREATE_ALIAS_COMMAND"

# Add the pr-merge alias if it doesn't exist
add_alias_if_not_exists "$PR_MERGE_ALIAS_NAME" "$PR_MERGE_ALIAS_COMMAND"

# 检查依赖
check_dependencies() {
    if ! command -v python3 &> /dev/null; then
        echo "错误：未安装 Python3"
    fi

    python3 -c "import requests" 2>/dev/null || {
        echo "正在安装 requests 包..."
        pip3 install requests
    }
}

check_dependencies

# Add environment variables
for env_var in "${ENV_VARS[@]}"; do
    add_env_var_if_not_exists "$env_var"
done

echo -e "${y} Installation complete. Please run 'source $RC_FILE' or restart your terminal to use the new aliases and environment variables."
echo -e "${y} Remember to update the environment variables in $RC_FILE with your actual values."