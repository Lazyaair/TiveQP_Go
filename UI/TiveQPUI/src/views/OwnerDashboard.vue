<template>
  <div class="dashboard-container">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <el-button-group>
        <el-button type="primary" @click="showUploadDialog">
          <el-icon><Upload /></el-icon> 上传数据
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon> 导出数据
        </el-button>
      </el-button-group>
    </div>

    <div class="main-content">
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
            <el-select v-model="filterForm.types" multiple placeholder="选择类型" clearable>
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
            <el-select v-model="filterForm.cities" multiple placeholder="选择城市" clearable>
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
          <el-table-column prop="name" label="店铺名称" />
          <el-table-column prop="type" label="类型" width="120" />
          <el-table-column prop="city" label="城市" width="120" />
          <el-table-column prop="address" label="地址" show-overflow-tooltip />
          <el-table-column label="营业时间" width="200">
            <template #default="{ row }">
              {{ row.openTime }} - {{ row.closeTime }}
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

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传数据" width="500px">
      <el-upload
        class="upload-demo"
        drag
        action="/api/upload"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        accept=".xlsx,.csv"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽文件到此处或 <em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .xlsx 或 .csv 格式的文件
          </div>
        </template>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
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
  UploadFilled
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

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
const tableData = ref([])
const tableLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 选项数据
const shopTypes = ref(['餐厅', '超市', '便利店', '药店'])
const cities = ref(['北京', '上海', '广州', '深圳'])

// 上传对话框
const uploadDialogVisible = ref(false)

// 获取数据
const fetchData = async () => {
  tableLoading.value = true
  try {
    // TODO: 调用后端API
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    tableData.value = [
      {
        id: '1',
        name: '示例店铺1',
        type: '餐厅',
        city: '北京',
        address: '示例地址1',
        openTime: '09:00',
        closeTime: '22:00'
      }
    ]
    total.value = 100
    
    // 更新统计数据
    statistics.totalShops = 100
    statistics.totalCities = 4
    statistics.totalTypes = 4
    statistics.openShops = 80
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    tableLoading.value = false
  }
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
const handleSelectionChange = (val: any[]) => {
  console.log('selected:', val)
}

// 处理编辑
const handleEdit = (row: any) => {
  console.log('edit:', row)
}

// 处理删除
const handleDelete = (row: any) => {
  ElMessageBox.confirm(
    `确定要删除店铺 "${row.name}" 吗？`,
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

// 处理上传
const showUploadDialog = () => {
  uploadDialogVisible.value = true
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  uploadDialogVisible.value = false
  fetchData()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 初始化
fetchData()
</script>

<style scoped>
.dashboard-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.action-bar {
  padding: 16px;
  background: white;
  border-bottom: 1px solid #ebeef5;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
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
  margin-bottom: 20px;
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

.upload-demo {
  text-align: center;
}
</style> 