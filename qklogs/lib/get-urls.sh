#!/bin/bash

# 从JIRA issue获取附件URL
get_attachment_urls() {
    local issue_key="$1"
    
    # 检查JIRA认证信息
    if [ -z "$JIRA_API_TOKEN" ]; then
        echo "错误: 未设置JIRA_API_TOKEN环境变量" >&2
        return 1
    fi

    echo "正在获取 $issue_key 的附件信息..." >&2

    # 保存原始输出到临时文件
    local temp_file=$(mktemp)
    jira issue view "$issue_key" --plain > "$temp_file"

    # 调试：打印找到的文件
    echo "DEBUG: 查找的文件:" >&2
    grep -A 1 "^[[:space:]]*[0-9]\+\. log\." "$temp_file" >&2

    # 使用 awk 处理文件，将多行 URL 合并为单行
    local attachments=$(awk '
        BEGIN { 
            in_file = 0 
            filename = ""
            url = ""
            RS = "\n  [0-9]+\\. "  # 使用序号作为记录分隔符
            FS = "\n"
        }
        # 只处理以 log. 开头的文件
        $1 ~ /^log\./ {
            filename = $1
            gsub(/[[:space:]]*$/, "", filename)
            
            # 合并剩余的所有行作为URL
            url = ""
            for (i = 2; i <= NF; i++) {
                if ($i ~ /^[[:space:]]*http/) {
                    line = $i
                    gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
                    if (url == "") {
                        url = line
                    } else {
                        url = url line
                    }
                }
            }
            
            if (filename != "" && url != "") {
                print filename "§" url
            }
        }
    ' "$temp_file" | sort | uniq)

    # 清理临时文件
    rm "$temp_file"

    if [ -z "$attachments" ]; then
        echo "未找到附件" >&2
        return 1
    fi

    echo "$attachments"
    return 0
} 