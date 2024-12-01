<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-form">
        <h1 class="title">TiveQP</h1>
        <p class="subtitle">组4</p>
        
        <el-form 
          ref="formRef"
          :model="loginForm" 
          :rules="rules"
          class="form"
        >
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username"
              placeholder="用户名"
              :prefix-icon="User"
              clearable
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input 
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              :prefix-icon="Lock"
              show-password
              clearable
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              @click="handleLogin"
              :loading="loading"
              class="login-button"
              round
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
})

const handleLogin = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 模拟登录成功
    if (loginForm.username === 'admin' && loginForm.password === '123456') {
      localStorage.setItem('token', 'dummy-token')
      ElMessage.success('登录成功')
      router.push('/map')
    } else {
      ElMessage.error('登录失败，用户名或密码错误')
    }
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error('登录失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 60%);
  pointer-events: none;
}

.login-box {
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.85);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transform: translateY(0);
  position: relative;
  z-index: 1;
}

.login-box:hover {
  transform: translateY(-5px);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
}

.title {
  text-align: center;
  font-size: 32px;
  background: linear-gradient(135deg, #5b86e5, #36d1dc);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 8px;
  font-weight: 600;
  letter-spacing: 2px;
}

.subtitle {
  text-align: center;
  color: #5b86e5;
  margin: 0 0 30px;
  font-size: 16px;
  font-weight: 500;
}

.form {
  margin-top: 30px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
  padding: 8px 15px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid transparent;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.el-input__wrapper:hover) {
  background: rgba(243, 244, 246, 0.9);
  border-color: rgba(102, 126, 234, 0.1);
}

:deep(.el-input__wrapper.is-focus) {
  background: #ffffff;
  border-color: #667eea;
  box-shadow: 0 0 0 1px #5b86e5;
}

:deep(.el-input__inner) {
  color: #374151;
}

:deep(.el-input__inner::placeholder) {
  color: #9ca3af;
}

:deep(.el-input__prefix-inner .el-icon) {
  color: #667eea;
  font-size: 18px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  margin-top: 20px;
  background: linear-gradient(135deg, #5b86e5 0%, #36d1dc 100%);
  border: none;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(91, 134, 229, 0.4);
}

.login-button:active {
  transform: translateY(0);
}

.login-button::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
  pointer-events: none;
}

@media screen and (max-width: 480px) {
  .login-box {
    width: 90%;
    margin: 0 20px;
    padding: 30px 20px;
  }

  .title {
    font-size: 32px;
  }
}

/* 添加动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-form {
  animation: fadeIn 0.6s ease-out;
}
</style> 