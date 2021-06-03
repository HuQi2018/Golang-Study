<template>
<el-card class="box-card">
  <section style="background-color:#EAEAEA;">
    <el-row style="padding:5px;">
      <label style="color:red;font-size:18px;">{{role_name}}</label>&nbsp; &nbsp; 权限分配
    </el-row>
    <el-row>
      <el-col :span="24" class="toolbar">
        <el-button @click="handleCheckedAll">全选</el-button>
        <el-button type="danger" @click="handleResetChecked">全不选</el-button>
        <el-button type="success" @click="handleSave">保存</el-button>
        <br>
        <br>
        <el-tree
          :data="permissions"
          :props="defaultProps"
          ref="tree"
          show-checkbox
          node-key="id"
          @check-change="handleCheckChange"
          highlight-current
          :expand-on-click-node="true"
          :default-expand-all="true"
          accordion>
        </el-tree>
        <br>
        <el-button @click="handleCheckedAll">全选</el-button>
        <el-button type="danger" @click="handleResetChecked">全不选</el-button>
        <el-button type="success" @click="handleSave">保存</el-button>
      </el-col>
    </el-row>
  </section>
</el-card>
</template>

<script>
  import { Notification } from 'element-ui'
  import {LocalAccount} from '@/api/local-account'
  import {deepDiffMapper} from '@/utils/deep_diff'
  import {findRolePermission, saveRolePermission} from '@/api/api'
  export default {
    data() {
      return {
        menu: {
          id: '',
          name: '',
          parent_id: 0,
          type: '',
          api_url: '',
          route_link: '',
          perms: ''
        },
        permissions: [],
        checked_keys: [],
        defaultProps: {
          children: 'children',
          label: 'name'
        }
      }
    },
    created() {
      this.role_id = this.$route.query.id
      this.role_name = this.$route.query.name
      this.fetchRolePermission()
    },
    methods: {
      fetchRolePermission() {
        findRolePermission({role_id: this.role_id}).then(result=>{
          this.permissions = result.data
          this.setCheckedKeys()
        })
      },
      setCheckedKeys() {
        this.permissions.forEach(group=>{
            if(group.checked == 1){
                this.checked_keys.push(group.id)
            }
            group.children.forEach(perm=>{
                if(perm.checked == 1){
                     this.checked_keys.push(perm.id)
                 }
            })
        })
        this.$refs.tree.setCheckedKeys(this.checked_keys)
      },
      handleCheckedAll() {
        this.checked_keys = []
        this.permissions.forEach(item=>{
          this.checked_keys.push(item.id)
        })
        this.$refs.tree.setCheckedKeys(this.checked_keys)
      },
      handleResetChecked() {
        this.$refs.tree.setCheckedKeys([])
      },
      handleCheckChange(data, checked, indeterminate) {
        if (data.group_id != '0' ) {
          data.checked = checked ? 1 : 0
        }
      },
      handleSave() {
          var checked_perm = this.$refs.tree.getCheckedKeys()
        let data = {role_id: this.role_id, permission_ids: checked_perm}
        saveRolePermission(data).then(result=>{
          if (result.code == 0) {
            Notification.success({title: '系统提示', message: '权限分配成功！', duration: 2000})
          } 
        })
      }
    }
  }
</script>
<style scoped>

</style>
