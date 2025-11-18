<template>
  <div ref="editorContainer" class="code-editor"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  code: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:code'])

const editorContainer = ref(null)
let editor = null
let isUpdatingFromProp = false

onMounted(() => {
  // 创建编辑器
  editor = monaco.editor.create(editorContainer.value, {
    value: props.code,
    language: 'cpp',
    theme: 'vs-dark',
    automaticLayout: true,
    fontSize: 14,
    minimap: {
      enabled: true
    },
    scrollBeyondLastLine: false,
    wordWrap: 'on'
  })

  // 监听内容变化
  editor.onDidChangeModelContent(() => {
    if (!isUpdatingFromProp) {
      const value = editor.getValue()
      emit('update:code', value)
    }
  })

  // 监听光标位置变化
  editor.onDidChangeCursorPosition((e) => {
    // TODO: 发送光标位置给其他用户
  })
})

// 监听prop变化
watch(() => props.code, (newCode) => {
  if (editor && editor.getValue() !== newCode) {
    isUpdatingFromProp = true
    const position = editor.getPosition()
    editor.setValue(newCode)
    if (position) {
      editor.setPosition(position)
    }
    isUpdatingFromProp = false
  }
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
})
</script>

<style scoped>
.code-editor {
  width: 100%;
  height: 400px;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
}
</style>
