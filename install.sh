#!/bin/bash

# Source the base.sh file
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "$SCRIPT_DIR/base.sh"

# Make all scripts executable
chmod +x "$SCRIPT_DIR"/*.sh

# Define the alias
ALIAS_NAME="qkupdate"
SCRIPT_PATH="$SCRIPT_DIR/get-pr-titile.sh"
ALIAS_COMMAND="alias $ALIAS_NAME=\"$SCRIPT_PATH\""

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
        echo -e "${n} Unsupported shell: $USER_SHELL. Please add the alias manually to your shell's rc file."
        exit 1
        ;;
esac

# Check if the alias already exists
if grep -q "$ALIAS_NAME" "$RC_FILE"; then
    echo -e "${y} Alias $ALIAS_NAME already exists in $RC_FILE. Updating..."
    sed -i '' "/alias $ALIAS_NAME=/c\\
$ALIAS_COMMAND" "$RC_FILE"
else
    echo -e "${y} Adding alias $ALIAS_NAME to $RC_FILE"
    echo "$ALIAS_COMMAND" >> "$RC_FILE"
fi

echo -e "${y} Installation complete. Please run 'source $RC_FILE' or restart your terminal to use the new alias."