import Layout from '@/views/layout/Layout'
export const constantRouterMap = [
    { path: '/login', component: () => import('@/views/login/index'), hidden: true },
    { path: '/authredirect', component: () => import('@/views/login/authredirect'), hidden: true },
    {
        path: '',
        component: Layout,
        redirect: 'dashboard',
        children: [{
            path: 'dashboard',
            component: () => import('@/views/dashboard/index'),
            name: 'dashboard',
            meta: { title: '首页', icon: 'dashboard', noCache: true }
        }]
    }
]

export const asyncRouterMap = [
    {
        path: '/sys',
        component: Layout,
        redirect: 'noredirect',
        name: 'sys',
        meta: {
            title: '系统管理',
            icon: 'xitong'
        },
        children: [
            { path: 'menus', component: () => import('@/views/sys/menus'), name: 'menus', meta: { title: '菜单管理', noCache: false } },
            { path: 'permission', component: () => import('@/views/sys/permission'), name: 'permission', meta: { title: '权限管理', noCache: false } },
            { path: 'org_types', component: () => import('@/views/sys/org_types'), name: 'org_types', meta: { title: '机构类型', noCache: true } },
            { path: 'roles', component: () => import('@/views/sys/roles'), name: 'roles', meta: { title: '角色管理', noCache: true } },
            { path: 'role_menu', component: () => import('@/views/sys/role_menu'), name: 'menus', meta: { title: '权限设置', noCache: false }, hidden: true },
            { path: 'role_permission', component: () => import('@/views/sys/role_permission'), name: 'role_permission', meta: { title: '权限设置', noCache: false }, hidden: true },
            { path: 'users', component: () => import('@/views/sys/users'), name: 'users', meta: { title: '用户管理', noCache: true } },
            { path: 'user_edit', component: () => import('@/views/sys/user_edit'), name: 'user_edit', meta: { title: '用户编辑', noCache: true }, hidden: true },
            { path: 'user_update_password', component: () => import('@/views/sys/user_update_password'), name: 'user_update_password', meta: { title: '修改密码', noCache: true }, hidden: true },
        ]
    },{
        path: '/sys',
        component: Layout,
        redirect: 'noredirect',
        name: 'movie',
        meta: {
            title: '电影管理',
            icon: 'component'
        },
        children: [
            { path: 'movies', component: () => import('@/views/sys/movies'), name: 'movies', meta: { title: '电影管理', noCache: true } },
            { path: 'movie_edit', component: () => import('@/views/sys/movie_edit'), name: 'movie_edit', meta: { title: '电影编辑', noCache: true }, hidden: true },
            { path: 'movie_type', component: () => import('@/views/sys/movie_type'), name: 'movie_type', meta: { title: '电影类型管理', noCache: true } },
            { path: 'movie_tag', component: () => import('@/views/sys/movie_tag'), name: 'movie_tag', meta: { title: '电影标签管理', noCache: true } },
        ]
    }
]

export const routes = constantRouterMap