<template>
  <div class="chat-panel">
    <div class="chat-header">
      <span>聊天室</span>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message', msg.type]"
      >
        <div v-if="msg.type === 'system'" class="system-message">
          <el-icon><InfoFilled /></el-icon>
          {{ msg.message }}
        </div>
        <div v-else class="chat-message">
          <div class="message-header">
            <span class="username">{{ msg.username }}</span>
            <span class="timestamp">{{ formatTime(msg.timestamp) }}</span>
          </div>
          <div class="message-content">{{ msg.message }}</div>
        </div>
      </div>
    </div>

    <div class="chat-input">
      <el-input
        v-model="inputMessage"
        placeholder="输入消息..."
        @keyup.enter="sendMessage"
      >
        <template #append>
          <el-button @click="sendMessage" :icon="Promotion">发送</el-button>
        </template>
      </el-input>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { InfoFilled, Promotion } from '@element-plus/icons-vue'

const props = defineProps({
  messages: {
    type: Array,
    default: () => []
  },
  currentUser: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['send-message'])

const inputMessage = ref('')
const messagesContainer = ref(null)

// 发送消息
const sendMessage = () => {
  if (!inputMessage.value.trim()) {
    return
  }

  emit('send-message', inputMessage.value)
  inputMessage.value = ''
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

.message {
  margin-bottom: 15px;
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
}

.chat-message {
  padding: 8px;
  background: #f0f9ff;
  border-radius: 8px;
  border-left: 3px solid #409eff;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.username {
  font-weight: bold;
  color: #409eff;
  font-size: 14px;
}

.timestamp {
  color: #909399;
  font-size: 12px;
}

.message-content {
  color: #303133;
  font-size: 14px;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.chat-input {
  margin-top: auto;
}
</style>
