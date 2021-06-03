<template>
  <el-card class="box-card">
    <el-form :inline="true" :model="search">
      <el-form-item>
        <el-input
          v-model="search.keyword"
          placeholder="类型名称"
          clearable
          style="width:260px;"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="fetchMovieTypes" icon="el-icon-search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleAdd" type="primary" icon="el-icon-plus">新增</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tags" border>
      <el-table-column type="index" align="center"></el-table-column>
			<el-table-column prop="id" label="类型编号" align="center"></el-table-column>
			<el-table-column prop="name" label="类型名称" align="center"></el-table-column>
			<el-table-column prop="op_user" label="操作人" align="center"></el-table-column>
      <el-table-column prop="updated_at" label="更新时间" align="center" :formatter="dateFormatter"></el-table-column>
      <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <el-button @click="handleEdit(scope.row)" type="text">编辑</el-button>
          <!-- <el-button type="text" @click="handleEdit(scope.$index, scope.row)">编辑</el-button> -->
          <el-button @click="handleDel(scope.row)" type="text" style="color:red">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      style="margin-top:10px"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page.sync="search.page"
      :page-size="search.limit"
      :page-sizes="page.sizes"
      :total="page.total"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>
    <el-dialog
      :title="dialog.title"
      :visible.sync="dialog.show"
      width="40%"
      @close="closeDialog">
      <el-form label-width="100px" :model="form" :rules="rules" ref="form">
				<el-form-item label="类型名称" prop="name">
					<el-input v-model="form.name" placeholder="请输入类型名称" maxlength="80" class="form-item"></el-input>
				</el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialog.show = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">保 存</el-button>
      </span>
    </el-dialog>
  </el-card>
</template>
<script>
import {
  // findMovieTypesSelect,
  // findMovieTypesSelect,
  findMoviesType,
  saveMovieType,
  delMovieType
} from "@/api/api"
import { LocalAccount } from "@/api/local-account"
import _ from "lodash"
import moment from "moment"

export default {
  data() {
    return {
      tags: [],
      page: {
        sizes: [10, 20, 30, 50],
        total: 0
      },
      search: {
        page: 1,
        limit: 10
      },
      form: {},
      selectVal:0,
      dialog: {
        show: false,
        title: ""
      },
      rules: {
		    name: [{ required: true, message: "类型名称不能为空", trigger: "blur" }],
      }
    }
  },
  created() {
    // this.fetchMovieTypes()
    this.fetchMovieTypes()
    // this.fetchMovies()
  },
  methods: {
    handleSizeChange(val) {
      this.search.limit = val
      this.fetchMovieTypes()
    },
    handleCurrentChange(val) {
      this.fetchMovieTypes()
    },
    dateFormatter(row, col) {
      return moment(row[col.property]).format("YYYY-MM-DD HH:mm:ss")
    },
    fetchMovieTypes() {
        // this.form.search = this.search
        findMoviesType(this.search).then(result => {
          this.tags = result.data
          this.page.total = result.total
      })
    },
    closeDialog() {
      this.$refs.form.clearValidate()
    },
    handleAdd() {
      this.form = {}
      this.dialog.show = true
      this.dialog.title = "新增"
    },
    handleEdit(item) {
      this.form = _.cloneDeep(item)
      this.dialog.show = true
      this.dialog.title = "编辑"
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
            let user_account = LocalAccount.getUserInfo().account
            if (!user_account) return
            this.form.op_user = user_account
          saveMovieType(this.form).then(res => {
            if (res.code == 0) {
              this.$message.success("已保存")
              this.fetchMovieTypes()
              this.dialog.show = false
            }
          })
        } else {
          console.log("error submit!!")
          return false
        }
      })
    },
    handleDel(row) {
      this.$confirm("您确定要删除"+row.name+"吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(() => {
        delMovieType({id:row.id}).then(result => {
          if (result.code == 0) {
            this.fetchMovieTypes()
            this.$message.success("已删除")
          }
        })
      }).catch(() => {
        //取消
        console.log("n")
      })
    },
  }
}
</script>
<style scoped>
.pagination {
  padding: 10px;
  float: right;
}
.form-item {
  width: 60%;
}
</style>
