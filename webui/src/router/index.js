import Vue from 'vue'
import Router from 'vue-router'
import store from '../store/index'

import Index from '@/components/Index'
import Console from '@/components/Console'
import Tables from '@/components/table/Tables'
// import Login from '@/components/Login'
// import Register from '@/components/Register'
// import Table from '@/components/table/Index'
// import Basic from '@/components/account/Basic'
import ViewUI from 'view-design';
// const Table = resolve => require(['@/components/table/Index.vue'], resolve)




Vue.use(Router)

const router = new Router({
    mode: 'history',
    routes: [{
            path: '/',
            name: 'index',
            component: Index
        },

        {
            path: '/login',
            name: 'login',
            component: resolve => require(['../components/Login.vue'], resolve)
        },
        {
            path: '/register',
            name: 'register',
            component: resolve => require(['../components/Register.vue'], resolve)
        },

        {
            path: '/console',
            name: 'console',
            component: Console,
            children: [{
                    path: 'tables',
                    name: 'tables',
                    component: Tables,
                    meta: { requireAuth: true },
                },
                {
                    path: 'tables/:tableId',
                    name: 'tablesDetail',
                    component: resolve => require(['../components/table/Index.vue'], resolve)
                },
                {
                    path: 'tables/:tableId/view/:viewId',
                    name: 'tablesView',
                    component: resolve => require(['../components/table/Index.vue'], resolve)
                },
                {
                    path: 'setting/basic',
                    name: 'userInfoSetting',
                    component: resolve => require(['../components/account/Basic'], resolve)
                },
                {
                    path: 'templates',
                    name: 'templates',
                    component: resolve => require(['../components/template/Index'], resolve)
                },
                {
                    path: 'templates/:templateId',
                    name: 'templatesDetail',
                    component: resolve => require(['../components/template/Detail'], resolve)
                },
                {
                    path: 'templates/:templateId/columns',
                    name: 'templatesColumn',
                    component: resolve => require(['../components/template/Column'], resolve)
                },
                {
                    path: 'setting/organizations',
                    name: 'organizations',
                    component: resolve => require(['../components/organizations/Index'], resolve)
                },
                {
                    path: 'setting/organizations/:organizationId',
                    name: 'organizationsBasic',
                    component: resolve => require(['../components/organizations/Basic'], resolve)
                },
                {
                    path: 'setting/organizations/:organizationId/team',
                    name: 'organizationsTeam',
                    component: resolve => require(['../components/organizations/Team'], resolve)
                },
                {
                    path: 'setting/organizations/:organizationId/members',
                    name: 'organizationsMember',
                    component: resolve => require(['../components/organizations/Member'], resolve)
                },
            ]
        },

        // 简单设置404页面
        {
            path: '*',
            component(resolve) {
                require.ensure(['@/components/404.vue'], () => {
                    resolve(require('@/components/404.vue'));
                });
            },
            hidden: true
        }
    ]
});

router.beforeEach((to, from, next) => {
    // iView.LoadingBar.start();
    let token = localStorage.getItem("token");
    if (to.meta.requireAuth) {
        if (token) {
            next()
        } else {
            next({
                path: '/login',
                query: { redirect: to.fullPath }
            })
        }
    } else {
        next()
    }
})

router.afterEach((to, from, next) => {
    // iView.LoadingBar.finish();
});

export default router;