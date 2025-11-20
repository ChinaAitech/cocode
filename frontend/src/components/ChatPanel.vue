<template>
  <div class="chat-panel">
    <div class="chat-header">
      <span>聊天室</span>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message-wrapper', msg.type, { 'is-self': msg.username === currentUser }]"
      >
        <div v-if="msg.type === 'system'" class="system-message">
          <el-icon><InfoFilled /></el-icon>
          {{ msg.message }}
        </div>
        <div v-else class="chat-message" :style="{ borderLeftColor: getUserColor(msg.username) }">
          <div class="message-header">
            <span class="username" :style="{ color: getUserColor(msg.username) }">
              {{ msg.username }}
            </span>
            <span class="timestamp">{{ formatTime(msg.timestamp) }}</span>
          </div>
          <div class="message-content">{{ msg.message }}</div>
        </div>
      </div>
    </div>

    <div class="chat-input">
      <div class="input-toolbar">
        <el-button size="small" @click="showEmojiPicker = !showEmojiPicker" :icon="ChatDotRound">
          表情
        </el-button>
        <el-upload
          action="/api/upload"
          :headers="{ 'X-Session-ID': sessionId }"
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          accept="image/*"
        >
          <el-button size="small" :icon="Picture">图片</el-button>
        </el-upload>
        <el-upload
          action="/api/upload"
          :headers="{ 'X-Session-ID': sessionId }"
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <el-button size="small" :icon="Paperclip">文件</el-button>
        </el-upload>
      </div>

      <!-- 表情选择器 -->
      <div v-if="showEmojiPicker" class="emoji-picker">
        <span
          v-for="emoji in emojis"
          :key="emoji"
          class="emoji-item"
          @click="insertEmoji(emoji)"
        >
          {{ emoji }}
        </span>
      </div>

      <el-input
        v-model="inputMessage"
        type="textarea"
        :rows="2"
        placeholder="输入消息... (Enter发送，Shift+Enter换行)"
        @keydown.enter="handleEnterKey"
      />
      <el-button
        type="primary"
        @click="sendMessage"
        :icon="Promotion"
        style="margin-top: 5px; width: 100%"
      >
        发送
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { InfoFilled, Promotion, ChatDotRound, Picture, Paperclip } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  currentUser: {
    type: String,
    default: ''
  },
  sessionId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['send-message'])

const inputMessage = ref('')
const messagesContainer = ref(null)
const showEmojiPicker = ref(false)

// 常用表情列表
const emojis = [
  '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂',
  '😊', '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰',
  '😘', '😗', '😙', '😚', '😋', '😛', '😝', '😜',
  '🤪', '🤨', '🧐', '🤓', '😎', '🤩', '🥳', '😏',
  '👍', '👎', '👌', '✌️', '🤞', '🤝', '👏', '🙌',
  '💪', '🎉', '🎊', '🎈', '🎁', '🏆', '🥇', '⭐',
  '❤️', '💕', '💖', '💗', '💙', '💚', '💛', '🧡'
]

// 用户颜色映射表
const userColors = ref({})
const colorPalette = [
  '#409eff', // 蓝色
  '#67c23a', // 绿色
  '#e6a23c', // 橙色
  '#f56c6c', // 红色
  '#909399', // 灰色
  '#c71585', // 紫红色
  '#ff6347', // 番茄色
  '#4682b4', // 钢蓝色
  '#32cd32', // 酸橙色
  '#ff8c00', // 深橙色
  '#9370db', // 中紫色
  '#20b2aa', // 浅海绿色
  '#ff1493', // 深粉色
  '#00ced1', // 深青色
  '#ff69b4', // 热粉色
  '#8a2be2', // 蓝紫色
  '#00bfff', // 深天蓝色
  '#adff2f', // 黄绿色
  '#ff4500', // 橙红色
  '#da70d6'  // 兰花紫
]

// 获取用户颜色
const getUserColor = (username) => {
  if (!username) return colorPalette[0]

  if (!userColors.value[username]) {
    // 为新用户分配颜色
    const existingColors = Object.keys(userColors.value).length
    userColors.value[username] = colorPalette[existingColors % colorPalette.length]
  }

  return userColors.value[username]
}

// 插入表情
const insertEmoji = (emoji) => {
  inputMessage.value += emoji
  showEmojiPicker.value = false
}

// 处理Enter键
const handleEnterKey = (e) => {
  if (!e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

// 文件上传前检查
const beforeUpload = (file) => {
  const maxSize = 50 * 1024 * 1024 // 50MB
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过 50MB')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.success) {
    emit('send-message', {
      type: 'file',
      fileUrl: response.fileUrl,
      fileName: response.fileName,
      fileSize: response.fileSize
    })
    ElMessage.success('文件上传成功')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

// 发送消息
const sendMessage = () => {
  if (!inputMessage.value.trim()) {
    return
  }

  emit('send-message', inputMessage.value)
  inputMessage.value = ''
  showEmojiPicker.value = false
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp * 1000)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 自动滚动到底部
watch(() => props.messages.length, async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
})
</script>

<style scoped>
.chat-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.chat-header {
  font-weight: bold;
  font-size: 16px;
  padding: 10px 0;
  border-bottom: 1px solid #ddd;
  margin-bottom: 10px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 10px;
}

.message-wrapper {
  margin-bottom: 15px;
  display: flex;
}

/* 其他人的消息靠左 */
.message-wrapper:not(.is-self) {
  justify-content: flex-start;
}

/* 自己的消息靠右 */
.message-wrapper.is-self {
  justify-content: flex-end;
}

.system-message {
  text-align: center;
  color: #909399;
  font-size: 12px;
  padding: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  margin: 0 auto;
}

.chat-message {
  padding: 10px 12px;
  background: #f0f9ff;
  border-radius: 8px;
  border-left: 3px solid #409eff;
  max-width: 80%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 自己的消息样式 */
.message-wrapper.is-self .chat-message {
  background: #e8f5e9;
  border-left: none;
  border-right: 3px solid #67c23a;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  align-items: center;
}

.username {
  font-weight: bold;
  font-size: 13px;
}

.timestamp {
  color: #909399;
  font-size: 11px;
  margin-left: 8px;
}

.message-content {
  color: #303133;
  font-size: 14px;
  word-wrap: break-word;
  white-space: pre-wrap;
  line-height: 1.5;
}

.chat-input {
  margin-top: auto;
}

.input-toolbar {
  display: flex;
  gap: 5px;
  margin-bottom: 8px;
}

.emoji-picker {
  background: white;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 8px;
  max-height: 150px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 5px;
}

.emoji-item {
  font-size: 24px;
  cursor: pointer;
  text-align: center;
  padding: 5px;
  border-radius: 4px;
  transition: all 0.2s;
}

.emoji-item:hover {
  background: #f5f7fa;
  transform: scale(1.2);
}
</style>
