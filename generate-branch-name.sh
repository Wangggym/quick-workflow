#!/bin/bash

generate_branch_name_from_input() {
    local input_text="$1"
    local api_key="$2"

    # OpenAI API 配置
    local API_URL="http://openai-proxy.brain.loocaa.com/v1/chat/completions"

    local json_data=$(cat <<EOF
{
    "model": "gpt-3.5-turbo",
    "stream": false,
    "messages": [
        {
            "role": "system",
            "content": "As a skilled linguist fluent in both English and Chinese, extract the key terms from the user's input (which may contain both Chinese and English) and generate a concise, descriptive branch name in English. Return only the branch name as a single string, and if other formats are present, convert them into a string."
        },
        {
            "role": "user",
            "content": "Based on the user's input '${input_text}', extract the key ideas and generate a concise branch name in English. The branch name should reflect the main concept of the suggestion without unnecessary detail."
        }
    ]
}
EOF
    )

    # 发送请求到 OpenAI API
    local response=$(curl -s --progress-bar --location --max-time 30 -w "%{http_code}" "$API_URL" \
        --header 'Content-Type: application/json' \
        --header "Authorization: Bearer $api_key" \
        --data "$json_data")

    # 提取状态码和响应内容
    local http_code="${response: -3}"
    local response_body="${response:0:${#response}-3}"

    # 检查状态码
    if [ "$http_code" -ne 200 ]; then
        echo "❌ Fetch branch name from AI failed, code: $http_code"
        return 1
    fi

    # 从响应中提取分支名
    local branch_name=$(echo "$response_body" | jq -r '.choices[0].message.content')

    # 清理分支名
    branch_name=$(echo "$branch_name" | \
        sed 's/^[ \t]*//;s/[ \t]*$//' | \
        tr '[:upper:]' '[:lower:]' | \
        tr ' ' '-' | \
        sed 's/[^a-z0-9-]//g' | \
        sed 's/-\+/-/g' | \
        sed 's/^-\|-$//g')

    # 检查 branch_name 是否为空
    if [ -z "$branch_name" ]; then
        echo "❌ Fetch branch name from AI failed, branch_name is empty"
        return 1
    fi

    echo "$branch_name"
    return 0
}
