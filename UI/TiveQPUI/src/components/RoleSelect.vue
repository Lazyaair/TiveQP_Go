<template>
  <div class="role-select-container">
    <div class="role-select-box">
      <h2 class="title">选择角色</h2>
      <div class="role-options">
        <div 
          class="role-card" 
          :class="{ active: selectedRole === 'user' }"
          @click="selectRole('user')"
        >
          <el-icon class="icon"><User /></el-icon>
          <h3>用户</h3>
          <p>查询和浏览店铺信息</p>
        </div>
        
        <div 
          class="role-card" 
          :class="{ active: selectedRole === 'owner' }"
          @click="selectRole('owner')"
        >
          <el-icon class="icon"><Shop /></el-icon>
          <h3>店主</h3>
          <p>管理和查看店铺数据</p>
        </div>
      </div>
      
      <el-button 
        type="primary" 
        class="confirm-btn" 
        :loading="loading"
        :disabled="!selectedRole"
        @click="handleConfirm"
      >
        确认
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { User, Shop } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const selectedRole = ref('')
const loading = ref(false)

const selectRole = (role: string) => {
  selectedRole.value = role
}

const handleConfirm = async () => {
  if (selectedRole.value) {
    loading.value = true
    try {
      localStorage.setItem('userRole', selectedRole.value)
      await router.push({
        path: '/map',
        query: { role: selectedRole.value }
      })
    } catch (error) {
      ElMessage.error('页面跳转失败，请重试')
    } finally {
      loading.value = false
    }
  }
}
</script>

<style scoped>
.role-select-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
}

.role-select-box {
  width: 600px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  text-align: center;
}

.title {
  font-size: 28px;
  color: #333;
  margin-bottom: 30px;
}

.role-options {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-bottom: 30px;
}

.role-card {
  flex: 1;
  max-width: 200px;
  padding: 30px 20px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.8);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.role-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.role-card.active {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

.icon {
  font-size: 40px;
  color: #409EFF;
  margin-bottom: 15px;
}

.role-card h3 {
  font-size: 20px;
  color: #333;
  margin: 0 0 10px;
}

.role-card p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.confirm-btn {
  width: 200px;
  height: 40px;
  font-size: 16px;
}
</style> 