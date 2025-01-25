<template>
  <div class="dashboard-container">
    <!-- 左侧查询面板 -->
    <div class="search-panel">
      <!-- 位置信息 -->
      <div class="location-card">
        <h3><el-icon><Location /></el-icon> 位置信息</h3>
        <div class="location-options">
          <el-tabs v-model="locationMode" class="location-tabs">
            <el-tab-pane label="自动定位" name="auto">
              <div class="location-info" v-if="currentLocation">
                <p>纬度: {{ currentLocation.latitude }}</p>
                <p>经度: {{ currentLocation.longitude }}</p>
              </div>
              <el-button @click="refreshLocation" :loading="locationLoading" class="location-btn">
                <el-icon><Refresh /></el-icon> 刷新位置
              </el-button>
            </el-tab-pane>
            <el-tab-pane label="选择城市" name="manual">
              <el-select 
                v-model="selectedCity" 
                placeholder="请选择城市"
                class="city-select"
                filterable
              >
                <el-option 
                  v-for="city in cities" 
                  :key="city.value" 
                  :label="city.label" 
                  :value="city.value"
                />
              </el-select>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>

      <!-- 查询条件 -->
      <div class="filter-card">
        <el-form :model="searchForm" label-position="top">
          <el-form-item label="店铺类型">
            <el-select v-model="searchForm.type" placeholder="选择店铺类型" clearable class="full-width">
              <el-option v-for="type in shopTypes" :key="type" :label="type" :value="type" />
            </el-select>
          </el-form-item>

          <el-form-item label="查询时间" class="time-select-item">
            <el-radio-group v-model="searchForm.timeMode" class="time-mode">
              <el-radio label="current">当前时间</el-radio>
              <el-radio label="specific">指定时间</el-radio>
            </el-radio-group>
            
            <div class="time-picker-wrapper" v-if="searchForm.timeMode === 'specific'">
              <el-time-picker
                v-model="searchForm.specificTime"
                format="HH:mm"
                placeholder="选择时间"
                class="full-width"
              />
            </div>
          </el-form-item>

          <el-form-item label="搜索范围" class="range-select-item">
            <div class="range-display">{{ searchForm.radius }}km</div>
            <el-slider
              v-model="searchForm.radius"
              :min="0.5"
              :max="5"
              :step="0.5"
              :marks="{
                0.5: '0.5km',
                1: '1km',
                3: '3km',
                5: '5km'
              }"
              class="range-slider"
            />
          </el-form-item>

          <el-button type="primary" @click="handleSearch" :loading="searchLoading" class="search-btn">
            <el-icon><Search /></el-icon> 搜索
          </el-button>
        </el-form>
      </div>

      <!-- 查询结果列表 -->
      <div class="result-card" v-if="searchResults.length">
        <h3>查询结果 ({{ searchResults.length }})</h3>
        <el-table :data="searchResults" style="width: 100%" height="300">
          <el-table-column prop="name" label="店铺名称" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="distance" label="距离" width="100">
            <template #default="{ row }">
              {{ formatDistance(row.distance) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <div class="no-result" v-else-if="hasSearched">
        <el-empty description="暂无符合条件的店铺" />
      </div>
    </div>

    <!-- 右侧地图 -->
    <div class="map-container">
      <ShopMap 
        :shops="searchResults"
        :currentLocation="currentLocation"
        :selectedShop="selectedShop"
        @shop-click="handleShopClick"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Location, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import ShopMap from '../components/ShopMap.vue'

interface Location {
  latitude: number
  longitude: number
}

interface Shop {
  id: string
  name: string
  type: string
  location: Location
  distance?: number
}

// 状态
const currentLocation = ref<Location | null>(null)
const locationLoading = ref(false)
const searchLoading = ref(false)
const hasSearched = ref(false)
const searchResults = ref<Shop[]>([])
const selectedShop = ref<Shop | null>(null)
const shopTypes = ref(['餐厅', '超市', '便利店', '药店'])
const locationMode = ref('auto')
const selectedCity = ref('')

// 城市数据
const cities = [
  { label: '北京', value: 'beijing' },
  { label: '上海', value: 'shanghai' },
  { label: '广州', value: 'guangzhou' },
  { label: '深圳', value: 'shenzhen' },
  // 可以添加更多城市
]

// 表单数据
const searchForm = reactive({
  type: '',
  timeMode: 'current',
  specificTime: null,
  radius: 1
})

// 获取位置
const refreshLocation = async () => {
  locationLoading.value = true
  try {
    const position = await new Promise<GeolocationPosition>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(
        resolve,
        reject,
        {
          enableHighAccuracy: true,  // 启用高精度
          timeout: 5000,            // 5秒超时
          maximumAge: 0             // 不使用缓存
        }
      )
    })
    
    currentLocation.value = {
      latitude: position.coords.latitude,
      longitude: position.coords.longitude
    }
    ElMessage.success('位置更新成功')
  } catch (error) {
    console.error('位置获取错误:', error)
    if (error instanceof GeolocationPositionError) {
      switch (error.code) {
        case error.PERMISSION_DENIED:
          ElMessage.error('获取位置失败：用户拒绝了位置请求')
          break
        case error.POSITION_UNAVAILABLE:
          ElMessage.error('获取位置失败：位置信息不可用')
          break
        case error.TIMEOUT:
          ElMessage.error('获取位置失败：请求超时')
          break
        default:
          ElMessage.error('获取位置失败：未知错误')
      }
    } else {
      ElMessage.error('获取位置失败，请检查浏览器定位权限')
    }
  } finally {
    locationLoading.value = false
  }
}

// 搜索
const handleSearch = async () => {
  if (locationMode.value === 'auto' && !currentLocation.value) {
    ElMessage.warning('请先获取位置信息')
    return
  }

  if (locationMode.value === 'manual' && !selectedCity.value) {
    ElMessage.warning('请选择城市')
    return
  }

  searchLoading.value = true
  hasSearched.value = true
  
  try {
    // TODO: 调用后端API
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    searchResults.value = [
      {
        id: '1',
        name: '示例店铺1',
        type: '餐厅',
        location: locationMode.value === 'auto' ? {
          latitude: currentLocation.value!.latitude + 0.001,
          longitude: currentLocation.value!.longitude + 0.001
        } : {
          latitude: 39.9042,
          longitude: 116.4074
        },
        distance: 200
      }
    ]
  } catch (error) {
    ElMessage.error('搜索失败，请重试')
  } finally {
    searchLoading.value = false
  }
}

// 处理店铺点击
const handleShopClick = (shop: Shop) => {
  selectedShop.value = shop
}

// 格式化距离
const formatDistance = (meters: number) => {
  return meters < 1000 ? `${meters}m` : `${(meters / 1000).toFixed(1)}km`
}

// 初始化获取位置
if (locationMode.value === 'auto') {
  refreshLocation()
}
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
  display: flex;
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
}

.search-panel {
  width: 400px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  overflow-y: auto;
}

.location-card,
.filter-card,
.result-card {
  margin-bottom: 20px;
  padding: 20px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.location-card:hover,
.filter-card:hover,
.result-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.15);
}

.location-options {
  margin-top: 15px;
}

.location-tabs {
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  padding: 10px;
}

.location-info {
  margin: 15px 0;
  padding: 15px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.location-btn,
.search-btn {
  width: 100%;
  margin-top: 10px;
  background: linear-gradient(135deg, #5b86e5 0%, #36d1dc 100%);
  border: none;
  height: 40px;
  transition: all 0.3s ease;
}

.location-btn:hover,
.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(91, 134, 229, 0.4);
}

.city-select {
  width: 100%;
  margin-top: 10px;
}

.time-select-item {
  margin-bottom: 25px;
}

.time-mode {
  margin-bottom: 15px;
  width: 100%;
  display: flex;
  justify-content: space-around;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  padding: 10px;
}

.time-picker-wrapper {
  padding: 10px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  margin-top: 10px;
}

.range-select-item {
  margin-bottom: 25px;
}

.range-display {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 15px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  padding: 10px;
}

.range-slider {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
}

:deep(.el-slider__marks-text) {
  margin-top: 15px;
  color: #606266;
}

:deep(.el-tabs__nav) {
  width: 100%;
  display: flex;
}

:deep(.el-tabs__item) {
  flex: 1;
  text-align: center;
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
  border-color: #409EFF;
}

:deep(.el-table) {
  background: transparent;
}

:deep(.el-table__header) {
  background: rgba(255, 255, 255, 0.8);
}

:deep(.el-table__body tr) {
  background: rgba(255, 255, 255, 0.6);
  transition: all 0.3s ease;
}

:deep(.el-table__body tr:hover) {
  background: rgba(255, 255, 255, 0.9);
}

.map-container {
  flex: 1;
  height: 100%;
  position: relative;
}

.full-width {
  width: 100%;
}

.no-result {
  margin-top: 40px;
  text-align: center;
}

@media screen and (max-width: 768px) {
  .dashboard-container {
    flex-direction: column;
  }

  .search-panel {
    width: 100%;
    height: 50vh;
  }

  .map-container {
    height: 50vh;
  }
}
</style> 