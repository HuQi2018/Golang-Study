import axios from '../utils/axios'

// 登录
export const loginRequest = params => { return axios.post(`/v1/admin_api/login`, params).then(res => res.data) }

//接口权限管理
export const findPermission = params => { return axios.get("/v1/admin_api/permission/query", { params: params }).then(res => res.data) }
export const savePermission = params => { return axios.post("/v1/admin_api/permission/save", params).then(res => res.data) }
export const delPermission  = params => { return axios.get("/v1/admin_api/permission/del", {params}).then(res => res.data) }

//角色接口权限管理
export const findRolePermission = params => { return axios.get(`/v1/admin_api/role_permission/query`, { params: params }).then(res => res.data) }
export const saveRolePermission = params => { return axios.post(`/v1/admin_api/role_permission/save`, params).then(res => res.data) }

//菜单管理
export const findSysMenus = params => { return axios.get(`/v1/admin_api/sys_menu/query`, { params: params }).then(res => res.data) }
export const saveSysMenu = params => { return axios.post(`/v1/admin_api/sys_menu/save`, params).then(res => res.data) }
export const deleteSysMenu = params => { return axios.get(`/v1/admin_api/sys_menu/del`, { params: params }).then(res => res.data) }

//角色菜单管理
export const findRoleMenus = params => { return axios.get(`/v1/admin_api/role_menu/query`, { params: params }).then(res => res.data) }
export const saveRoleMenus = params => { return axios.post(`/v1/admin_api/role_menu/save`, params).then(res => res.data) }

//角色管理
export const findRoles = params => { return axios.get(`/v1/admin_api/role/query`, { params: params }).then(res => res.data) }
export const saveRole = params => { return axios.post(`/v1/admin_api/role/save`, params).then(res => res.data) }
export const delRole = params => { return axios.get(`/v1/admin_api/role/del`, {params}).then(res => res.data) }

// 用户管理
export const findEmployees = params => { return axios.get(`/v1/admin_api/user/query`, { params: params }).then(res => res.data) }
export const getEmployee = params => { return axios.get(`/v1/admin_api/user/get`, { params: params }).then(res => res.data) }
export const saveEmployee = params => { return axios.post(`/v1/admin_api/user/save`, params).then(res => res.data) }
export const deleteEmployee = params => { return axios.delete(`/v1/admin_api/user/delete`, { params: params }).then(res => res.data) } //此删除接口不存在？
export const delEmployees = params => { return axios.post(`/v1/admin_api/user/dels`, params).then(res => res.data) }
export const updatePasswordEmployee = params => { return axios.post(`/v1/admin_api/user/update_password`, params).then(res => res.data) }
export const resetPasswordEmployee = params => { return axios.post(`/v1/admin_api/user/reset_password`, params).then(res => res.data) }

// 机构类型
export const findOrgTypesTree = params => { return axios.get(`/v1/admin_api/org_type/query_by_tree`, { params: params }).then(res => res.data) }
export const findOrgTypesSelect = params => { return axios.get(`/v1/admin_api/org_type/query_by_select`, { params: params }).then(res => res.data) }
export const saveOrgType = params => { return axios.post(`/v1/admin_api/org_type/save`, params).then(res => res.data) }

//电影管理
export const findMovies = params => { return axios.get(`/v1/admin_api/movie/query`, { params: params }).then(res => res.data) }
export const getMovie = params => { return axios.get(`/v1/admin_api/movie/get`, { params: params }).then(res => res.data) }
export const saveMovie = params => { return axios.post(`/v1/admin_api/movie/save`, params).then(res => res.data) }
export const delMovie = params => { return axios.get(`/v1/admin_api/movie/del`, {params}).then(res => res.data) }

//电影类型管理
export const findMoviesType = params => { return axios.get(`/v1/admin_api/movie_type/query`, { params: params }).then(res => res.data) }
export const saveMovieType = params => { return axios.post(`/v1/admin_api/movie_type/save`, params).then(res => res.data) }
export const delMovieType = params => { return axios.get(`/v1/admin_api/movie_type/del`, {params}).then(res => res.data) }


//电影标签管理
export const findMoviesTag = params => { return axios.get(`/v1/admin_api/movie_tag/query`, { params: params }).then(res => res.data) }
export const saveMovieTag = params => { return axios.post(`/v1/admin_api/movie_tag/save`, params).then(res => res.data) }
export const delMovieTag = params => { return axios.get(`/v1/admin_api/movie_tag/del`, {params}).then(res => res.data) }


