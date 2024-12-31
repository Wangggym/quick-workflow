#!/bin/bash

# 获取测试文件所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 获取库文件的路径
LIB_SCRIPT="$SCRIPT_DIR/../lib/handle-single-zip.sh"

# 检查库文件是否存在
if [ ! -f "$LIB_SCRIPT" ]; then
    echo "❌ Cannot find library script: $LIB_SCRIPT"
    exit 1
fi

# 导入被测试的函数
source "$LIB_SCRIPT"

# 设置测试环境
setup() {
    TEST_DIR=$(mktemp -d)
    mkdir -p "$TEST_DIR/output"
}

# 清理测试环境
cleanup() {
    rm -rf "$TEST_DIR"
}

# 测试：单个zip文件的情况
test_single_zip_file() {
    local test_name="test_single_zip_file"
    echo "Running $test_name..."
    
    # 创建测试数据
    touch "$TEST_DIR/output/test.zip"
    
    # 运行测试
    if handle_single_zip_file "$TEST_DIR/output" "test_logs"; then
        echo "✅ $test_name: PASSED"
    else
        echo "❌ $test_name: FAILED - Expected success for single zip file"
        return 1
    fi
}

# 测试：多个zip文件的情况
test_multiple_zip_files() {
    local test_name="test_multiple_zip_files"
    echo "Running $test_name..."
    
    # 创建测试数据
    touch "$TEST_DIR/output/test1.zip"
    touch "$TEST_DIR/output/test2.zip"
    
    # 运行测试
    if ! handle_single_zip_file "$TEST_DIR/output" "test_logs"; then
        echo "✅ $test_name: PASSED"
    else
        echo "❌ $test_name: FAILED - Expected failure for multiple zip files"
        return 1
    fi
}

# 测试：没有zip文件的情况
test_no_zip_files() {
    local test_name="test_no_zip_files"
    echo "Running $test_name..."
    
    # 运行测试
    if ! handle_single_zip_file "$TEST_DIR/output" "test_logs"; then
        echo "✅ $test_name: PASSED"
    else
        echo "❌ $test_name: FAILED - Expected failure for no zip files"
        return 1
    fi
}

# 运行所有测试
run_tests() {
    setup
    
    local failed=0
    test_single_zip_file || failed=1
    test_multiple_zip_files || failed=1
    test_no_zip_files || failed=1
    
    cleanup
    
    if [ $failed -eq 0 ]; then
        echo "🎉 All tests passed!"
        return 0
    else
        echo "💥 Some tests failed!"
        return 1
    fi
}

# 执行测试
run_tests