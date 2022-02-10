<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.info"
        placeholder="appId/租户名称"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter.native="handleFilter"
      />
      <el-button
        v-waves
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        style="margin-left: 20px"
        @click="handleFilter"
      >搜索</el-button>

      <router-link :to="'/app/appAddUser'">
        <el-button
          class="filter-item"
          style="margin-left: 20px;"
          type="primary"
          icon="el-icon-edit"
        >
          添加租户
        </el-button>
      </router-link>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" prop="id" align="center" width="50">
        <template slot-scope="{ row }">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="appId" min-width="50px">
        <template slot-scope="{row}">
          <span>{{ row.appId }}</span>
        </template>
      </el-table-column>
      <el-table-column label="租户名称" min-width="80px">
        <template slot-scope="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="密钥" min-width="120px">
        <template slot-scope="{ row }">
          <span>{{ row.secret }}</span>
        </template>
      </el-table-column>
      <el-table-column label="QPS" width="100px">
        <template slot-scope="{ row }">
          <span>{{ row.realQps }} / {{ row.qps }}</span>
        </template>
      </el-table-column>
      <el-table-column label="日请求数" width="110px">
        <template slot-scope="{ row }">
          <span>{{ row.realQpd }} / {{ row.qpd }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="操作"
        align="center"
        min-width="90"
        class-name="small-padding fixed-width"
      >
        <template slot-scope="{ row }">
          <router-link :to="'/app/appStat/'+row.id">
            <el-button
              type="primary"
              size="small"
              style="margin-left: 10px"
            >
              统计
            </el-button>
          </router-link>
          <router-link :to="'/app/appEditUser/'+row.id">
            <el-button
              type="primary"
              size="small"
              style="margin-left: 10px"
            >
              修改
            </el-button>
          </router-link>
          <el-button
            type="primary"
            size="small"
            style="margin-left: 10px"
            @click="handleDelete(row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total > 0"
      :total="total"
      :page.sync="listQuery.pageNo"
      :limit.sync="listQuery.pageSize"
      @pagination="getList"
    />
  </div>
</template>

<script>
import { appList, appDelete } from '@/api/app'
import waves from '@/directive/waves'
import Pagination from '@/components/Pagination'
export default {
  name: 'AppList',
  components: { Pagination },
  directives: { waves },
  data() {
    return {
      tableKey: 'qq',
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        pageNo: 1,
        pageSize: 20,
        info: ''
      },
      parentsProps: { checkStrictly: true, value: 'id', label: 'name', children: 'children' },
      parentsKey: 0,
      temp: {
        'id': undefined,
        'appId': undefined,
        'name': undefined,
        'secret': undefined,
        'white_ips': undefined,
        'qpd': undefined,
        'qps': undefined,
        'created_at': undefined,
        'updated_at': undefined
      },
      dialogFormVisible: false,
      dialogStatus: '',
      rules: {}
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      appList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.page_no = 1
      this.getList()
    },
    handleDelete(row) {
      const tempData = {
        id: row.id
      }
      this.$confirm('请确认是否删除?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        appDelete(tempData).then(() => {
          this.dialogFormVisible = false
          this.$notify({
            title: 'Success',
            message: '成功删除',
            type: 'success',
            duration: 1500
          })
          this.getList()
        })
      }).catch(() => {})
    }
  }
}
</script>
