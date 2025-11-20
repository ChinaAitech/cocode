<template>
  <el-dialog
    v-model="visible"
    title="用户管理"
    width="800px"
    :before-close="handleClose"
  >
    <el-tabs v-model="activeTab">
      <!-- 在线用户标签页 -->
      <el-tab-pane label="在线用户" name="online">
        <el-table :data="onlineUsers" style="width: 100%">
          <el-table-column prop="username" label="用户名" width="200" />
          <el-table-column label="状态" width="100">
            <template #default>
              <el-tag type="success">在线</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" v-if="isAdmin">
            <template #default="{ row }">
              <el-button
                v-if="row.username !== currentUser"
                size="small"
                type="danger"
                @click="kickUser(row.username)"
              >
                踢出
              </el-button>
              <el-tag v-else type="info">当前用户</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 用户清单标签页 -->
      <el-tab-pane label="用户清单" name="list" v-if="isAdmin">
        <div style="margin-bottom: 10px">
          <el-button type="primary" @click="showAddUserDialog" :icon="Plus">
            添加用户
          </el-button>
          <el-button @click="loadUserList" :icon="Refresh">
            刷新
          </el-button>
        </div>
        <el-table :data="userList" style="width: 100%">
          <el-table-column prop="username" label="用户名" width="150" />
          <el-table-column prop="displayName" label="显示名称" width="150" />
          <el-table-column prop="password" label="密码" width="150">
            <template #default="{ row }">
              <span>{{ row.showPassword ? row.password : '••••••••' }}</span>
              <el-button
                link
                size="small"
                @click="row.showPassword = !row.showPassword"
                style="margin-left: 5px"
              >
                {{ row.showPassword ? '隐藏' : '显示' }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="{ row }">
              <el-button size="small" @click="editUser(row)">编辑</el-button>
              <el-button
                size="small"
                type="danger"
                @click="deleteUser(row.username)"
                :disabled="row.username === 'admin'"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      :title="editingUser ? '编辑用户' : '添加用户'"
      width="400px"
    >
      <el-form :model="userForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input
            v-model="userForm.username"
            :disabled="!!editingUser"
            placeholder="请输入用户名"
          />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="userForm.displayName" placeholder="请输入显示名称" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">保存</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  onlineUsers: {
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

const emit = defineEmits(['update:modelValue', 'kick-user'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isAdmin = computed(() => props.currentUser === 'admin')
const activeTab = ref('online')
const userList = ref([])
const userDialogVisible = ref(false)
const editingUser = ref(null)
const userForm = ref({
  username: '',
  displayName: '',
  password: ''
})

// 踢出用户
const kickUser = (username) => {
  ElMessageBox.confirm(
    `确定要踢出用户 ${username} 吗？`,
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    emit('kick-user', username)
    ElMessage.success(`已踢出用户 ${username}`)
  }).catch(() => {
    // 取消操作
  })
}

// 加载用户列表
const loadUserList = async () => {
  try {
    const response = await fetch('/api/users', {
      headers: {
        'X-Session-ID': props.sessionId
      }
    })
    if (response.ok) {
      const data = await response.json()
      userList.value = data.users.map(u => ({ ...u, showPassword: false }))
    } else {
      ElMessage.error('加载用户列表失败')
    }
  } catch (error) {
    ElMessage.error('加载用户列表失败: ' + error.message)
  }
}

// 显示添加用户对话框
const showAddUserDialog = () => {
  editingUser.value = null
  userForm.value = {
    username: '',
    displayName: '',
    password: ''
  }
  userDialogVisible.value = true
}

// 编辑用户
const editUser = (user) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    displayName: user.displayName,
    password: user.password
  }
  userDialogVisible.value = true
}

// 保存用户
const saveUser = async () => {
  if (!userForm.value.username || !userForm.value.password) {
    ElMessage.error('用户名和密码不能为空')
    return
  }

  try {
    const url = editingUser.value ? '/api/users/update' : '/api/users/create'
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Session-ID': props.sessionId
      },
      body: JSON.stringify(userForm.value)
    })

    if (response.ok) {
      ElMessage.success(editingUser.value ? '用户更新成功' : '用户创建成功')
      userDialogVisible.value = false
      loadUserList()
    } else {
      const data = await response.json()
      ElMessage.error(data.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败: ' + error.message)
  }
}

// 删除用户
const deleteUser = (username) => {
  if (username === 'admin') {
    ElMessage.error('不能删除管理员账号')
    return
  }

  ElMessageBox.confirm(
    `确定要删除用户 ${username} 吗？此操作不可恢复。`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await fetch('/api/users/delete', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Session-ID': props.sessionId
        },
        body: JSON.stringify({ username })
      })

      if (response.ok) {
        ElMessage.success('用户删除成功')
        loadUserList()
      } else {
        const data = await response.json()
        ElMessage.error(data.message || '删除失败')
      }
    } catch (error) {
      ElMessage.error('删除失败: ' + error.message)
    }
  }).catch(() => {
    // 取消操作
  })
}

const handleClose = () => {
  visible.value = false
}

// 监听 tab 切换,加载用户列表
const handleTabChange = () => {
  if (activeTab.value === 'list' && isAdmin.value) {
    loadUserList()
  }
}

// 暴露方法给父组件
defineExpose({
  loadUserList
})
</script>

<style scoped>
.el-table {
  margin-top: 10px;
}
</style>
