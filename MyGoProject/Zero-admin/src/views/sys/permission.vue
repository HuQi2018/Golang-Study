<template>
    <el-card class="box-card">
        <!--工具条-->
        <el-col :span="24" class="toolbar">
            <el-form :inline="true" :model="search" class="demo-form-inline">
                <el-form-item>
                    <el-input v-model="search.key" placeholder="接口组名称" clearable></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button @click="fetchPermission" icon="el-icon-search" >查询</el-button>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleAdd(1)" icon="el-icon-plus" type="primary">新增</el-button>
                </el-form-item>
            </el-form>
        </el-col>
        <!--列表-->
        <template>
        <el-table
          :data="tableData"          
          row-key="id"
          border
          default-expand-all
          :tree-props="{children: 'children', hasChildren: 'hasChildren'}">
          <el-table-column prop="name" label="名称"></el-table-column>
          <el-table-column prop="url" label="接口地址"></el-table-column>
          <el-table-column prop="org_types" label="所属机构" :formatter="orgTypeFormatter"></el-table-column>
          <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleAdd(0, scope.row)" v-if="scope.row.group_id == 0">新增接口</el-button>
          <el-button type="text" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="text" @click="handleDelete(scope.row)">删除</el-button> 
        </template>
      </el-table-column>
        </el-table>
        </template>
        <!--分页-->
        <el-pagination 
            class="pull-right clearfix"
            @size-change="handleSizeChange" style="margin-top:10px"
            @current-change="handleCurrentChange"
            :current-page.sync="search.page"
            :page-size="search.limit"
            :page-sizes="page.sizes"
            :total="page.total"
            layout="total, sizes, prev, pager, next, jumper">
        </el-pagination>
        <el-dialog
          :title="(editType=='new'?'新增':'更新') + (editLevel == 'group'?'接口':'接口组')"
          :visible.sync="dialogVisible"
          size="tiny"
          width="40%">
          <el-form :model="permissionForm" :rules="rules" ref="permissionForm" label-width="80px" label-position="left">
              <el-form-item prop="name" label="名称" style="margin:30px 30px 0 30px">
                  <el-input v-model="permissionForm.name" placeholder="请输入名称"></el-input>
              </el-form-item>
              <el-form-item prop="org_types" label="所属机构" style="margin:15px 30px ">
                  <el-select v-model="permissionForm.org_types" placeholder="请选择所属机构" style="width:100%;" multiple>
                      <el-option v-for="item in orgTypes" :key="item.id" :label="item.name" :value="item.id"></el-option>
                  </el-select>
              </el-form-item>
              <el-form-item prop="url" label="接口地址" style="margin:15px 30px 15px 30px" v-if="editLevel == 'group'">
                  <el-input v-model="permissionForm.url" placeholder="请输入接口地址"></el-input>
              </el-form-item>
            </el-form>
          <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSave">确 定</el-button>
      </span>
        </el-dialog>
    </el-card>
</template>

<script>
    import {LocalAccount} from '@/api/local-account'
    import {findPermission, savePermission, delPermission , findOrgTypesSelect} from '@/api/api'
    import moment from 'moment'
    import _ from 'lodash'
import org_typesVue from './org_types.vue'

    export default {
        data() {
            return {
               editType:'new',
               editLevel:"group",
                formInline: {},
                roles: [],
                orgTypes: [],
                permissionForm: {
                    id: 0,
                    org_types: [],
                    org_type_name: '',
                    name: ''
                },
                dialogVisible: false,
                search: {
                    page: 1,
                    limit: 50,
                },
                page: {
                    sizes: [50, 100, 200, 500],
                    total: 0
                },
                rules: {
                   name: [{required: true, message: '名称不能为空', trigger: 'blur'}],
                   org_types: [{required: true, message: '机构类型不能为空', trigger: 'blur'}],
                   url: [{required: true, message: '接口地址不能为空', trigger: 'blur'}],
                },
                tableData: [],
            }
        },
        created() {
            this.fetchPermission()
            this.fetchOrgTypes()
        },
        methods: {          
            fetchOrgTypes() {
                findOrgTypesSelect().then(result => {
                    this.orgTypes = result.data
                })
            },
            fetchPermission() {
                findPermission(this.search).then(result => {
                    this.tableData = result.data
                    this.page.total = result.total
                })
            },
            handleAdd(flag, row) {
              this.permissionForm = {}
              if(flag == 0){ //新增接口              
                this.permissionForm.group_id = row.id
                this.editLevel = 'group' 
              }else if (flag == 1) { //新增接口组
                this.editLevel = 'api'
              }              
              this.editType = 'new'                         
              this.dialogVisible = true
            },
            handleEdit(row) {
              this.permissionForm = {
                    id: row.id,
                    group_id: row.group_id,
                    org_types: [],
                    org_type_name: '',
                    name: row.name,
                    url:row.url
                }
                var org_types = row.org_types.split(",")
                for (var i = 0; i < org_types.length; i++){
                  this.permissionForm.org_types.push(parseInt(org_types[i]))
                }
                this.editType = 'update'
                if (row.group_id == 0){
                    this.editLevel = 'api'
                }else{
                    this.editLevel = 'group'
                }
                this.dialogVisible = true
            },
            disabledEdit(role) {
                return role.buildin === 1
            },
            handleEditPerms(row) {
                this.$router.push({path: '/sys/role_menu', query: {id: row.id, name: row.name}})
            },
            handleSave: function () {
                this.$refs.permissionForm.validate((valid) => {
                  if (valid) {
                      let org_type = _.find(this.orgTypes, {id: this.permissionForm.org_types})
                      if (org_type) {
                          this.permissionForm.org_type_name = org_type.name
                      }

                      savePermission(this.permissionForm).then(result => {
                          if (result.code == 0) {
                              this.$message.success((this.permissionForm.id ? '编辑' : '新增') + '成功！')
                              this.fetchPermission()
                              this.dialogVisible = false
                          } else {
                          }
                      })
                  }
                })
            },
            orgTypeFormatter(row, column) {
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
            handleSizeChange(val) {
                this.search.limit = val
                this.fetchPermission()
            },
            handleCurrentChange(val) {
                this.fetchPermission()
            },
            handleDelete(row){
                this.$confirm("您确定要删除 "+row.name+" 吗?", "提示", {
                    confirmButtonText: "确定",
                    cancelButtonText: "取消",
                    type: "warning"
                }).then(() => {
                    if(row.group_id == 0){
                        this.$confirm("注意！此操作将删除该组下所有接口,确认删除接口组？", "确认删除", {
                           confirmButtonText: "确定",
                           cancelButtonText: "取消",
                           type: "warning"
                       }).then(() => {                                                 
                            delPermission({id:row.id}).then(res=>{
                               if(res.code == 0){
                                   this.$message.success('已删除')
                                   this.fetchPermission()
                               }
                            }
                           )
                       })
                       .catch(() => {
                           //取消
                           console.log("n")
                       })
                    }else{
                      delPermission({id:row.id}).then(res=>{
                          if(res.code == 0){
                              this.$message.success('已删除')
                              this.fetchPermission()
                          }
                      })
                    }
                })
                .catch(() => {
                    //取消
                    console.log("n")
                })
            }
        }
    }
</script>

<style scoped>
    .toolbar .el-form-item {
        margin-bottom: 10px
    }
</style>
