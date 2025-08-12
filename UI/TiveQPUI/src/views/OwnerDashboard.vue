<template>
  <div class="dashboard-container">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <el-button-group>
        <el-button type="primary" @click="loadShopData">
          <el-icon><Upload /></el-icon> 加载数据
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon> 导出数据
        </el-button>
      </el-button-group>
    </div>

    <div class="main-content">
      <!-- 数据上传区域 -->
      <el-card class="upload-card">
        <template #header>
          <div class="upload-header">
            <el-icon><Upload /></el-icon>
            <span>数据文件上传</span>
          </div>
        </template>
        
        <div class="upload-content">
          <div class="upload-section">
            <el-upload
              class="upload-area"
              :auto-upload="false"
              :limit="1"
              :on-change="handleFileChange"
              :file-list="fileList"
              :disabled="uploading"
            >
              <div class="upload-placeholder">
                <el-icon class="upload-icon"><FolderOpened /></el-icon>
                <div class="upload-text">
                  <p>点击选择文件或拖拽文件到此处</p>
                  <p class="upload-hint">支持与 20k.txt 相同格式，每行 8 段用 ** 分隔</p>
                </div>
              </div>
            </el-upload>
            
            <el-button 
              type="primary" 
              class="upload-btn" 
              :loading="uploading" 
              :disabled="!selectedFile || uploading"
              @click="submitUpload"
            >
              {{ uploading ? '处理中...' : '开始上传并构建索引' }}
            </el-button>
          </div>
          
          <!-- 进度显示区域 -->
          <div v-if="uploading || uploadProgress.length > 0" class="progress-section">
            <h4>处理进度</h4>
            <div class="progress-list">
              <div 
                v-for="(step, index) in uploadProgress" 
                :key="index"
                class="progress-item"
                :class="{ 
                  'completed': step.status === 'success',
                  'error': step.status === 'error',
                  'pending': step.status === 'pending'
                }"
              >
                <el-icon class="progress-icon">
                  <Check v-if="step.status === 'success'" />
                  <Close v-else-if="step.status === 'error'" />
                  <Loading v-else-if="step.status === 'pending'" />
                  <CircleCheck v-else />
                </el-icon>
                <span class="progress-text">{{ step.message }}</span>
                <span v-if="step.status === 'pending'" class="progress-time">{{ step.time }}</span>
              </div>
            </div>
            
            <el-progress 
              v-if="uploading"
              :percentage="progressPercentage" 
              :status="progressStatus"
              class="progress-bar"
            />
          </div>
        </div>
      </el-card>

      <!-- 统计卡片 -->
      <div class="stat-cards">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="stat-header">
                  <el-icon><Shop /></el-icon>
                  <span>总店铺数</span>
                </div>
              </template>
              <div class="stat-value">{{ statistics.totalShops }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="stat-header">
                  <el-icon><Location /></el-icon>
                  <span>覆盖城市</span>
                </div>
              </template>
              <div class="stat-value">{{ statistics.totalCities }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="stat-header">
                  <el-icon><Menu /></el-icon>
                  <span>店铺类型</span>
                </div>
              </template>
              <div class="stat-value">{{ statistics.totalTypes }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="stat-header">
                  <el-icon><Clock /></el-icon>
                  <span>营业中</span>
                </div>
              </template>
              <div class="stat-value">{{ statistics.openShops }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 筛选条件 -->
      <el-card class="filter-card">
        <el-form :model="filterForm" :inline="true">
          <el-form-item label="店铺类型">
            <el-select 
              v-model="filterForm.types" 
              multiple 
              placeholder="选择类型" 
              clearable
              :popper-append-to-body="false"
              class="custom-select"
            >
              <el-option v-for="type in shopTypes" :key="type" :label="type" :value="type" />
            </el-select>
          </el-form-item>

          <el-form-item label="营业时间">
            <el-time-picker
              v-model="filterForm.timeRange"
              is-range
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              format="HH:mm"
            />
          </el-form-item>

          <el-form-item label="城市">
            <el-select 
              v-model="filterForm.cities" 
              multiple 
              placeholder="选择城市" 
              clearable
              :popper-append-to-body="false"
              class="custom-select"
            >
              <el-option v-for="city in cities" :key="city" :label="city" :value="city" />
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleFilter">
              <el-icon><Search /></el-icon> 筛选
            </el-button>
            <el-button @click="resetFilter">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 数据表格 -->
      <el-card class="table-card">
        <el-table
          :data="tableData"
          style="width: 100%"
          v-loading="tableLoading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="type" label="店铺类型" width="120" />
          <el-table-column prop="city" label="城市" width="120" />
          <el-table-column prop="address" label="位置" show-overflow-tooltip>
            <template #default="{ row }">
              {{ `(${row.lat}, ${row.lng})` }}
            </template>
          </el-table-column>
          <el-table-column label="营业时间" width="200">
            <template #default="{ row }">
              {{ formatTime(row.hourStart, row.minStart) }} - {{ formatTime(row.hourClose, row.minClose) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button size="small" @click="handleEdit(row)">
                  <el-icon><Edit /></el-icon>
                </el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Upload,
  Download,
  Shop,
  Location,
  Menu,
  Clock,
  Search,
  Edit,
  Delete,
  FolderOpened,
  Check,
  Close,
  Loading,
  CircleCheck
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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

interface ProgressStep {
  message: string
  status: 'pending' | 'success' | 'error'
  time?: string
}

// 上传相关
const fileList = ref<any[]>([])
const selectedFile = ref<File | null>(null)
const uploading = ref(false)
const uploadProgress = ref<ProgressStep[]>([])

const progressPercentage = computed(() => {
  if (uploadProgress.value.length === 0) return 0
  const completed = uploadProgress.value.filter(step => step.status === 'success').length
  return Math.round((completed / uploadProgress.value.length) * 100)
})

const progressStatus = computed(() => {
  const hasError = uploadProgress.value.some(step => step.status === 'error')
  return hasError ? 'exception' : 'success'
})

const handleFileChange = (_file: any, files: any[]) => {
  fileList.value = files
  selectedFile.value = files && files.length ? (files[0].raw as File) : null
  // 重置进度
  uploadProgress.value = []
}

const submitUpload = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  
  // 初始化进度
  uploadProgress.value = [
    { message: '准备上传文件...', status: 'pending', time: new Date().toLocaleTimeString() },
    { message: '文件上传中...', status: 'pending' },
    { message: '验证文件格式...', status: 'pending' },
    { message: '构建索引中...', status: 'pending' },
    { message: '完成索引构建', status: 'pending' }
  ]
  
  uploading.value = true
  try {
    // 步骤1: 文件上传
    uploadProgress.value[0].status = 'success'
    uploadProgress.value[0].time = new Date().toLocaleTimeString()
    uploadProgress.value[1].status = 'pending'
    uploadProgress.value[1].time = new Date().toLocaleTimeString()
    
    const form = new FormData()
    form.append('file', selectedFile.value)
    
    // 模拟上传延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const resp = await fetch('http://localhost:8080/api/upload', {
      method: 'POST',
      body: form
    })
    
    // 步骤2: 上传完成
    uploadProgress.value[1].status = 'success'
    uploadProgress.value[1].time = new Date().toLocaleTimeString()
    uploadProgress.value[2].status = 'pending'
    uploadProgress.value[2].time = new Date().toLocaleTimeString()
    
    if (!resp.ok) {
      uploadProgress.value[2].status = 'error'
      uploadProgress.value[2].time = new Date().toLocaleTimeString()
      throw new Error(`上传失败: ${resp.status}`)
    }
    
    // 步骤3: 验证完成
    uploadProgress.value[2].status = 'success'
    uploadProgress.value[2].time = new Date().toLocaleTimeString()
    uploadProgress.value[3].status = 'pending'
    uploadProgress.value[3].time = new Date().toLocaleTimeString()
    
    // 模拟构建延迟
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // 步骤4: 构建完成
    uploadProgress.value[3].status = 'success'
    uploadProgress.value[3].time = new Date().toLocaleTimeString()
    uploadProgress.value[4].status = 'success'
    uploadProgress.value[4].time = new Date().toLocaleTimeString()
    
    ElMessage.success('上传并构建索引成功')
    await loadStatistics()
    await fetchData()
  } catch (e) {
    // 标记当前步骤为错误
    const currentStep = uploadProgress.value.find(step => step.status === 'pending')
    if (currentStep) {
      currentStep.status = 'error'
      currentStep.time = new Date().toLocaleTimeString()
    }
    ElMessage.error(e instanceof Error ? e.message : '上传失败')
  } finally {
    uploading.value = false
  }
}

// 统计数据
const statistics = reactive({
  totalShops: 0,
  totalCities: 0,
  totalTypes: 0,
  openShops: 0
})

// 筛选表单
const filterForm = reactive({
  types: [],
  timeRange: [],
  cities: []
})

// 表格数据
const tableData = ref<Shop[]>([])
const tableLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 选项数据
const shopTypes = ref(['餐厅', '超市', '便利店', '药店'])
const cities = ref(['北京', '上海', '广州', '深圳'])

// 获取数据
const fetchData = async () => {
  tableLoading.value = true
  try {
    // 构建查询参数
    const params = new URLSearchParams({
      page: currentPage.value.toString(),
      pageSize: pageSize.value.toString(),
    })

    if (filterForm.types.length === 1) {
      params.append('type', filterForm.types[0])
    }
    if (filterForm.cities.length === 1) {
      params.append('city', filterForm.cities[0])
    }
    if (filterForm.timeRange && filterForm.timeRange.length === 2) {
      const time = new Date(filterForm.timeRange[0])
      params.append('time', `${time.getHours()}:${time.getMinutes()}`)
    }

    const response = await fetch(`http://localhost:8080/api/shops?${params}`)
    if (!response.ok) {
      throw new Error('加载数据失败')
    }
    
    const data = await response.json()
    tableData.value = data.data
    total.value = data.total
    
    ElMessage.success('数据加载成功')
  } catch (error) {
    ElMessage.error('加载数据失败：' + (error instanceof Error ? error.message : String(error)))
  } finally {
    tableLoading.value = false
  }
}

// 加载店铺数据
const loadShopData = async () => {
  await Promise.all([
    loadStatistics(),
    fetchData()
  ])
}

// 处理筛选
const handleFilter = () => {
  currentPage.value = 1
  fetchData()
}

// 重置筛选
const resetFilter = () => {
  filterForm.types = []
  filterForm.timeRange = []
  filterForm.cities = []
  handleFilter()
}

// 处理分页
const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}

// 处理选择
const handleSelectionChange = (val: Shop[]) => {
  console.log('selected:', val)
}

// 处理编辑
const handleEdit = (row: Shop) => {
  console.log('edit:', row)
}

// 处理删除
const handleDelete = (row: Shop) => {
  ElMessageBox.confirm(
    `确定要删除该店铺吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // TODO: 调用后端API
      await new Promise(resolve => setTimeout(resolve, 500))
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

// 处理导出
const handleExport = () => {
  ElMessage.success('开始导出数据')
}

// 格式化时间
const formatTime = (hour: number, min: number) => {
  return `${hour.toString().padStart(2, '0')}:${min.toString().padStart(2, '0')}`
}

// 修改统计数据获取方式
const loadStatistics = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/shops/stats')
    if (!response.ok) {
      throw new Error('加载统计数据失败')
    }
    const data = await response.json()
    statistics.totalShops = data.totalShops
    statistics.totalCities = data.totalCities
    statistics.totalTypes = data.totalTypes
    statistics.openShops = data.openShops
    
    // 更新选项数据
    shopTypes.value = data.types
    cities.value = data.cities
  } catch (error) {
    ElMessage.error('加载统计数据失败：' + (error instanceof Error ? error.message : String(error)))
  }
}

// 初始化
onMounted(() => {
  loadStatistics()
  fetchData()
})
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  padding: 20px;
}

.action-bar {
  margin-bottom: 20px;
}

.main-content {
  flex: 1;
  overflow-y: auto;
}

.upload-card {
  margin-bottom: 20px;
}

.upload-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.upload-content {
  padding: 20px 0;
}

.upload-section {
  margin-bottom: 20px;
}

.upload-area {
  margin-bottom: 16px;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s;
}

.upload-placeholder:hover {
  border-color: #409EFF;
  background: #f0f9ff;
}

.upload-icon {
  font-size: 48px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.upload-text {
  text-align: center;
}

.upload-text p {
  margin: 4px 0;
  color: #606266;
}

.upload-hint {
  font-size: 12px;
  color: #909399;
}

.upload-btn {
  width: 100%;
  height: 40px;
}

.progress-section {
  border-top: 1px solid #e4e7ed;
  padding-top: 20px;
}

.progress-section h4 {
  margin: 0 0 16px 0;
  color: #303133;
}

.progress-list {
  margin-bottom: 16px;
}

.progress-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-radius: 4px;
}

.progress-item.completed {
  background: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}

.progress-item.error {
  background: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.progress-item.pending {
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.progress-icon {
  font-size: 16px;
}

.progress-text {
  flex: 1;
  font-size: 14px;
}

.progress-time {
  font-size: 12px;
  color: #909399;
}

.progress-bar {
  margin-top: 8px;
}

.stat-cards {
  margin-bottom: 20px;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  color: #409EFF;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  background: white;
  border-radius: 8px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-card__header) {
  padding: 12px 20px;
}

:deep(.el-card__body) {
  padding: 20px;
}

.custom-select {
  min-width: 200px;
  
  :deep(.el-select__tags) {
    flex-wrap: nowrap;
    overflow-x: auto;
    max-width: calc(100% - 80px) !important;
    scrollbar-width: none;
    -ms-overflow-style: none;
    &::-webkit-scrollbar {
      display: none;
    }
  }

  :deep(.el-select__tags-text) {
    display: inline-block;
    max-width: none;
  }

  :deep(.el-input__wrapper) {
    padding-right: 65px !important;
  }

  :deep(.el-input__suffix) {
    right: 8px;
    display: flex;
    align-items: center;
    gap: 12px;
    position: absolute;
    height: 100%;
    top: 0;
  }

  :deep(.el-input__suffix-inner) {
    display: flex !important;
    align-items: center;
    gap: 12px;
    position: static !important;
  }

  :deep(.el-select__caret.el-icon) {
    margin: 0;
    position: static;
    order: 2;
    z-index: 1;
  }

  :deep(.el-select__clear.el-icon) {
    position: static;
    order: 1;
    margin-right: 4px;
    z-index: 1;
  }
}
</style> 