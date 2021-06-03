<template>
<el-card class="box-card">
  <el-row>
    <el-col :span="10">
      <el-form
        :model="movieForm"
        label-width="90px"
        :rules="editFormRules"
        ref="movieForm"
        class="form"
      >
				<el-form-item label="电影编号" prop="movie_id">
					<el-input v-model="movieForm.movie_id" placeholder="请输入电影编号" maxlength="50" class="form-item" :disabled="isUpdate"></el-input>
				</el-form-item>
				<el-form-item label="电影名称" prop="name">
					<el-input v-model="movieForm.name" placeholder="请输入电影名称" maxlength="80" class="form-item" :disabled="isUpdate"></el-input>
				</el-form-item>
				<el-form-item label="电影评分" prop="rating">
					<el-input v-model="movieForm.rating" placeholder="请输入电影评分" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影时长" prop="durations">
					<el-input v-model="movieForm.durations" placeholder="请输入电影时长" maxlength="80" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影上映年份" prop="year">
					<el-input v-model="movieForm.year" placeholder="请输入电影上映年份" maxlength="20" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影上映日期数据" prop="pubdates">
					<el-input v-model="movieForm.pubdates" placeholder="请输入电影上映日期数据" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影上映日期" prop="pubdate">
					<el-input v-model="movieForm.pubdate" placeholder="请输入电影上映日期" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影又名" prop="aka">
					<el-input v-model="movieForm.aka" placeholder="请输入电影又名" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影原始标题" prop="original_title">
					<el-input v-model="movieForm.original_title" placeholder="请输入电影原始标题" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影语言" prop="language">
					<el-input v-model="movieForm.language" placeholder="请输入电影语言" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="制片国家" prop="country">
					<el-input v-model="movieForm.country" placeholder="请输入制片国家" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="演员信息" prop="actors">
					<el-input v-model="movieForm.actors" placeholder="请输入演员信息" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="作者信息" prop="writers">
					<el-input v-model="movieForm.writers" placeholder="请输入作者信息" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="导演信息" prop="directors">
					<el-input v-model="movieForm.directors" placeholder="请输入导演信息" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="电影简介" prop="summary">
					<el-input v-model="movieForm.summary" placeholder="请输入电影简介" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="图像数据" prop="photos">
					<el-input v-model="movieForm.photos" placeholder="请输入图像数据" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="海报数据" prop="images">
					<el-input v-model="movieForm.images" placeholder="请输入海报数据" class="form-item"></el-input>
				</el-form-item>
				<el-form-item label="视频数据" prop="videos">
					<el-input v-model="movieForm.videos" placeholder="请输入视频数据" class="form-item"></el-input>
				</el-form-item>
        
        <el-form-item prop="types" label="电影类型">
            <el-select v-model="movieForm.types" placeholder="请选择电影类型" style="width:100%;" multiple>
                <el-option v-for="item in movie_types" :key="item.id" :label="item.name" :value="item.id"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item prop="tags" label="电影标签">
            <el-select v-model="movieForm.tags" placeholder="请选择电影标签" style="width:100%;" multiple>
                <el-option v-for="item in movie_tags" :key="item.id" :label="item.name" :value="item.id"></el-option>
            </el-select>
        </el-form-item>
        <!-- <el-table-column prop="movie_types" label="电影类型" :formatter="orgTypeFormatter"></el-table-column>
        <el-table-column prop="movie_tags" label="电影标签" :formatter="orgTypeFormatter"></el-table-column> -->
        <!-- <el-form-item label="电影类型" prop="types">
          <el-select v-model="movieForm.types" style="width:100%">
            <el-option v-for="item in movie_tags" :key="item.id" :label="item.name" :value="item.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="电影标签" prop="tags">
          <el-select v-model="movieForm.tags" style="width:100%">
            <el-option v-for="item in movie_types" :key="item.id" :label="item.name" :value="item.id"></el-option>
          </el-select>
        </el-form-item> -->
        <el-form-item>
          <el-button @click="submit" type="primary">确定</el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</el-card>
</template>
<script>
import { Notification } from "element-ui"
import { LocalAccount } from "@/api/local-account"
import { saveMovie, getMovie, findMoviesType, findMoviesTag } from "@/api/api"
import moment from 'moment'
import _ from 'lodash'

export default {
  data() {
    return {
      isUpdate: false,
      disabled:false,
      roles: [],
      movieForm: {
          id: 0,
          types: [],
          type_name: '--请选择电影类型--',
          tags: [],
          tag_name: '--请选择电影标签--',
          name: ''
      },
      movie_types: [],
      movie_tags: [],
      editFormRules: {
        name: [{ required: true, message: "电影名称不能为空", trigger: "blur" }],
        movie_id: [{required: true, message: "电影ID不能为空", trigger: "blur"}]
      },
      loading: false,
      isBuildin:false,
      // movie: getMovie(),
      movie: [],
    }
  },
  created() {
    this.fetchTags()
    this.fetchTypes()
    this.disabled = false
    if (this.$route.query.id) {
      this.isUpdate = true
      getMovie({ id: this.$route.query.id }).then(result => {
        if (result.code === 0) {
          this.movieForm = result.data
          // var tags = result.data.tags.split(",")
          // console.log(tags)
          // this.movieForm.tags = []
          // for (var i = 0; i < this.movie_tags.length; i++){
          //   if(tags.indexOf(this.movie_tags[i].id)){
          //     this.movieForm.tags.push(this.movie_tags[i].name)
          //   }
          // }
          // console.log(this.movieForm.tags)
          // var types = result.data.types.split(",")
          // console.log(types)
          // this.movieForm.types = []
          // for (var i = 0; i < this.movie_types.length; i++){
          //   if(types.indexOf(this.movie_types[i].id)){
          //     this.movieForm.types.push(this.movie_types[i].name)
          //   }
          // }
          // console.log(this.movieForm.types)
        }
      })
    }
  },
  methods: {
    submit() {
      let user_id = LocalAccount.getUserInfo().account
      if (!user_id) return

      this.movieForm.op_user = user_id
      this.$refs.movieForm.validate(valid => {
        if (valid) {
          // this.movieForm.tags = this.movieForm.tags.join(",")
          // this.movieForm.types = this.movieForm.types.join(",")
          saveMovie(this.movieForm).then(result => {
            if (result.code == 0) {
              Notification.success({
                title: "系统提示",
                message: (this.isUpdate ? "修改" : "新增") + "成功！",
                duration: 2000
              })
              this.$router.push({ path: "/sys/movies" })
            }
          })
        }
      })
    },
    fetchTags() {
      findMoviesTag().then(result => {
        if(result.code == 0) {
          this.movie_tags = result.data
        }
      })
    },
    fetchTypes() {
      findMoviesType().then(result => {
        if(result.code == 0) {
          this.movie_types = result.data
        }
      })
    },
    multiTypeFormatter(row, column) {
        var orgTypeList = new Array(); 
        var OrgTypeStr = ""
        orgTypeList = row.org_types.split(","); 
                        
        for (var i = 0; i < orgTypeList.length; i++ ){
            for(var k = 0; k < this.orgTypes.length; k++){
              if (orgTypeList[i] == this.orgTypes[k].id) {
                OrgTypeStr += this.orgTypes[k].name
                break;
              }
            }                    
            if (i != orgTypeList.length - 1){
              OrgTypeStr += ", "
            }
        }
        return OrgTypeStr
    },
  }
}
</script>
<style scoped>
.form {
  padding-top: 20px
}
</style>
