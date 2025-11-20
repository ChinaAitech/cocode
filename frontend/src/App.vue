<template>
  <div id="app">
    <!-- 登录页面 -->
    <!-- 登录页面 -->
    <div v-if="!isLoggedIn" class="login-container">
      <div class="login-background"></div>
      <el-card class="login-card glass-effect">
        <template #header>
          <div class="card-header">
            <h2 class="gradient-text">协同编程平台</h2>
            <p class="subtitle">Next Gen Collaborative Coding</p>
          </div>
        </template>
        <el-form :model="loginForm" @submit.prevent="handleLogin" class="login-form">
          <el-form-item>
            <el-input 
              v-model="loginForm.username" 
              placeholder="用户名" 
              prefix-icon="User"
              class="glass-input"
            />
          </el-form-item>
          <el-form-item>
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              prefix-icon="Lock"
              @keyup.enter="handleLogin"
              class="glass-input"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleLogin" :loading="loginLoading" class="login-btn">
              登 录
            </el-button>
          </el-form-item>
        </el-form>
        <div class="login-tips">
          <p>默认账号：admin / admin123</p>
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
          <el-button
            v-if="currentUser.username === 'admin'"
            @click="showUserManagement"
            size="small"
            style="margin-left: 10px"
          >
            用户管理
          </el-button>
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
            <div class="header-left-controls">
              <span>代码编辑器</span>
              <el-switch
                v-model="showModules"
                active-text="显示附加模块"
                inactive-text="隐藏附加模块"
                style="margin-left: 20px"
              />
              <div class="font-size-controls" style="margin-left: 20px">
                <el-button-group>
                  <el-button size="small" @click="decreaseFontSize" :icon="Minus" />
                  <el-button size="small" disabled>{{ fontSize }}px</el-button>
                  <el-button size="small" @click="increaseFontSize" :icon="Plus" />
                </el-button-group>
              </div>
            </div>
            <div class="header-actions">
              <el-button @click="downloadCode" size="small" :icon="Download">
                下载代码
              </el-button>
              <div class="online-users">
                在线用户:
                <el-tag
                  v-for="user in onlineUsers"
                  :key="user.username"
                  size="small"
                  style="margin-left: 5px"
                >
                  {{ user.displayName }}
                </el-tag>
              </div>
            </div>
          </div>
          <div class="code-editor-container">
            <CodeEditor
              ref="codeEditor"
              :code="code"
              :font-size="fontSize"
              @update:code="handleCodeUpdate"
            />
          </div>

          <!-- 输入输出标答三列布局 -->
          <div v-if="showModules" class="modules-container">
            <div class="io-section-three">
              <div class="io-column">
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
                    <el-button size="small" :icon="Upload">文件</el-button>
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
              <div class="io-column">
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
              <div class="io-column">
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
                    <el-button size="small" :icon="Upload">文件</el-button>
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
            </div>

            <!-- 编译和对比按钮 -->
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
                type="success"
                @click="compareOutput"
                :icon="Check"
              >
                对比输出
              </el-button>
            </div>

            <!-- 日志和对比结果区 -->
            <div class="log-compare-section">
              <div class="log-half">
                <div class="log-header">
                  <span>编译日志</span>
                  <el-button @click="clearLog" size="small" text>清空</el-button>
                </div>
                <div class="log-content" ref="logContent">
                  <pre>{{ reversedCompileLog }}</pre>
                </div>
              </div>
              <div class="compare-half">
                <div class="log-header">
                  <span>对比结果</span>
                </div>
                <div class="compare-content">
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
                    点击"对比输出"按钮查看结果
                  </div>
                </div>
              </div>
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

      <!-- Copyright footer -->
      <div class="copyright-footer">
        © 2024 Frank Guo | 上海诶爱科技有限公司
      </div>

      <!-- 用户管理对话框 -->
      <UserManagement
        v-model="userManagementVisible"
        :online-users="onlineUsersWithDetails"
        :current-user="currentUser.username"
        :session-id="sessionId"
        @kick-user="handleKickUser"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { VideoPlay, Download, Upload, Check, Plus, Minus, User, Lock } from '@element-plus/icons-vue'
import CodeEditor from './components/CodeEditor.vue'
import ChatPanel from './components/ChatPanel.vue'
import UserManagement from './components/UserManagement.vue'

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
const showModules = ref(false)
const fontSize = ref(14)
const inputData = ref('')
const outputData = ref('')
const compileLog = ref('等待编译...\n')
const compiling = ref(false)
const answerData = ref('') // 标准答案
const compareResult = ref(null) // 对比结果

// 聊天相关
const chatMessages = ref([])

// 用户管理相关
const userManagementVisible = ref(false)

// 组件引用
const codeEditor = ref(null)
const chatPanel = ref(null)
const logContent = ref(null)

// 在线用户详情(用于用户管理)
const onlineUsersWithDetails = computed(() => {
  return onlineUsers.value
})

// 倒序显示编译日志（最新的在上面）
const reversedCompileLog = computed(() => {
  const lines = compileLog.value.split('\n')
  return lines.reverse().join('\n')
})

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

// 显示用户管理
const showUserManagement = () => {
  userManagementVisible.value = true
}

// 踢出用户
const handleKickUser = (username) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'kick_user',
      data: { username }
    }))
  }
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

    case 'user_kicked':
      // 用户被踢出
      if (message.data && message.data.username) {
        const kickedUser = message.data.username
        chatMessages.value.push({
          type: 'system',
          message: `${kickedUser} 被管理员踢出房间`,
          timestamp: message.timestamp
        })

        // 如果被踢的是当前用户，强制登出
        if (kickedUser === currentUser.value.username) {
          ElMessage.warning('您已被管理员踢出房间')
          setTimeout(() => {
            if (ws) {
              ws.close()
            }
            isLoggedIn.value = false
            sessionId.value = ''
            currentUser.value = { username: '', displayName: '' }
          }, 1000)
        }
      }
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
        displayName: message.displayName || message.username,
        message: message.data.message || '',
        fileUrl: message.data.fileUrl,
        fileName: message.data.fileName,
        fileSize: message.data.fileSize,
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

  // 浏览器通知
  if (document.hidden && message.type === 'chat') {
    sendBrowserNotification(message)
  }
}

// 发送浏览器通知
const sendBrowserNotification = (message) => {
  if (Notification.permission === 'granted') {
    const notification = new Notification(`来自 ${message.displayName || message.username} 的消息`, {
      body: message.data.message || '发送了一个文件',
      icon: '/favicon.ico'
    })
    
    notification.onclick = () => {
      window.focus()
      notification.close()
    }
  }
}

// 字体大小控制
const increaseFontSize = () => {
  if (fontSize.value < 30) fontSize.value += 2
}

const decreaseFontSize = () => {
  if (fontSize.value > 10) fontSize.value -= 2
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
  // 请求通知权限
  if (Notification.permission !== 'granted' && Notification.permission !== 'denied') {
    Notification.requestPermission()
  }
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
  position: relative;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('/login-bg.png');
  background-size: cover;
  background-position: center;
  z-index: -1;
}

.login-card {
  width: 420px;
  border: none !important;
  border-radius: 16px !important;
}

.glass-effect {
  background: rgba(255, 255, 255, 0.1) !important;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37) !important;
  border: 1px solid rgba(255, 255, 255, 0.18) !important;
}

.card-header {
  text-align: center;
  padding: 10px 0;
}

.gradient-text {
  background: linear-gradient(135deg, #00f260 0%, #0575e6 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 28px;
  font-weight: 800;
  margin: 0;
  letter-spacing: 1px;
}

.subtitle {
  color: rgba(255, 255, 255, 0.8);
  font-size: 12px;
  margin: 5px 0 0;
  letter-spacing: 2px;
  text-transform: uppercase;
}

.login-form {
  padding: 20px 10px;
}

.glass-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 8px 15px;
  transition: all 0.3s ease;
}

.glass-input :deep(.el-input__wrapper:hover),
.glass-input :deep(.el-input__wrapper.is-focus) {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.5);
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.1);
}

.glass-input :deep(.el-input__inner) {
  color: white;
  height: 30px;
}

.glass-input :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.6);
}

.login-btn {
  width: 100%;
  height: 45px;
  border-radius: 8px;
  background: linear-gradient(90deg, #00f260 0%, #0575e6 100%);
  border: none;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
  margin-top: 10px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(5, 117, 230, 0.4);
}

.login-tips {
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  margin-top: 10px;
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
  padding: 0;
  overflow: hidden;
  height: 100%;
}

.code-editor-container {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: bold;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.online-users {
  display: flex;
  align-items: center;
  font-size: 14px;
}

/* 输入输出标答三列 */
.io-section-three {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
  margin-top: 10px;
}

.io-column {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.area-header {
  font-weight: bold;
  margin-bottom: 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  order: 2;
}

.io-column .el-textarea {
  order: 1;
}

/* 操作按钮 */
.action-buttons {
  margin: 10px 0;
  text-align: center;
}

/* 日志和对比结果区 */
.log-compare-section {
  margin-top: 10px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  height: 150px;
}

.log-half, .compare-half {
  display: flex;
  flex-direction: column;
  min-height: 0;
  justify-content: flex-end;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  margin-bottom: 5px;
  order: 2;
}

.log-content {
  flex: 1;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  overflow-y: auto;
  overflow-x: hidden;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  order: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.compare-content {
  flex: 1;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  overflow-y: auto;
  overflow-x: hidden;
  order: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

/* 聊天区 */
.chat-section {
  border-left: 1px solid #ddd;
  background: #fafafa;
}

/* 对比结果区 */
.compare-result {
  height: 100%;
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

/* Copyright footer */
.copyright-footer {
  text-align: center;
  padding: 10px;
  background: #f5f5f5;
  color: #666;
  font-size: 12px;
  border-top: 1px solid #ddd;
}
</style>
