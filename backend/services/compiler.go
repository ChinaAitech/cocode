package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cocode/backend/config"
	"cocode/backend/models"
)

// CompileAndRun 编译并运行C++代码
func CompileAndRun(code string, input string) *models.CompileResult {
	result := &models.CompileResult{
		Success: false,
	}

	// 确保临时目录存在
	tempDir := config.AppConfig.Compiler.TempDir
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		result.Message = fmt.Sprintf("创建临时目录失败: %v", err)
		return result
	}

	// 创建临时源文件
	timestamp := time.Now().UnixNano()
	sourceFile := filepath.Join(tempDir, fmt.Sprintf("code_%d.cpp", timestamp))
	execFile := filepath.Join(tempDir, fmt.Sprintf("exec_%d", timestamp))
	inputFile := filepath.Join(tempDir, fmt.Sprintf("input_%d.txt", timestamp))
	outputFile := filepath.Join(tempDir, fmt.Sprintf("output_%d.txt", timestamp))

	// 清理临时文件
	defer func() {
		os.Remove(sourceFile)
		os.Remove(execFile)
		os.Remove(inputFile)
		os.Remove(outputFile)
	}()

	// 写入源代码
	if err := os.WriteFile(sourceFile, []byte(code), 0644); err != nil {
		result.Message = fmt.Sprintf("写入源文件失败: %v", err)
		return result
	}

	// 写入输入数据
	if err := os.WriteFile(inputFile, []byte(input), 0644); err != nil {
		result.Message = fmt.Sprintf("写入输入文件失败: %v", err)
		return result
	}

	// 编译代码
	compileTimeout := time.Duration(config.AppConfig.Compiler.CompileTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), compileTimeout)
	defer cancel()

	compileArgs := append(config.AppConfig.Compiler.CompileFlags, sourceFile, "-o", execFile)
	compileCmd := exec.CommandContext(ctx, config.AppConfig.Compiler.Compiler, compileArgs...)

	var compileOut bytes.Buffer
	var compileErr bytes.Buffer
	compileCmd.Stdout = &compileOut
	compileCmd.Stderr = &compileErr

	if err := compileCmd.Run(); err != nil {
		result.Message = fmt.Sprintf("编译失败:\n%s%s", compileOut.String(), compileErr.String())
		return result
	}

	result.Message = "编译成功!\n" + compileOut.String()

	// 运行程序
	runTimeout := time.Duration(config.AppConfig.Compiler.RunTimeout) * time.Second
	runCtx, runCancel := context.WithTimeout(context.Background(), runTimeout)
	defer runCancel()

	runCmd := exec.CommandContext(runCtx, execFile)

	// 设置输入
	inputData, _ := os.ReadFile(inputFile)
	runCmd.Stdin = bytes.NewReader(inputData)

	var runOut bytes.Buffer
	var runErr bytes.Buffer
	runCmd.Stdout = &runOut
	runCmd.Stderr = &runErr

	if err := runCmd.Run(); err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			result.Message += "\n运行超时!"
			result.Output = runOut.String()
			return result
		}
		result.Message += fmt.Sprintf("\n运行时错误:\n%s", runErr.String())
		result.Output = runOut.String()
		return result
	}

	// 运行成功
	result.Success = true
	result.Message += "\n运行成功!"
	result.Output = runOut.String()

	return result
}
