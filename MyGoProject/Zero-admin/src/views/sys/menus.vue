<template>
  <el-card class="box-card">
    <el-col :span="12" class="toolbar">
      <el-button @click="handleAddLevel1Menu">新增根目录</el-button>
      <br>
      <br>
      <el-tree
        :data="menus"
        :props="defaultProps"
        node-key="id"
        highlight-current
        :expand-on-click-node="true"
        default-expand-all
        accordion
        :render-content="renderContent"
        style="margin-bottom:50px">
      </el-tree>
    </el-col>
    <el-col :span="12">
      <el-form :model="menuForm" ref="menuForm">
        <el-form-item label="节点名称">
          <el-input v-model="menuForm.name" placeholder="请填写节点名称"></el-input>
        </el-form-item>
        <el-form-item label="适用机构" v-if="menuForm.level=='level1' || menuForm.level=='level2'">
          <el-select v-model="menuForm.org_type_ids" placeholder="请选择机构" style="width:100%;" multiple>
            <el-option v-for="item in org_types" :key="item.id" :label="item.name" :value="item.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="页面路由" v-if="menuForm.level=='level2'">
          <el-input v-model="menuForm.route_link" placeholder="请填写节点路由地址"></el-input>
        </el-form-item>
        <el-form-item label="排序编号" v-if="menuForm.level=='level1' || menuForm.level=='level2'">
          <el-input v-model="menuForm.index" placeholder="请填写排序编码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleSave">保存</el-button>
          <el-button type="danger" @click="handleDelete">删除节点</el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-card>
</template>

<script>
  let id = 1000;
  import { Notification } from 'element-ui'
  import {LocalAccount} from '@/api/local-account'
  import {findSysMenus, saveSysMenu, deleteSysMenu, findOrgTypesSelect} from '@/api/api'
  export default {
    data() {
      return {
        menuForm: {
          id: '',
          name: '',
          parent_id: 0,
          node_type: '',
          route_link: '',
          org_type_ids: [],
          index: '',
        },
        menus: [],
        org_types: [],
        node_types: [
          {type_name: '菜单', node_type: 'menu'},
          {type_name: '权限', node_type: 'permission'},
        ],
        defaultProps: {
          children: 'children',
          label: 'name'
        },
        new_insert: {}
      }
    },
    created() {
      this.fetchOrgTypes()
      this.fetchMenus()
    },
    methods: {
      fetchOrgTypes() {
          findOrgTypesSelect().then(result => {
              this.org_types = result.data
          })
      },
      fetchMenus() {
        findSysMenus().then(result=>{
          this.menus = result.data
        })
      },
      handleSave() {
        saveSysMenu(this.menuForm).then(result=>{
          if(result.code === 0) {
            Notification.success({title: '系统提示', message: '修改节点成功!', duration: 2000})
            this.clearMenuForm()
            this.$refs.menuForm.resetFields()
            this.fetchMenus()
          }
        })
      },
      handleDelete() {
        this.$confirm('删除<'+this.menuForm.name+'>节点将会删除<'+this.menuForm.name+'>下的所有子节点, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          deleteSysMenu({menu_id: this.menuForm.id}).then(result=>{
            if(result.code === 0) {
              this.$message({
                type: 'success',
                message: '删除成功!'
              })
              this.$refs.menuForm.resetFields()
              this.fetchMenus()
            }
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消删除'
          })
        })
      },
      handleAddLevel1Menu() {
        let menu = {}
        menu.name = '一级节点'
        menu.parent_id = 0
        menu.children = []
        menu.level = 'level1'
        this.menus.push(menu)
      },
      append(store, data, node, e) {
        if (!node.data.id) {
          Notification.info({title: '系统提示', message: '请保存新增的节点!', duration: 2000})
          return
        }

        this.clearMenuForm()

        let menu = {}
        menu.org_type_ids = node.data.org_type_ids
        menu.parent_id = node.data.id
        menu.children = []
        if (node.data.parent_id === 0) {
          menu.name = '二级节点'
        } else {
          menu.name = '权限名称'
        }

        if (node.data.level == 'level1') {
          menu.level = 'level2'
        }

        this.menuForm = menu
        store.append(menu, data)
        e.stopPropagation()
      },
      update(store, data, e) {
        this.clearMenuForm()

        _.assign(this.menuForm, data)
        this.api_urls = this.menuForm.api_urls || []

        e.stopPropagation()
      },
      delete(store, data, e) {
        _.assign(this.menuForm, data)
        e.stopPropagation()
      },
      clearMenuForm() {
        this.menuForm = {
          id: '',
          name: '',
          parent_id: 0,
          node_type: '',
          route_link: '',
          api_urls: 'null',
          // api_urls: [],
          org_type_ids: [],
          index: ''
        }
        this.api_urls = 'null'
        // this.api_urls = []
      },
      renderContent(h, { node, data, store }) {
        if (node.data.level === 'level1' ) {
          return (
            <span style="margin-top:5px;margin-bottom:5px;">
              <span>
                <span>{node.label}</span>
              </span>
              <span style="margin-left: 20px;">
                <el-button type="text" size="mini" on-click={ (e) => this.append(store, data, node, e) }>新增页面</el-button>
                <el-button type="text" size="mini" on-click={ (e) => this.update(store, data, e) }>修改</el-button>
              </span>
            </span>
          )
        }

        if (node.data.level === 'level2') {
          return (
            <span style="margin-top:5px;margin-bottom:5px;">
              <span>
                <span>{node.label}</span>
              </span>
              <span style="margin-left: 20px;">
                <el-button type="text" size="mini" on-click={ (e) => this.update(store, data, e) }>修改</el-button>
              </span>
            </span>
          )
        }
      }
    }
  }
</script>
