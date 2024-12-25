#!/usr/bin/env python3

import requests
import json
import sys
import re

def make_stream_request(input_text, api_key):

    if not input_text:
      print("Input text is required")
      return ""
    if not api_key:
      print("API key is required")
      return ""

    url = 'http://openai-proxy.brain.loocaa.com/v1/chat/completions'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {api_key}'
    }

    payload = {
        "model": "gpt-3.5-turbo",
        "stream": True,
        "messages": [
            {
                "role": "system",
                "content": "As a skilled linguist fluent in both English and Chinese, extract the key terms from the user's input (which may contain both Chinese and English) and generate a concise, descriptive branch name in English. Return only the branch name as a single string, and if other formats are present, convert them into a string."
            },
            {
                "role": "user",
                "content": f"Based on the user's input '{input_text}', extract the key ideas and generate a concise branch name in English. The branch name should reflect the main concept of the suggestion without unnecessary detail."
            }
        ]
    }

    response = requests.post(url, headers=headers, json=payload, stream=True)

    # 检查响应状态码
    if response.status_code != 200:
      return ""

    # 用于存储完整的响应
    full_response = ""

    for line in response.iter_lines():
        if line:
            # 先解码字节字符串
            line = line.decode('utf-8')
            if line.startswith('data: '):
                line = line[6:]

            # 跳过心跳消息
            if line == '[DONE]':
                break

            try:
                json_response = json.loads(line)
                if 'choices' in json_response:
                    content = json_response['choices'][0].get('delta', {}).get('content', '')
                    if content:
                        full_response += content
            except json.JSONDecodeError:
                continue

    return full_response

def clean_branch_name(branch_name):
    # 去除前后空格
    branch_name = branch_name.strip()

    # 转换为小写
    branch_name = branch_name.lower()

    # 替换空格为短横线
    branch_name = branch_name.replace(' ', '-')

    # 移除非字母数字和短横线的字符
    branch_name = re.sub(r'[^a-z0-9-]', '', branch_name)

    # 替换多个短横线为一个短横线
    branch_name = re.sub(r'-+', '-', branch_name)

    # 去除开头和结尾的短横线
    branch_name = branch_name.strip('-')

    return branch_name

def save_branch_name(result, identifier):
    file_path = f"/tmp/branch_name_{identifier}.txt"
    with open(file_path, 'w') as file:
        file.write(result)
    print(f"Branch name saved to {file_path}")

if __name__ == "__main__":
    if len(sys.argv) > 2:
        input_text = sys.argv[1]
        api_key = sys.argv[2]
        identifier = sys.argv[3]
    else:
      print("Usage: python stream_request.py <input_text> <api_key> <identifier>")
      sys.exit(1)

    result = make_stream_request(input_text, api_key)
    if result:
        branch_name = clean_branch_name(result)
        save_branch_name(branch_name, identifier)
        sys.exit(0)
    else:
        print(f"Failed to get branch name {result}")
        sys.exit(1)
