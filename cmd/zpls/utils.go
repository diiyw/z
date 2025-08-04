package main

import (
	"runtime"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func creatURI(filename string) protocol.DocumentUri {
	if runtime.GOOS == "windows" {
		// 在Windows上，文件路径使用反斜杠，需要转换为正斜杠
		filename = strings.ReplaceAll(filename, "\\", "/")
	}
	return "file://" + filename
}

func offsetToPosition(offset int, content string) protocol.Position {
	lines := strings.Split(content, "\n")
	line := 0
	character := 0

	for i, l := range lines {
		if offset < len(l) {
			line = i
			character = offset
			break
		}
		offset -= len(l) + 1 // +1 for the newline character
	}

	return protocol.Position{
		Line:      uint32(line),
		Character: uint32(character),
	}
}

func createRange(startLine, startChar, endLine, endChar uint32) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      startLine,
			Character: startChar,
		},
		End: protocol.Position{
			Line:      endLine,
			Character: endChar,
		},
	}
}

func fullFileRange(content string) protocol.Range {
	lines := strings.Split(content, "\n")
	return protocol.Range{
		Start: protocol.Position{},
		End: protocol.Position{
			Line:      uint32(len(lines) - 1),
			Character: uint32(len(lines[len(lines)-1])),
		},
	}
}

// getTextFromVersion 从文档URI和内容变更中提取完整的文档文本
func getTextFromVersion(uri string, changes []any) string {
	// 获取当前文档
	doc, ok := Document().Get(uri)

	// 如果文档不存在，则从变更中获取文本
	if !ok {
		if len(changes) > 0 {
			return changes[0].(protocol.TextDocumentContentChangeEvent).Text
		}
		return ""
	}

	// 获取当前文档文本
	text := doc.Text

	// 应用所有变更
	for _, change := range changes {
		if v, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			// 根据范围应用变更
			text = applyChange(text, v)
		} else {
			// 如果没有指定范围，则替换整个文档
			text = change.(protocol.TextDocumentContentChangeEventWhole).Text
		}
	}

	return text
}

// applyChange 根据变更范围应用文本变更
func applyChange(original string, change protocol.TextDocumentContentChangeEvent) string {
	// 将文本按行分割
	lines := strings.Split(original, "\n")

	// 获取变更范围
	startLine := int(change.Range.Start.Line)
	startChar := int(change.Range.Start.Character)
	endLine := int(change.Range.End.Line)
	endChar := int(change.Range.End.Character)

	// 提取范围前的文本
	before := ""
	for i := 0; i < startLine; i++ {
		before += lines[i] + "\n"
	}
	if startLine < len(lines) {
		before += lines[startLine][:startChar]
	}

	// 提取范围后的文本
	after := ""
	if endLine < len(lines) {
		after += lines[endLine][endChar:]
	}
	for i := endLine + 1; i < len(lines); i++ {
		after += "\n" + lines[i]
	}

	// 组合变更前的文本、变更文本和变更后的文本
	return before + change.Text + after
}
