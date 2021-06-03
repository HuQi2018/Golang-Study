<template>
  <el-card class="box-card">
    <el-form :inline="true" :model="search">
      <el-form-item>
        <el-input
          v-model="search.keyword"
          placeholder="电影名称/电影ID/电影类型/电影标签/电影年份"
          clearable
          style="width:260px;"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="fetchMovies" icon="el-icon-search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="handleAdd" type="primary" icon="el-icon-plus">新增</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="movies" border>
      <el-table-column type="index" align="center"></el-table-column>
			<el-table-column prop="movie_id" label="电影ID" align="center"></el-table-column>
			<el-table-column prop="name" label="电影名称" align="center"></el-table-column>
			<el-table-column prop="types" label="电影类型" align="center"></el-table-column>
			<el-table-column prop="rating" label="电影评分" align="center"></el-table-column>
			<el-table-column prop="year" label="电影年份" align="center"></el-table-column>
      <el-table-column prop="durations" label="电影时长" align="center"></el-table-column>
      <el-table-column prop="aka" label="电影又名" align="center"></el-table-column>
      <el-table-column prop="tags" label="电影标签" align="center"></el-table-column>
      <el-table-column prop="original_title" label="原始标题" align="center"></el-table-column>
      <el-table-column prop="language" label="电影语言" align="center"></el-table-column>
      <el-table-column prop="country" label="制片国家" align="center"></el-table-column>
      <el-table-column prop="op_user" label="操作人" align="center"></el-table-column>
      <el-table-column prop="updated_at" label="更新时间" align="center" :formatter="dateFormatter"></el-table-column>
      <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <!-- <el-button @click="handleEdit(scope.row)" type="text">编辑</el-button> -->
          <!-- <el-button type="text" :disabled="scope.row.role_id === 1" @click="handleEdit(scope.$index, scope.row)">编辑</el-button> -->
          <!-- <el-button @click="handleDel(scope.row)" type="text" style="color:red">删除</el-button> -->
          <el-button type="text" :disabled="scope.row.role_id === 1" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button type="text" @click="handleDel(scope.row)">删除</el-button> 
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
  </el-card>
</template>
<script>
import {
  // findMovieTypesSelect,
  // findMovieTagsSelect,
  findMovies,
  getEmployee,
  saveMovie,
  delMovie
} from "@/api/api"
import { LocalAccount } from "@/api/local-account"
import _ from "lodash"
import moment from "moment"

export default {
  data() {
    return {
      formInline: {},
      movies: [],
      list:[],
      search: {
        page: 1,
        limit: 10,
      },
      page: {
        sizes: [10, 20, 30, 50],
        total: 0
      },
      qrCode: "",
      show_dialog: false,
    }
  },
  created() {
    // this.fetchMovieTypes()
    // this.fetchMovieTags()
    this.fetchMovies()
  },
  methods: {
    dateFormatter(row, column) {
        return moment(row.created_at).format('YYYY-MM-DD HH:MM')
    },
    handleSizeChange(val) {
      this.search.limit = val
      this.fetchMovies()
    },
    handleCurrentChange(val) {
      this.fetchMovies()
    },
    handleEdit(index, item) {
      this.$router.push({
        path: "/sys/movie_edit",
        query: {
          id: item.id
        }
      })
    },
    fetchMovies() {
      this.list.length = 0
      findMovies(this.search).then(result => {
        // var str1 = JSON.parse(localStorage.getItem("webadmin_account"))
        // var my_role_id = str1.user.role_id

        //筛选出role_id大于等于自己的记录，即平级和下级的记录
        this.movies = result.data
        for (var i = 0; i < result.data.length; i++){
          // if(my_role_id <= this.movies[i].role_id){
            this.list.push(this.movies[i])
          // }
        }
        this.movies = this.list
        this.page.total = result.total
        
      })
    },
    handleAdd: function () {
      this.$router.push({
        path: "/sys/movie_edit"
      })
    },
    handleDel(item) {
      this.$confirm("您确定要删除"+item.name+"吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          delMovies([item.id]).then(result => {
            if (result.code == 0) {
              this.fetchMovies()
              this.$message({
                type: "success",
                message: "删除成功!"
              })
            } else {
              this.$message({
                type: "error",
                message: "删除失败!"
              })
            }
          })
        })
        .catch(() => {
          //取消
        })
    },
    handleQrCode(row) {

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
