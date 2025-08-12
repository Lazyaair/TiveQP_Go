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
                <p>纬度: {{ currentLocation?.latitude }}</p>
                <p>经度: {{ currentLocation?.longitude }}</p>
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
                  :key="city"
                  :label="city" 
                  :value="city"
                />
              </el-select>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>

      <!-- 查询条件 -->
      <div class="filter-card">
        <el-form :model="searchForm" label-position="top">
          <el-form-item label="店铺类型（可多选）">
            <el-select
              v-model="searchForm.types"
              placeholder="请选择店铺类型（可多选）"
              class="full-width"
              filterable
              multiple
              collapse-tags
              collapse-tags-tooltip
            >
              <el-option v-for="type in shopTypes" :key="type" :label="type" :value="type" />
            </el-select>
            <div v-if="searchForm.types.length" class="per-type-settings">
              <div class="per-type-row" v-for="t in searchForm.types" :key="'row-' + t">
                <span class="per-type-label">{{ t }}</span>
                <el-input-number
                  v-model.number="radiusByType[t]"
                  :min="RADIUS_MIN"
                  :max="RADIUS_MAX"
                  :step="1"
                  :controls="false"
                  class="per-type-input"
                  placeholder="距离 (km)"
                  @change="() => validateRadius(t)"
                  @blur="() => validateRadius(t)"
                />
                <span class="per-type-unit">km</span>
              </div>
              <div class="hint">每个类型可设置独立范围；未设置时使用下方默认范围</div>
            </div>
          </el-form-item>

          <el-form-item label="最大搜索数量">
            <el-input-number
              v-model="searchForm.maxShops"
              :min="1"
              :max="100"
              :step="1"
              class="full-width"
              placeholder="请输入最大搜索数量"
              :controls="false"
            />
          </el-form-item>

          <el-form-item label="搜索时间" class="time-select-item">
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

          <el-form-item label="默认搜索范围" class="range-select-item">
            <div class="range-display">{{ searchForm.radius }}km</div>
            <el-input
              v-model.number="searchForm.radius"
              type="number"
              placeholder="请输入搜索范围（km）"
              class="full-width"
              @change="() => validateRadius()"
              @blur="() => validateRadius()"
            >
              <template #suffix>km</template>
            </el-input>
          </el-form-item>

          <el-button type="primary" @click="handleSearch" :loading="searchLoading" class="search-btn">
            <el-icon><Search /></el-icon> 搜索
          </el-button>
          <div v-if="lastQueryTime" class="query-time">上次查询时间：{{ lastQueryTime }}</div>
        </el-form>
      </div>

      <!-- 陷门显示按钮 -->
      <div class="trapdoor-card" v-if="trapdoorInfo">
        <el-button 
          type="primary" 
          class="trapdoor-btn"
          @click="showTrapdoorDialog = true"
        >
          <el-icon><Key /></el-icon>
          显示陷门
        </el-button>
      </div>

      <!-- 查询结果列表 -->
      <div class="result-card" v-if="searchResults.length">
        <h3>查询结果 ({{ searchResults.length }})</h3>
        <el-table :data="searchResults" style="width: 100%" height="300">
          <el-table-column prop="type" label="类型" width="120" />
          <el-table-column label="开始时间" width="120">
            <template #default="{ row }">
              {{ formatTime(row.hourStart, row.minStart) }}
            </template>
          </el-table-column>
          <el-table-column label="结束时间" width="120">
            <template #default="{ row }">
              {{ formatTime(row.hourClose, row.minClose) }}
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
        ref="mapRef"
      />
    </div>

    <!-- 陷门内容弹窗 - 移到最外层 -->
    <teleport to="#app">
      <el-dialog
        v-model="showTrapdoorDialog"
        title="查询陷门"
        width="80%"
        :close-on-click-modal="false"
        class="trapdoor-dialog"
        append-to-body
      >
        <div class="trapdoor-content">
          <div class="trapdoor-title">
            <h3>陷门 T</h3>
          </div>
          <div class="trapdoor-matrices">
            <div v-for="(matrix, key) in parsedTrapdoor" :key="key" class="matrix-section">
              <div class="matrix-header">
                <h4>{{ key }}</h4>
                <div class="matrix-type">{{ getMatrixType(key) }}</div>
              </div>
              <div class="matrix-content">
                <div v-for="(row, rowIndex) in matrix" :key="rowIndex" class="matrix-row">
                  <div v-for="(cell, cellIndex) in row" :key="cellIndex" class="matrix-cell">
                    {{ cell }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-dialog>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Location, Refresh, Search, Key } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ShopMap from '../components/ShopMap.vue'

interface Location {
  latitude: number
  longitude: number
}

interface Shop {
  type: string
  city: string
  lat: number
  lng: number
  hourStart: number
  minStart: number
  hourClose: number
  minClose: number
}

// 状态
const currentLocation = ref<Location | null>(null)
const locationLoading = ref(false)
const searchLoading = ref(false)
const hasSearched = ref(false)
const searchResults = ref<Shop[]>([])
const selectedShop = ref<Shop | null>(null)
const shopTypes = ref<string[]>([])
const locationMode = ref('auto')
const selectedCity = ref('')
const cities = ref<string[]>([])
const trapdoorInfo = ref<string>('')
const showTrapdoorDialog = ref(false)
const lastQueryTime = ref<string>('')
const mapRef = ref()

// 表单数据
const searchForm = reactive({
  type: '',
  timeMode: 'current',
  specificTime: new Date(),
  city: '',
  radius: 1,
  maxShops: undefined,  // 使用undefined
  types: [] as string[]
})

const RADIUS_MIN = 1
const RADIUS_MAX = 5

// 每类型的范围设置
const radiusByType = reactive<Record<string, number>>({})

const validateRadius = (typeKey?: string): boolean => {
  const applyClamp = (val: number) => {
    const rounded = Math.round(val)
    return Math.min(Math.max(rounded, RADIUS_MIN), RADIUS_MAX)
  }

  if (typeKey) {
    const raw = Number(radiusByType[typeKey])
    if (isNaN(raw)) {
      ElMessageBox.alert(`请输入有效的范围（${RADIUS_MIN} - ${RADIUS_MAX} km）`, '范围无效')
      radiusByType[typeKey] = searchForm.radius
      return false
    }
    const clamped = applyClamp(raw)
    if (clamped !== raw) {
      ElMessageBox.alert(`搜索范围应在 ${RADIUS_MIN} - ${RADIUS_MAX} km 之间，已自动调整为 ${clamped} km`, '范围超出')
      radiusByType[typeKey] = clamped
      return true
    }
    radiusByType[typeKey] = clamped
    return true
  }

  // 校验默认范围
  const raw = Number(searchForm.radius)
  if (isNaN(raw)) {
    ElMessageBox.alert(`请输入有效的范围（${RADIUS_MIN} - ${RADIUS_MAX} km）`, '范围无效')
    searchForm.radius = RADIUS_MIN
    return false
  }
  const clamped = applyClamp(raw)
  if (clamped !== raw) {
    ElMessageBox.alert(`搜索范围应在 ${RADIUS_MIN} - ${RADIUS_MAX} km 之间，已自动调整为 ${clamped} km`, '范围超出')
    searchForm.radius = clamped
    return true
  }
  searchForm.radius = clamped
  return true
}

// 解析陷门数据
const parsedTrapdoor = computed(() => {
  if (!trapdoorInfo.value) return null
  try {
    const data = JSON.parse(trapdoorInfo.value)
    return {
      T1: JSON.parse(data.T1),
      T2: JSON.parse(data.T2),
      T3: JSON.parse(data.T3)
    }
  } catch (error) {
    console.error('解析陷门数据失败:', error)
    return null
  }
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

  // 校验搜索范围
  validateRadius()

  searchLoading.value = true
  hasSearched.value = true
  
  try {
    // 获取当前时间或指定时间
    const time = searchForm.timeMode === 'specific' && searchForm.specificTime 
      ? searchForm.specificTime 
      : new Date()
    
    // 选择查询的类型：优先使用多选；否则使用单选；都没有则用 ALL
    const selectedTypes: string[] = (searchForm.types && searchForm.types.length > 0)
      ? [...searchForm.types]
      : (searchForm.type ? [searchForm.type] : ['ALL'])

    const allResults: Shop[] = []

    // 依次发送请求并合并
    for (const typeItem of selectedTypes) {
      const effectiveRadius = typeItem !== 'ALL' && radiusByType[typeItem] ? radiusByType[typeItem] : searchForm.radius
      // 对每个类型也校验一次范围
      validateRadius(typeItem !== 'ALL' ? typeItem : undefined)
      const params = [
        typeItem,
        locationMode.value === 'auto' ? '自动获取的城市' : selectedCity.value,
        locationMode.value === 'auto' ? currentLocation.value!.latitude.toString() : '33.846335',
        locationMode.value === 'auto' ? currentLocation.value!.longitude.toString() : '-84.3635778',
        time.getHours().toString(),
        time.getMinutes().toString(),
        (searchForm.maxShops === undefined ? 1 : searchForm.maxShops).toString(),
        effectiveRadius.toString()
      ].join('**')

      console.log('Sending search request with params:', params)
      const response = await fetch(`http://localhost:8080/api/message?params=${encodeURIComponent(params)}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      console.log('Received data from API:', data)

      if (!data) {
        throw new Error('返回数据格式错误')
      }

      if (data.trapdoor) {
        const formattedTrapdoor = {
          T1: JSON.stringify(data.trapdoor.t1, null, 2),
          T2: JSON.stringify(data.trapdoor.t2, null, 2),
          T3: JSON.stringify(data.trapdoor.t3, null, 2)
        }
        // 显示最后一次查询的陷门
        trapdoorInfo.value = JSON.stringify(formattedTrapdoor, null, 2)
      }

      if (Array.isArray(data.shops)) {
        const processedResults = data.shops.map((shop: {
          type: string;
          city: string;
          lat: number;
          lng: number;
          hourStart: number;
          minStart: number;
          hourClose: number;
          minClose: number;
        }) => ({
          type: shop.type,
          city: shop.city,
          lat: shop.lat,
          lng: shop.lng,
          hourStart: shop.hourStart,
          minStart: shop.minStart,
          hourClose: shop.hourClose,
          minClose: shop.minClose
        }))
        allResults.push(...processedResults)
      }
    }

    // 去重合并结果
    const seen = new Set<string>()
    const merged: Shop[] = []
    for (const s of allResults) {
      const key = `${s.type}|${s.city}|${s.lat}|${s.lng}|${s.hourStart}|${s.minStart}|${s.hourClose}|${s.minClose}`
      if (!seen.has(key)) {
        seen.add(key)
        merged.push(s)
      }
    }

    console.log('Merged search results:', merged)
    searchResults.value = merged

    // 更新地图标记
    if (mapRef.value) {
      console.log('Updating map with merged results')
      mapRef.value.updateShops(merged)
    } else {
      console.error('Map reference not found')
    }

    if (merged.length === 0) {
      ElMessage.info('未找到符合条件的店铺')
    } else {
      ElMessage.success(`找到 ${merged.length} 家店铺（已合并）`)
    }
  } catch (error) {
    console.error('搜索错误:', error)
    ElMessage.error(error instanceof Error ? error.message : '搜索失败，请重试')
    searchResults.value = []
    trapdoorInfo.value = ''  // 清空陷门信息
  } finally {
    searchLoading.value = false
    // 记录查询时间（当前时间）
    const currentTime = new Date().toLocaleString()
    const previousTime = lastQueryTime.value
    lastQueryTime.value = currentTime
    
    // 显示查询完成弹窗
    ElMessageBox.alert(
      '<div style="text-align: center; padding: 20px;">' +
        '<h3 style="color: #409EFF; margin-bottom: 16px;">查询完成</h3>' +
        '<div style="margin-bottom: 12px;">' +
          '<strong>当前查询时间：</strong><br>' +
          '<span style="color: #67c23a;">' + currentTime + '</span>' +
        '</div>' +
        (previousTime ? '<div><strong>上次查询时间：</strong><br><span style="color: #909399;">' + previousTime + '</span></div>' : '') +
      '</div>',
      '查询时间信息',
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '确定',
        customClass: 'query-time-dialog'
      }
    )
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

// 格式化时间
const formatTime = (hour: number, minute: number) => {
  return `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
}

// 添加地图组件引用


// 添加加载选项数据的方法
const loadOptions = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/shops/stats')
    if (!response.ok) {
      throw new Error('加载数据失败')
    }
    const data = await response.json()
    
    // 更新选项数据
    shopTypes.value = data.types
    cities.value = data.cities
  } catch (error) {
    ElMessage.error('加载选项数据失败：' + (error instanceof Error ? error.message : String(error)))
  }
}

// 在初始化时调用
onMounted(() => {
  loadOptions()
  if (locationMode.value === 'auto') {
    refreshLocation()
  }
})

const getMatrixType = (key: string) => {
  switch (key) {
    case 'T1':
      return '类型'
    case 'T2':
      return '位置'
    case 'T3':
      return '时间'
    default:
      return ''
  }
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

.multi-type {
  margin-top: 10px;
}

.per-type-settings { margin-top: 10px; }
.per-type-row { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.per-type-label { min-width: 140px; color: #606266; }
.per-type-input { width: 120px; }
.per-type-unit { color: #909399; }
.hint { color: #909399; font-size: 12px; }

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

.query-time {
  margin-top: 8px;
  text-align: center;
  color: #606266;
  font-size: 12px;
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

.trapdoor-card {
  margin-bottom: 20px;
  padding: 20px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  text-align: center;
}

.trapdoor-btn {
  width: 100%;
  height: 40px;
  background: linear-gradient(135deg, #5b86e5 0%, #36d1dc 100%);
  border: none;
  transition: all 0.3s ease;
}

.trapdoor-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(91, 134, 229, 0.4);
}

.trapdoor-content {
  max-height: 70vh;
  overflow-y: auto;
  padding: 20px;
  background: #f5f7fa;
}

.trapdoor-title {
  text-align: center;
  margin-bottom: 30px;
  padding: 15px;
  background: linear-gradient(135deg, #5b86e5 0%, #36d1dc 100%);
  border-radius: 8px;
  color: white;
}

.trapdoor-title h3 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.trapdoor-matrices {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  padding: 20px;
}

.matrix-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.matrix-header {
  background: rgba(64, 158, 255, 0.1);
  padding: 15px;
  border-bottom: 1px solid rgba(64, 158, 255, 0.2);
}

.matrix-header h4 {
  margin: 0;
  color: #409EFF;
  font-size: 18px;
  font-weight: 600;
}

.matrix-type {
  font-size: 14px;
  color: #606266;
  margin-top: 5px;
}

.matrix-content {
  padding: 15px;
  background: white;
}

.matrix-row {
  display: flex;
  margin-bottom: 8px;
  gap: 8px;
}

.matrix-row:last-child {
  margin-bottom: 0;
}

.matrix-cell {
  flex: 1;
  background: rgba(64, 158, 255, 0.05);
  padding: 8px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  word-break: break-all;
  color: #606266;
  text-align: center;
}

@media screen and (max-width: 1200px) {
  .trapdoor-matrices {
    grid-template-columns: 1fr;
  }
}

.trapdoor-dialog {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #5b86e5 0%, #36d1dc 100%);
  padding: 20px;
  margin-right: 0;
}

:deep(.el-dialog__title) {
  color: white;
  font-size: 18px;
  font-weight: 600;
}

:deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
}

:deep(.el-dialog__body) {
  padding: 0;
  background: #f5f7fa;
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