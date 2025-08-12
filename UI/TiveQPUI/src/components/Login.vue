<template>
  <div class="login-container">
    <!-- 粒子背景容器 -->
    <div id="particles-js" class="particles-container"></div>
    
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
})

// 初始化粒子背景
const initParticles = () => {
  // 动态加载 particles.js
  const script = document.createElement('script')
  script.src = '/js/particles.min.js'
  script.onload = () => {
    // @ts-ignore
    if (window.particlesJS) {
      // @ts-ignore
      window.particlesJS('particles-js', {
        "particles": {
          "number": {
            "value": 80,
            "density": {
              "enable": true,
              "value_area": 800
            }
          },
          "color": {
            "value": "#ffffff"
          },
          "shape": {
            "type": "circle",
            "stroke": {
              "width": 0,
              "color": "#000000"
            }
          },
          "opacity": {
            "value": 0.3,
            "random": false,
            "anim": {
              "enable": false,
              "speed": 1,
              "opacity_min": 0.1,
              "sync": false
            }
          },
          "size": {
            "value": 3,
            "random": true,
            "anim": {
              "enable": false,
              "speed": 40,
              "size_min": 0.1,
              "sync": false
            }
          },
          "line_linked": {
            "enable": true,
            "distance": 150,
            "color": "#ffffff",
            "opacity": 0.2,
            "width": 1
          },
          "move": {
            "enable": true,
            "speed": 4,
            "direction": "none",
            "random": false,
            "straight": false,
            "out_mode": "out",
            "attract": {
              "enable": false,
              "rotateX": 600,
              "rotateY": 1200
            }
          }
        },
        "interactivity": {
          "detect_on": "canvas",
          "events": {
            "onhover": {
              "enable": true,
              "mode": "repulse"
            },
            "onclick": {
              "enable": true,
              "mode": "push"
            },
            "resize": true
          },
          "modes": {
            "grab": {
              "distance": 400,
              "line_linked": {
                "opacity": 1
              }
            },
            "bubble": {
              "distance": 400,
              "size": 40,
              "duration": 2,
              "opacity": 8,
              "speed": 3
            },
            "repulse": {
              "distance": 200
            },
            "push": {
              "particles_nb": 4
            },
            "remove": {
              "particles_nb": 2
            }
          }
        },
        "retina_detect": true
      })
    }
  }
  document.head.appendChild(script)
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid, fields) => {
    if (valid) {
      loading.value = true
      try {
        // TODO: 调用后端登录API
        // 模拟API调用
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // 存储token
        localStorage.setItem('token', 'mock-token')
        
        // 清除可能存在的角色信息
        localStorage.removeItem('role')
        
        // 跳转到角色选择页面
        await router.push('/role-select')
        
        ElMessage.success('登录成功')
      } catch (error) {
        ElMessage.error('登录失败，请重试')
      } finally {
        loading.value = false
      }
    } else {
      console.log('验证失败:', fields)
    }
  })
}

onMounted(() => {
  initParticles()
})
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  position: relative;
  overflow: hidden;
}

.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.login-box {
  position: relative;
  z-index: 2;
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.login-box:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.login-form {
  text-align: center;
}

.title {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
  font-size: 16px;
  color: #666;
  margin-bottom: 30px;
}

.form {
  margin-top: 20px;
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  background: rgba(255, 255, 255, 0.95);
}

:deep(.el-input__wrapper.is-focus) {
  background: #ffffff;
  border-color: #667eea;
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
  height: 40px;
  font-size: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
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