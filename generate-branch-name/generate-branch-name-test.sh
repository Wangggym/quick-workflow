SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source $SCRIPT_DIR/generate-branch-name.sh

# 主函数
__main__() {
    check_dependencies
    commit_title="Fetch branch name from AI cost:"
    echo "SCRIPT_DIR: $SCRIPT_DIR"
    generate_id=$(echo -n $commit_title | md5sum | awk '{print $1}')
    echo "generate_id: $generate_id"
    generate_branch_name "$commit_title" "$BRAIN_AI_KEY" "$generate_id"
    echo "branch name: $(cat /tmp/branch_name_$generate_id.txt)"
    rm -f /tmp/branch_name_$generate_id.txt
}

# 执行主函数，传递所有命令行参数
__main__