<template>
  <div class="map-container">
    <!-- 用户界面 -->
    <div v-if="currentRole === 'user'" class="search-panel">
      <div class="panel-header">
        <h2>店铺查询</h2>
        <p class="stats">共 {{ filteredShopsCount }} 家店铺</p>
      </div>
      
      <el-form :model="userFilters" class="filter-form" label-position="top">
        <el-form-item label="店铺类型">
          <el-select 
            v-model="userFilters.type" 
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

        <el-form-item label="查询时间">
          <el-radio-group v-model="userFilters.timeMode" class="time-mode">
            <el-radio label="current">当前时间</el-radio>
            <el-radio label="specific">指定时间</el-radio>
          </el-radio-group>
          
          <el-time-picker
            v-if="userFilters.timeMode === 'specific'"
            v-model="userFilters.specificTime"
            format="HH:mm"
            placeholder="选择时间"
            class="full-width"
          />
        </el-form-item>

        <el-form-item label="城市">
          <el-select 
            v-model="userFilters.city" 
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

        <el-button type="primary" @click="handleUserSearch" class="search-btn">
          查询
        </el-button>
      </el-form>

      <!-- 查询结果列表 -->
      <div class="result-card" v-if="searchResults.length">
        <h3>查询结果 ({{ searchResults.length }})</h3>
        <el-table 
          ref="tableRef"
          :data="searchResults" 
          style="width: 100%" 
          height="300"
          @row-click="handleRowClick"
          :highlight-current-row="true"
        >
          <el-table-column 
            prop="type" 
            label="类型" 
            width="120"
            class-name="clickable-column"
          />
          <el-table-column 
            label="开始时间" 
            width="120"
            class-name="clickable-column"
          >
            <template #default="{ row }">
              {{ formatTime(row.hourStart, row.minStart) }}
            </template>
          </el-table-column>
          <el-table-column 
            label="结束时间" 
            width="120"
            class-name="clickable-column"
          >
            <template #default="{ row }">
              {{ formatTime(row.hourClose, row.minClose) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 店主界面 -->
    <div v-else-if="currentRole === 'owner'" class="filter-panel">
      <div class="panel-header">
        <h2>店铺管理</h2>
        <p class="stats">共 {{ filteredShopsCount }} 家店铺</p>
      </div>

      <el-form :model="ownerFilters" class="filter-form" label-position="top">
        <el-form-item label="店铺类型">
          <el-select 
            v-model="ownerFilters.type" 
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

        <el-form-item label="时间段">
          <div class="time-range">
            <el-time-select
              v-model="ownerFilters.timeStart"
              start="00:00"
              step="00:30"
              end="23:30"
              placeholder="开始时间"
              clearable
            />
            <span class="time-separator">至</span>
            <el-time-select
              v-model="ownerFilters.timeEnd"
              start="00:00"
              step="00:30"
              end="23:30"
              placeholder="结束时间"
              clearable
            />
          </div>
        </el-form-item>

        <el-form-item label="城市">
          <el-select 
            v-model="ownerFilters.city" 
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

        <el-button type="primary" @click="handleOwnerFilter" class="filter-btn">
          筛选
        </el-button>
      </el-form>
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
import { ref, onMounted, reactive, computed, defineExpose, watch } from 'vue'
import { useRoute } from 'vue-router'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { Plus, Minus } from '@element-plus/icons-vue'

interface Shop {
  type: string
  city: string  // 这个字段将同时作为id使用
  lat: number
  lng: number
  hourStart: number
  minStart: number
  hourClose: number
  minClose: number
}

const route = useRoute()
const currentRole = computed(() => route.query.role as string)

interface UserFilters {
  type: string
  timeMode: 'current' | 'specific'
  specificTime: string | null
  city: string
}

interface OwnerFilters {
  type: string
  timeStart: string
  timeEnd: string
  city: string
}

const mapContainer = ref<HTMLElement | null>(null)
const map = ref<L.Map | null>(null)
const markers = ref<L.Marker[]>([])
const shopList = ref<Shop[]>([])
const searchResults = ref<Shop[]>([])
const tableRef = ref()

const userFilters = reactive<UserFilters>({
  type: '',
  timeMode: 'current',
  specificTime: null,
  city: ''
})

const ownerFilters = reactive<OwnerFilters>({
  type: '',
  timeStart: '',
  timeEnd: '',
  city: ''
})

// 计算过滤后的店铺数量
const filteredShopsCount = computed(() => {
  return currentRole.value === 'user' ? searchResults.value.length : getFilteredShops().length
})

// 获取当前时间是否在营业时间内
const isShopOpen = (shop: Shop, checkTime: string) => {
  const time = new Date(`1970-01-01T${checkTime}`)
  const open = new Date(`1970-01-01T${shop.hourStart}:${shop.minStart}`)
  const close = new Date(`1970-01-01T${shop.hourClose}:${shop.minClose}`)
  return time >= open && time <= close
}

// 用户查询处理
const handleUserSearch = () => {
  const currentTime = new Date().toLocaleTimeString('en-US', { 
    hour12: false, 
    hour: '2-digit', 
    minute: '2-digit'
  })
  
  const checkTime = userFilters.timeMode === 'current' 
    ? currentTime 
    : userFilters.specificTime || ''

  searchResults.value = shopList.value.filter(shop => {
    const matchType = !userFilters.type || shop.type === userFilters.type
    const matchCity = !userFilters.city || shop.city === userFilters.city
    const matchTime = checkTime ? isShopOpen(shop, checkTime) : true

    return matchType && matchCity && matchTime
  })

  addMarkers(searchResults.value)
}

// 店主筛选处理
const handleOwnerFilter = () => {
  const filteredShops = getFilteredShops()
  addMarkers(filteredShops)
}

// 获取过滤后的店铺列表（店主视角）
const getFilteredShops = () => {
  return shopList.value.filter((shop: Shop) => {
    const matchType = !ownerFilters.type || shop.type === ownerFilters.type
    const matchCity = !ownerFilters.city || shop.city === ownerFilters.city
    const matchTime = ownerFilters.timeStart && ownerFilters.timeEnd 
      ? new Date(`1970-01-01T${ownerFilters.timeStart}`).getTime() >= new Date(`1970-01-01T${shop.hourStart}:${shop.minStart}`).getTime() &&
        new Date(`1970-01-01T${ownerFilters.timeEnd}`).getTime() <= new Date(`1970-01-01T${shop.hourClose}:${shop.minClose}`).getTime()
      : true

    return matchType && matchCity && matchTime
  })
}

// 处理行点击
const handleRowClick = (row: Shop, column: any, event: Event) => {
  console.log('表格行被点击')
  console.log('点击的行数据:', JSON.stringify(row))
  
  if (!row) {
    console.error('行数据为空')
    return
  }
  
  if (!map.value) {
    console.error('地图未初始化')
    return
  }
  
  // 设置当前选中行
  if (tableRef.value) {
    console.log('设置当前选中行')
    tableRef.value.setCurrentRow(row)
  }
  
  // 聚焦到选中的店铺
  console.log('调用 focusShop')
  focusShop(row)
}

// 聚焦到特定店铺
const focusShop = (shop: Shop) => {
  console.log('focusShop called with shop:', shop)
  
  if (!map.value) {
    console.error('Map is not initialized in focusShop')
    return
  }
  
  if (!shop || typeof shop.lat !== 'number' || typeof shop.lng !== 'number') {
    console.error('Invalid shop data:', shop)
    return
  }
  
  try {
    const currentZoom = map.value.getZoom()
    console.log('Current map zoom level:', currentZoom)
    console.log('Flying to coordinates:', [shop.lat, shop.lng])
    
    // 先关闭所有已打开的弹窗
    markers.value.forEach(marker => {
      console.log('Closing popup for marker')
      marker.closePopup()
    })
    
    // 使用平滑动画效果
    map.value.flyTo([shop.lat, shop.lng], 15, {
      duration: 1.5,
      easeLinearity: 0.25
    })
    console.log('Map flyTo executed')
    
    // 找到对应的标记并打开弹窗
    let markerFound = false
    console.log('Total markers:', markers.value.length)
    
    markers.value.forEach(marker => {
      const latLng = marker.getLatLng()
      console.log('Checking marker position:', latLng, 'against shop position:', [shop.lat, shop.lng])
      
      if (Math.abs(latLng.lat - shop.lat) < 0.0001 && Math.abs(latLng.lng - shop.lng) < 0.0001) {
        console.log('Matching marker found, opening popup')
        marker.openPopup()
        markerFound = true
      }
    })
    
    if (!markerFound) {
      console.warn('No matching marker found for shop:', shop)
    }
  } catch (error) {
    console.error('Error in focusShop:', error)
  }
}

const formatTime = (hour: number, min: number) => 
  `${hour.toString().padStart(2, '0')}:${min.toString().padStart(2, '0')}`

function initMap() {
  const mapElement = document.getElementById('map')
  console.log('Map element:', mapElement)
  
  if (!mapElement) {
    console.error('Map element not found')
    return
  }

  try {
    map.value = L.map(mapElement).setView([39.8283, -98.5795], 4)
    console.log('Map instance created:', map.value)
    
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors'
    }).addTo(map.value)
    
    console.log('Tile layer added to map')
  } catch (error) {
    console.error('Error in map initialization:', error)
  }
}

function addMarkers(shops: Shop[]) {
  console.log('Adding markers for shops:', shops)
  if (!map.value) {
    console.error('Map not initialized when adding markers')
    return
  }

  try {
    console.log('Clearing existing markers:', markers.value.length)
    markers.value.forEach((marker: L.Marker) => marker.remove())
    markers.value = []

    shops.forEach(shop => {
      console.log('Creating marker for shop:', shop)
      const marker = L.marker([shop.lat, shop.lng])
        .bindPopup(`
          <div class="popup-content">
            <h3>${shop.city}</h3>
            <p>类型: ${shop.type}</p>
            <p>营业时间: ${formatTime(shop.hourStart, shop.minStart)} - ${formatTime(shop.hourClose, shop.minClose)}</p>
          </div>
        `)
        .addTo(map.value!)
      
      markers.value.push(marker)
      console.log('Marker added successfully at:', [shop.lat, shop.lng])
    })

    console.log('Total markers after addition:', markers.value.length)
  } catch (error) {
    console.error('Error in addMarkers:', error)
  }
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

// 获取所有唯一的店铺类型和城市
const uniqueTypes = computed(() => 
  Array.from(new Set(shopList.value.map(shop => shop.type)))
)

const uniqueCities = computed(() => 
  Array.from(new Set(shopList.value.map(shop => shop.city)))
)

onMounted(() => {
  console.log('Component mounted')
  setTimeout(() => {
    console.log('Initializing map after delay')
    initMap()
  }, 100) // 给DOM一些时间来渲染
})

// 更新店铺数据和标记
const updateShops = (shops: Shop[]) => {
  console.log('updateShops called with:', shops)
  
  const processedShops = shops.map(shop => ({
    ...shop,
    city: shop.city.toUpperCase()
  }))
  
  console.log('Processed shops:', processedShops)
  shopList.value = processedShops
  searchResults.value = processedShops

  console.log('Adding markers for processed shops')
  addMarkers(processedShops)

  if (processedShops.length > 0) {
    console.log('Fitting map to bounds')
    const latLngs = processedShops.map(shop => [shop.lat, shop.lng])
    const bounds = L.latLngBounds(latLngs)
    map.value?.fitBounds(bounds, { padding: [50, 50] })
  }
}

// 监听店铺列表变化
watch(shopList, (newShops) => {
  console.log('Shop list changed:', newShops)
  if (newShops.length > 0) {
    console.log('Adding markers for updated shop list')
    addMarkers(newShops)
  }
}, { deep: true })

// 监听搜索结果变化
watch(searchResults, (newResults) => {
  console.log('Search results changed:', newResults)
  if (newResults.length > 0) {
    console.log('Adding markers for search results')
    addMarkers(newResults)
  }
}, { deep: true })

// 导出更新方法供父组件调用
defineExpose({
  updateShops
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
  background: rgba(255, 255, 255, 0.8);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  z-index: 1000;
  transition: all 0.3s ease;
}

.panel-header {
  margin-bottom: 20px;
}

.panel-header h2 {
  margin: 0 0 8px;
  color: rgba(48, 49, 51, 0.9);
  font-size: 20px;
  font-weight: 600;
}

.stats {
  margin: 0;
  color: rgba(144, 147, 153, 0.9);
  font-size: 14px;
}

.search-input {
  margin-bottom: 15px;
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px rgba(64, 158, 255, 0.5);
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

.time-range .el-select {
  flex: 1;
  min-width: 120px;
}

:deep(.el-select__wrapper) {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

:deep(.el-select__selection) {
  padding: 0 4px;
}

:deep(.el-select__input) {
  margin: 0;
  padding: 0;
}

:deep(.el-select__placeholder) {
  padding: 0 4px;
}

.time-separator {
  color: rgba(144, 147, 153, 0.9);
  padding: 0 4px;
  flex-shrink: 0;
}

:deep(.el-select-dropdown__list) {
  min-width: 120px;
}

.map-controls {
  position: absolute;
  right: 20px;
  top: 20px;
  z-index: 1000;
}

:deep(.el-button) {
  background: rgba(64, 158, 255, 0.9);
  border: none;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

:deep(.el-button:hover) {
  background: rgba(64, 158, 255, 1);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
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

/* 折叠面板样式 */
:deep(.el-collapse) {
  border: none;
  background: transparent;
}

:deep(.el-collapse-item__header) {
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  color: rgba(48, 49, 51, 0.9);
}

:deep(.el-collapse-item__content) {
  background: transparent;
  color: rgba(48, 49, 51, 0.9);
}

:deep(.el-form-item__label) {
  color: rgba(48, 49, 51, 0.9);
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .search-panel {
    width: calc(100% - 40px);
    margin: 0 20px;
  }

  .time-range {
    flex-direction: column;
    gap: 12px;
  }

  .time-range .el-select {
    width: 100%;
  }

  .time-separator {
    display: none;
  }
}

/* 表单样式 */
:deep(.filter-form) {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 表单项样式 */
:deep(.el-form-item) {
  margin-bottom: 16px;
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

:deep(.el-form-item__label) {
  color: rgba(48, 49, 51, 0.9);
  font-weight: 500;
  padding-bottom: 4px;
}

/* Select 下拉框样式 */
:deep(.el-select .el-input__wrapper) {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: none;
}

:deep(.el-select .el-input__wrapper:hover) {
  background: rgba(255, 255, 255, 0.8);
}

:deep(.el-select .el-input__wrapper.is-focus) {
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 0 0 1px rgba(64, 158, 255, 0.5);
}

/* 下拉选项样式 */
:deep(.el-select-dropdown) {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

:deep(.el-select-dropdown__item) {
  background: transparent;
}

:deep(.el-select-dropdown__item.hover) {
  background: rgba(64, 158, 255, 0.1);
}

:deep(.el-select-dropdown__item.selected) {
  background: rgba(64, 158, 255, 0.2);
  color: #409EFF;
}

/* 时间选择器样式 */
.time-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-range .el-select {
  flex: 1;
  min-width: 120px;
}

:deep(.el-select__selection) {
  background: transparent;
}

:deep(.el-select__placeholder) {
  color: rgba(128, 128, 128, 0.8);
}

/* 折叠面板样式 */
:deep(.el-collapse-item__header) {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 0 16px;
  height: 48px;
  line-height: 48px;
  color: rgba(48, 49, 51, 0.9);
  font-weight: 500;
}

:deep(.el-collapse-item__header:hover) {
  background: rgba(255, 255, 255, 0.7);
}

:deep(.el-collapse-item__wrap) {
  background: transparent;
  border: none;
}

:deep(.el-collapse-item__content) {
  padding: 16px 0 0;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .time-range {
    flex-direction: column;
    gap: 12px;
  }

  .time-range .el-select {
    width: 100%;
  }

  .time-separator {
    display: none;
  }
}

.search-results {
  margin-top: 20px;
  border-top: 1px solid #eee;
  padding-top: 15px;
}

.search-results h3 {
  margin: 0 0 15px;
  color: #333;
}

.shop-item {
  padding: 12px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.shop-item:hover {
  background: rgba(255, 255, 255, 1);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.shop-item.active {
  background: rgba(64, 158, 255, 0.1);
  border-color: #409EFF;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.shop-item h4 {
  margin: 0 0 8px;
  color: #303133;
  font-size: 16px;
}

.shop-item p {
  margin: 4px 0;
  color: #606266;
  font-size: 14px;
}

.time-mode {
  margin-bottom: 10px;
  width: 100%;
  display: flex;
  gap: 20px;
}

.filter-panel {
  position: absolute;
  left: 20px;
  top: 20px;
  width: 350px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  z-index: 1000;
}

:deep(.el-table__row) {
  cursor: pointer !important;
}

:deep(.el-table__body) {
  cursor: pointer !important;
}

.result-card {
  margin-top: 20px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.result-card h3 {
  margin: 0 0 15px;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

:deep(.clickable-row) {
  cursor: pointer !important;
}

:deep(.el-table) {
  --el-table-row-hover-bg-color: var(--el-color-primary-light-9);
  --el-table-current-row-bg-color: var(--el-color-primary-light-8);
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: var(--el-table-row-hover-bg-color) !important;
}

:deep(.el-table__row.current-row) {
  background-color: var(--el-table-current-row-bg-color) !important;
  color: var(--el-color-primary);
}

:deep(.clickable-column) {
  cursor: pointer;
}

:deep(.el-table__body td) {
  cursor: pointer;
}
</style>