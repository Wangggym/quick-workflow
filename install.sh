#!/bin/bash

# Source the base.sh file
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "$SCRIPT_DIR/base.sh"

# Make all scripts executable
chmod +x "$SCRIPT_DIR"/*.sh

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
        echo -e "${w} Alias $alias_name already exists in $RC_FILE. Skipping..."
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

echo -e "${y} Installation complete. Please run 'source $RC_FILE' or restart your terminal to use the new aliases."