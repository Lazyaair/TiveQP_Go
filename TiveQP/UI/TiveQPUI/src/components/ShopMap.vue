<template>
  <div class="map-container">
    <div class="search-panel">
      <div class="panel-header">
        <h2>店铺查询</h2>
        <p class="stats">共 {{ filteredShopsCount }} 家店铺</p>
      </div>
      
      <el-input
        v-model="searchQuery"
        placeholder="搜索店铺..."
        class="search-input"
        :prefix-icon="Search"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button @click="handleSearch">搜索</el-button>
        </template>
      </el-input>
      
      <el-collapse>
        <el-collapse-item title="高级筛选" name="1">
          <el-form :model="filters" class="filter-form" label-position="top">
            <el-form-item label="店铺类型">
              <el-select 
                v-model="filters.type" 
                placeholder="选择店铺类型" 
                clearable
                class="full-width"
              >
                <el-option 
                  v-for="type in uniqueTypes" 
                  :key="type" 
                  :label="type" 
                  :value="type" 
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="城市">
              <el-select 
                v-model="filters.city" 
                placeholder="选择城市" 
                clearable
                class="full-width"
              >
                <el-option 
                  v-for="city in uniqueCities" 
                  :key="city" 
                  :label="city" 
                  :value="city" 
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="营业时间">
              <div class="time-range">
                <el-time-select
                  v-model="filters.timeStart"
                  start="00:00"
                  step="00:30"
                  end="23:30"
                  placeholder="开始时间"
                  clearable
                />
                <span class="time-separator">至</span>
                <el-time-select
                  v-model="filters.timeEnd"
                  start="00:00"
                  step="00:30"
                  end="23:30"
                  placeholder="结束时间"
                  clearable
                />
              </div>
            </el-form-item>
          </el-form>
        </el-collapse-item>
      </el-collapse>
    </div>
    
    <div class="map-controls">
      <el-button-group>
        <el-button @click="handleZoomIn" :icon="Plus" circle></el-button>
        <el-button @click="handleZoomOut" :icon="Minus" circle></el-button>
      </el-button-group>
    </div>
    
    <div id="map" ref="mapContainer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { Search, Plus, Minus } from '@element-plus/icons-vue'
import { shops, type Shop } from '../data/shopLoader'

interface Filters {
  type: string
  city: string
  timeStart: string
  timeEnd: string
}

const mapContainer = ref<HTMLElement | null>(null)
const searchQuery = ref('')
const map = ref<L.Map | null>(null)
const markers = ref<L.Marker[]>([])
const shopList = ref<Shop[]>([])

const filters = reactive<Filters>({
  type: '',
  city: '',
  timeStart: '',
  timeEnd: ''
})

// 计算过滤后的店铺数量
const filteredShopsCount = computed(() => {
  return getFilteredShops().length
})

// 获取过滤后的店铺列表
const getFilteredShops = () => {
  return shopList.value.filter((shop: Shop) => {
    const matchSearch = searchQuery.value 
      ? shop.type.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        shop.city.toLowerCase().includes(searchQuery.value.toLowerCase())
      : true

    const matchType = filters.type ? shop.type === filters.type : true
    const matchCity = filters.city ? shop.city === filters.city : true
    const matchTime = filters.timeStart && filters.timeEnd 
      ? new Date(`1970-01-01T${filters.timeStart}`).getTime() >= new Date(`1970-01-01T${shop.openTime}`).getTime() &&
        new Date(`1970-01-01T${filters.timeEnd}`).getTime() <= new Date(`1970-01-01T${shop.closeTime}`).getTime()
      : true

    return matchSearch && matchType && matchCity && matchTime
  })
}

const addMarkers = (shops: Shop[]) => {
  markers.value.forEach(marker => marker.remove())
  markers.value = []

  shops.forEach(shop => {
    if (map.value) {
      const marker = L.marker([shop.location.lat, shop.location.lng])
        .bindPopup(`
          <h3>${shop.type}</h3>
          <p>${shop.city}</p>
          <p>营业时间: ${shop.openTime} - ${shop.closeTime}</p>
        `)
        .addTo(map.value)
      markers.value.push(marker)
    }
  })
}

// 地图缩放控制
const handleZoomIn = () => {
  if (map.value) {
    map.value.setZoom(map.value.getZoom() + 1)
  }
}

const handleZoomOut = () => {
  if (map.value) {
    map.value.setZoom(map.value.getZoom() - 1)
  }
}

onMounted(() => {
  if (!mapContainer.value) return

  map.value = L.map(mapContainer.value, {
    zoomControl: false
  }).setView([28.4498736, -81.4863524], 10)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map.value)

  // 使用模拟数据
  shopList.value = [
    {
      id: '1',
      type: 'Hair Salons',
      city: 'ORLANDO',
      location: {
        lat: 28.4498736,
        lng: -81.4863524
      },
      openTime: '10:00',
      closeTime: '18:00'
    }
  ]
  
  addMarkers(shopList.value)
})

const handleSearch = () => {
  const filteredShops = getFilteredShops()
  addMarkers(filteredShops)
}

// 获取所有唯一的店铺类型和城市
const uniqueTypes = computed(() => 
  Array.from(new Set(shopList.value.map(shop => shop.type)))
)

const uniqueCities = computed(() => 
  Array.from(new Set(shopList.value.map(shop => shop.city)))
)

defineExpose({
  searchQuery,
  filters,
  handleSearch
})
</script>

<style scoped>
.map-container {
  position: relative;
  height: 100vh;
  background: #f5f7fa;
}

.search-panel {
  position: absolute;
  left: 20px;
  top: 20px;
  width: 350px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  z-index: 1000;
  transition: all 0.3s ease;
}

.panel-header {
  margin-bottom: 20px;
}

.panel-header h2 {
  margin: 0 0 8px;
  color: #303133;
  font-size: 20px;
}

.stats {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.search-input {
  margin-bottom: 15px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

.filter-form {
  padding: 15px 0;
}

.full-width {
  width: 100%;
}

.time-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-separator {
  color: #909399;
}

.map-controls {
  position: absolute;
  right: 20px;
  top: 20px;
  z-index: 1000;
}

:deep(.el-button) {
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.9);
}

#map {
  height: 100%;
  width: 100%;
}

/* 自定义地图标记样式 */
:deep(.leaflet-popup-content-wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

:deep(.leaflet-popup-content) {
  margin: 12px;
}

:deep(.leaflet-popup-content h3) {
  margin: 0 0 8px;
  color: #303133;
}

:deep(.leaflet-popup-content p) {
  margin: 4px 0;
  color: #606266;
}

@media screen and (max-width: 768px) {
  .search-panel {
    width: calc(100% - 40px);
    margin: 0 20px;
  }
}
</style>