<template>
  <div class="role-select-container">
    <!-- 粒子背景容器 -->
    <div id="particles-js-role" class="particles-container"></div>
    
    <div class="role-select-box">
      <h2 class="title">选择角色</h2>
      <div class="role-options">
        <div 
          class="role-card" 
          :class="{ active: selectedRole === 'user' }"
          @click="selectRole('user')"
        >
          <el-icon class="icon"><User /></el-icon>
          <h3>数据使用者</h3>
          <p>上传数据、查询并浏览店铺信息</p>
        </div>
        
        <div 
          class="role-card" 
          :class="{ active: selectedRole === 'owner' }"
          @click="selectRole('owner')"
        >
          <el-icon class="icon"><Shop /></el-icon>
          <h3>数据拥有者</h3>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Shop } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const selectedRole = ref('')
const loading = ref(false)

// 初始化粒子背景
const initParticles = () => {
  // 动态加载 particles.js
  const script = document.createElement('script')
  script.src = '/js/particles.min.js'
  script.onload = () => {
    // @ts-ignore
    if (window.particlesJS) {
      // @ts-ignore
      window.particlesJS('particles-js-role', {
        "particles": {
          "number": {
            "value": 60,
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
            "value": 0.4,
            "random": true,
            "anim": {
              "enable": true,
              "speed": 1,
              "opacity_min": 0.1,
              "sync": false
            }
          },
          "size": {
            "value": 4,
            "random": true,
            "anim": {
              "enable": true,
              "speed": 2,
              "size_min": 0.1,
              "sync": false
            }
          },
          "line_linked": {
            "enable": true,
            "distance": 200,
            "color": "#ffffff",
            "opacity": 0.3,
            "width": 1
          },
          "move": {
            "enable": true,
            "speed": 3,
            "direction": "none",
            "random": true,
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
              "mode": "grab"
            },
            "onclick": {
              "enable": true,
              "mode": "push"
            },
            "resize": true
          },
          "modes": {
            "grab": {
              "distance": 300,
              "line_linked": {
                "opacity": 0.8
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

const selectRole = (role: string) => {
  selectedRole.value = role
}

const handleConfirm = async () => {
  if (!selectedRole.value) return
  
  loading.value = true
  try {
    // 存储角色信息
    localStorage.setItem('role', selectedRole.value)
    
    // 根据角色跳转到对应的仪表板
    if (selectedRole.value === 'user') {
      await router.push('/user-dashboard')
    } else {
      await router.push('/owner-dashboard')
    }
    
    ElMessage.success('角色选择成功')
  } catch (error) {
    ElMessage.error('跳转失败，请重试')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initParticles()
})
</script>

<style scoped>
.role-select-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

.role-select-box {
  position: relative;
  z-index: 2;
  width: 600px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  text-align: center;
  transition: all 0.3s ease;
}

.role-select-box:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.title {
  font-size: 28px;
  color: #333;
  margin-bottom: 30px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
  background: rgba(255, 255, 255, 0.85);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(4px);
}

.role-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  background: rgba(255, 255, 255, 0.95);
}

.role-card.active {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.15);
  box-shadow: 0 4px 20px rgba(64, 158, 255, 0.3);
}

.icon {
  font-size: 40px;
  color: #409EFF;
  margin-bottom: 15px;
  transition: all 0.3s ease;
}

.role-card:hover .icon {
  transform: scale(1.1);
}

.role-card h3 {
  font-size: 20px;
  color: #333;
  margin: 0 0 10px;
  font-weight: 600;
}

.role-card p {
  font-size: 14px;
  color: #666;
  margin: 0;
  line-height: 1.4;
}

.confirm-btn {
  width: 200px;
  height: 40px;
  font-size: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}
</style> 