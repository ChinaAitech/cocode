<template>
  <div id="app">
    <!-- 登录页面 -->
    <div v-if="!isLoggedIn" class="login-container">
      <el-card class="login-card">
        <template #header>
          <div class="card-header">
            <h2>协同编程平台</h2>
          </div>
        </template>
        <el-form :model="loginForm" @submit.prevent="handleLogin">
          <el-form-item label="用户名">
            <el-input v-model="loginForm.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleLogin" :loading="loginLoading" style="width: 100%">
              登录
            </el-button>
          </el-form-item>
        </el-form>
        <div class="login-tips">
          <p>默认账号：</p>
          <p>admin / admin123</p>
          <p>userA / passwordA</p>
          <p>userB / passwordB</p>
        </div>
      </el-card>
    </div>

    <!-- 主界面 -->
    <div v-else class="main-container">
      <!-- 顶部导航栏 -->
      <el-header class="header">
        <div class="header-left">
          <h2>C++ 协同编程平台</h2>
        </div>
        <div class="header-right">
          <el-tag type="success">{{ currentUser.displayName }}</el-tag>
          <el-button @click="handleLogout" size="small" style="margin-left: 10px">
            退出
          </el-button>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-container class="content-container">
        <!-- 代码编辑区 -->
        <el-main class="editor-section">
          <div class="editor-header">
            <span>代码编辑器</span>
            <div class="online-users">
              在线用户:
              <el-tag
                v-for="user in onlineUsers"
                :key="user"
                size="small"
                style="margin-left: 5px"
              >
                {{ user }}
              </el-tag>
            </div>
          </div>
          <CodeEditor
            ref="codeEditor"
            :code="code"
            @update:code="handleCodeUpdate"
          />

          <!-- 输入输出区 -->
          <div class="io-section">
            <div class="input-area">
              <div class="area-header">
                输入区 (stdin)
                <el-upload
                  action="/api/upload"
                  :headers="{ 'X-Session-ID': sessionId }"
                  :show-file-list="false"
                  :before-upload="beforeInputUpload"
                  :on-success="handleInputFileSuccess"
                  accept=".txt"
                >
                  <el-button size="small" :icon="Upload">上传文件</el-button>
                </el-upload>
              </div>
              <el-input
                v-model="inputData"
                type="textarea"
                :rows="6"
                placeholder="输入程序的标准输入..."
                @input="handleInputChange"
              />
            </div>
            <div class="output-area">
              <div class="area-header">
                输出区 (stdout)
                <el-button @click="clearOutput" size="small" text>清空</el-button>
              </div>
              <el-input
                v-model="outputData"
                type="textarea"
                :rows="6"
                readonly
                placeholder="程序输出将显示在这里..."
              />
            </div>
          </div>

          <!-- 标答和对比区 -->
          <div class="io-section" style="margin-top: 10px">
            <div class="input-area">
              <div class="area-header">
                标准答案
                <el-upload
                  action="/api/upload"
                  :headers="{ 'X-Session-ID': sessionId }"
                  :show-file-list="false"
                  :before-upload="beforeInputUpload"
                  :on-success="handleAnswerFileSuccess"
                  accept=".txt"
                >
                  <el-button size="small" :icon="Upload">上传文件</el-button>
                </el-upload>
              </div>
              <el-input
                v-model="answerData"
                type="textarea"
                :rows="6"
                placeholder="输入标准答案..."
                @input="handleAnswerChange"
              />
            </div>
            <div class="output-area">
              <div class="area-header">
                对比结果
                <el-button @click="compareOutput" size="small" type="primary">对比</el-button>
              </div>
              <div class="compare-result" v-if="compareResult">
                <div v-if="compareResult.isMatch" class="match-success">
                  ✅ 输出完全匹配！
                </div>
                <div v-else class="match-failed">
                  <div class="diff-summary">❌ 输出不匹配（{{ compareResult.differentLines.length }} 行不同）</div>
                  <div class="diff-details">
                    <div v-for="(diff, idx) in compareResult.differentLines" :key="idx" class="diff-line">
                      <div class="line-number">第 {{ diff.line }} 行：</div>
                      <div class="expected">期望: {{ diff.expected || '(空行)' }}</div>
                      <div class="actual">实际: {{ diff.actual || '(空行)' }}</div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="no-compare">
                点击"对比"按钮查看结果
              </div>
            </div>
          </div>

          <!-- 编译和运行按钮 -->
          <div class="action-buttons">
            <el-button
              type="primary"
              @click="compileAndRun"
              :loading="compiling"
              :icon="VideoPlay"
            >
              编译并运行
            </el-button>
            <el-button
              @click="downloadCode"
              :icon="Download"
            >
              下载代码
            </el-button>
          </div>

          <!-- 日志区 -->
          <div class="log-section">
            <div class="log-header">
              <span>编译日志</span>
              <el-button @click="clearLog" size="small" text>清空日志</el-button>
            </div>
            <div class="log-content" ref="logContent">
              <pre>{{ compileLog }}</pre>
            </div>
          </div>
        </el-main>

        <!-- 聊天侧边栏 -->
        <el-aside width="300px" class="chat-section">
          <ChatPanel
            ref="chatPanel"
            :messages="chatMessages"
            :current-user="currentUser.username"
            :session-id="sessionId"
            @send-message="sendChatMessage"
          />
        </el-aside>
      </el-container>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { VideoPlay, Download, Upload } from '@element-plus/icons-vue'
import CodeEditor from './components/CodeEditor.vue'
import ChatPanel from './components/ChatPanel.vue'

// 登录相关
const isLoggedIn = ref(false)
const loginForm = ref({ username: '', password: '' })
const loginLoading = ref(false)
const currentUser = ref({ username: '', displayName: '' })
const sessionId = ref('')

// WebSocket
let ws = null
const onlineUsers = ref([])

// 代码相关
const code = ref('')
const inputData = ref('')
const outputData = ref('')
const compileLog = ref('等待编译...\n')
const compiling = ref(false)
const answerData = ref('') // 标准答案
const compareResult = ref(null) // 对比结果

// 聊天相关
const chatMessages = ref([])

// 组件引用
const codeEditor = ref(null)
const chatPanel = ref(null)
const logContent = ref(null)

// 登录处理
const handleLogin = async () => {
  if (!loginForm.value.username || !loginForm.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loginLoading.value = true
  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginForm.value)
    })

    const result = await response.json()
    if (result.success) {
      sessionId.value = result.sessionId
      currentUser.value = {
        username: result.username,
        displayName: result.displayName
      }
      isLoggedIn.value = true
      ElMessage.success('登录成功')

      // 连接WebSocket
      connectWebSocket()
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    ElMessage.error('登录失败: ' + error.message)
  } finally {
    loginLoading.value = false
  }
}

// 登出处理
const handleLogout = async () => {
  try {
    await fetch('/api/logout', {
      method: 'POST',
      headers: { 'X-Session-ID': sessionId.value }
    })
  } catch (error) {
    console.error('登出失败:', error)
  }

  if (ws) {
    ws.close()
  }

  isLoggedIn.value = false
  sessionId.value = ''
  currentUser.value = { username: '', displayName: '' }
  ElMessage.success('已退出登录')
}

// WebSocket连接
const connectWebSocket = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/ws?session=${sessionId.value}`

  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    console.log('WebSocket已连接')
    ElMessage.success('已连接到协同服务器')
  }

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)
      handleWebSocketMessage(message)
    } catch (error) {
      console.error('消息解析失败:', error)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket错误:', error)
    ElMessage.error('连接错误')
  }

  ws.onclose = () => {
    console.log('WebSocket已断开')
    if (isLoggedIn.value) {
      ElMessage.warning('连接已断开')
      setTimeout(connectWebSocket, 3000) // 3秒后重连
    }
  }
}

// 处理WebSocket消息
const handleWebSocketMessage = (message) => {
  switch (message.type) {
    case 'init':
      // 初始化所有数据
      if (message.data) {
        code.value = message.data.code || ''
        inputData.value = message.data.inputData || ''
        outputData.value = message.data.outputData || ''
        compileLog.value = message.data.compileLog || '等待编译...\n'
        answerData.value = message.data.answer || ''
      }
      break

    case 'user_join':
      onlineUsers.value = message.data.users || []
      chatMessages.value.push({
        type: 'system',
        message: `${message.username} 加入了房间`,
        timestamp: message.timestamp
      })
      break

    case 'user_leave':
      onlineUsers.value = message.data.users || []
      chatMessages.value.push({
        type: 'system',
        message: `${message.username} 离开了房间`,
        timestamp: message.timestamp
      })
      break

    case 'edit':
      // 其他用户的编辑
      if (message.username !== currentUser.value.username && message.data) {
        code.value = message.data.code
      }
      break

    case 'input_change':
      // 输入数据变化
      if (message.username !== currentUser.value.username && message.data) {
        inputData.value = message.data.input
      }
      break

    case 'answer_change':
      // 标准答案变化
      if (message.username !== currentUser.value.username && message.data) {
        answerData.value = message.data.answer
      }
      break

    case 'cursor':
      // 光标位置更新
      if (message.username !== currentUser.value.username) {
        // TODO: 显示其他用户的光标
      }
      break

    case 'chat':
      // 聊天消息
      chatMessages.value.push({
        type: 'chat',
        username: message.username,
        message: message.data.message,
        timestamp: message.timestamp
      })
      break

    case 'compile_result':
      // 编译结果（已经包含所有数据）
      compiling.value = false
      const result = message.data

      // 更新输出和日志（所有用户同步）
      if (result.output !== undefined) {
        outputData.value = result.output
      }
      if (result.compileLog !== undefined) {
        compileLog.value = result.compileLog
      }

      // 只对编译发起者显示提示
      if (message.username === currentUser.value.username) {
        if (result.success) {
          ElMessage.success('编译运行成功')
        } else {
          ElMessage.error('编译失败')
        }
      } else {
        // 其他用户看到的通知
        ElMessage.info(`${result.compiledBy || message.username} 执行了编译`)
      }

      // 滚动到日志底部
      setTimeout(() => {
        if (logContent.value) {
          logContent.value.scrollTop = logContent.value.scrollHeight
        }
      }, 100)
      break
  }
}

// 代码更新
const handleCodeUpdate = (newCode) => {
  code.value = newCode

  // 发送编辑消息
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'edit',
      data: { code: newCode }
    }))
  }
}

// 编译并运行
const compileAndRun = () => {
  if (!code.value.trim()) {
    ElMessage.warning('请先编写代码')
    return
  }

  compiling.value = true
  compileLog.value += `\n[${new Date().toLocaleTimeString()}] 开始编译...\n`

  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'compile',
      data: {
        code: code.value,
        input: inputData.value
      }
    }))
  }
}

// 发送聊天消息
const sendChatMessage = (message) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    // 支持文本和文件消息
    const msgData = typeof message === 'string'
      ? { message, messageType: 'text' }
      : { ...message, messageType: message.type || 'text' }

    ws.send(JSON.stringify({
      type: 'chat',
      data: msgData
    }))
  }
}

// 下载代码
const downloadCode = () => {
  const url = `/api/download/code?session=${sessionId.value}`
  window.location.href = url
}

// 输入数据变化
const handleInputChange = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'input_change',
      data: { input: inputData.value }
    }))
  }
}

// 标准答案变化
const handleAnswerChange = () => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'answer_change',
      data: { answer: answerData.value }
    }))
  }
}

// 文件上传前检查
const beforeInputUpload = (file) => {
  const maxSize = 10 * 1024 * 1024 // 10MB
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

// 输入文件上传成功
const handleInputFileSuccess = (response) => {
  if (response.success) {
    // 读取文件内容
    fetch(response.fileUrl)
      .then(res => res.text())
      .then(content => {
        inputData.value = content
        handleInputChange() // 同步到其他用户
        ElMessage.success('输入文件加载成功')
      })
      .catch(() => {
        ElMessage.error('读取文件内容失败')
      })
  }
}

// 标答文件上传成功
const handleAnswerFileSuccess = (response) => {
  if (response.success) {
    // 读取文件内容
    fetch(response.fileUrl)
      .then(res => res.text())
      .then(content => {
        answerData.value = content
        handleAnswerChange() // 同步到其他用户
        ElMessage.success('标答文件加载成功')
      })
      .catch(() => {
        ElMessage.error('读取文件内容失败')
      })
  }
}

// 对比输出
const compareOutput = () => {
  const outputLines = outputData.value.split('\n')
  const answerLines = answerData.value.split('\n')

  const maxLines = Math.max(outputLines.length, answerLines.length)
  const differentLines = []

  for (let i = 0; i < maxLines; i++) {
    const outputLine = outputLines[i] || ''
    const answerLine = answerLines[i] || ''

    if (outputLine.trim() !== answerLine.trim()) {
      differentLines.push({
        line: i + 1,
        expected: answerLine,
        actual: outputLine
      })
    }
  }

  compareResult.value = {
    isMatch: differentLines.length === 0,
    differentLines: differentLines
  }

  if (compareResult.value.isMatch) {
    ElMessage.success('输出完全匹配！')
  } else {
    ElMessage.warning(`发现 ${differentLines.length} 行不匹配`)
  }
}

// 清空输出
const clearOutput = () => {
  outputData.value = ''
}

// 清空日志
const clearLog = () => {
  compileLog.value = '日志已清空\n'
}

// 生命周期
onMounted(() => {
  // 可以在这里添加初始化逻辑
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})
</script>

<style scoped>
#app {
  height: 100vh;
  overflow: hidden;
}

/* 登录页面 */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
}

.card-header {
  text-align: center;
}

.login-tips {
  text-align: center;
  color: #909399;
  font-size: 12px;
  margin-top: 20px;
}

.login-tips p {
  margin: 5px 0;
}

/* 主界面 */
.main-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #409eff;
  color: white;
  padding: 0 20px;
  height: 60px !important;
}

.header h2 {
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.content-container {
  flex: 1;
  overflow: hidden;
}

/* 代码编辑区 */
.editor-section {
  display: flex;
  flex-direction: column;
  padding: 10px;
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: bold;
}

.online-users {
  display: flex;
  align-items: center;
  font-size: 14px;
}

/* 输入输出区 */
.io-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 10px;
}

.input-area, .output-area {
  display: flex;
  flex-direction: column;
}

.area-header {
  font-weight: bold;
  margin-bottom: 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 操作按钮 */
.action-buttons {
  margin: 10px 0;
  text-align: center;
}

/* 日志区 */
.log-section {
  margin-top: 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 150px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  margin-bottom: 5px;
}

.log-content {
  flex: 1;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 聊天区 */
.chat-section {
  border-left: 1px solid #ddd;
  background: #fafafa;
}

/* 对比结果区 */
.compare-result {
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  max-height: 200px;
  overflow-y: auto;
}

.match-success {
  color: #67c23a;
  font-size: 16px;
  font-weight: bold;
  text-align: center;
  padding: 20px;
}

.match-failed {
  color: #f56c6c;
}

.diff-summary {
  font-weight: bold;
  margin-bottom: 10px;
  font-size: 14px;
}

.diff-details {
  max-height: 150px;
  overflow-y: auto;
}

.diff-line {
  background: white;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 8px;
  margin-bottom: 8px;
  font-size: 12px;
}

.line-number {
  font-weight: bold;
  color: #909399;
  margin-bottom: 4px;
}

.expected {
  color: #67c23a;
  font-family: 'Courier New', monospace;
  padding: 2px 4px;
  background: #f0f9ff;
  border-radius: 2px;
  margin-bottom: 2px;
}

.actual {
  color: #f56c6c;
  font-family: 'Courier New', monospace;
  padding: 2px 4px;
  background: #fef0f0;
  border-radius: 2px;
}

.no-compare {
  text-align: center;
  color: #909399;
  padding: 40px 20px;
  font-size: 14px;
}
</style>
