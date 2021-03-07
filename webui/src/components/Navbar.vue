<template>
    <section id="navbar" class="layout" :class="{ slideInDown: show, slideOutDown: !show }">
        <div class="layout-logo">
            <a href="/console/tables" class="wrapper-header-nav-logo router-link-active">
                <img src="../assets/logo.png" style="vertical-align:middle;width:20px;border-radius:2px;">

                <strong>{{$config.APP_NAME}}</strong>
                <span style="font-size:10px;">beta</span>
            </a>
        </div>    
        <div class="layout-spin" v-if="isLoading">
            <Spin>
                <Icon type="load-c" size=18 class="demo-spin-icon-load"></Icon>
            </Spin>
        </div>
        <div class="layout-tables" v-if="tables && tables.length && showNavTables">
            <div class="layout-ceiling-main">
                <Dropdown @on-click="selectTable">
                    <a href="javascript:void(0)">
                        <span v-if="!table">选择表格</span>
                        <span v-if="table" v-text="table.name"></span>
                        <Icon type="arrow-down-b"></Icon>
                    </a>
                    <Dropdown-menu slot="list" >
                        <Dropdown-item  v-for="(table, index) in tables" @click="selectTable(table)" :name="index" :key="table.id">{{table.name}}</Dropdown-item>
                    </Dropdown-menu>
                </Dropdown>
            </div>
        </div>

        <div class="layout-ceiling" >
            <div class="item layout-ceiling-main" v-if="!user || !user.id">
                <router-link to="/login">登录</router-link>
            </div>
            <div class="layout-ceiling-main" style="padding-top:2px;" v-if="user && user.id">
                <Poptip trigger="click" placement="bottom-end" v-if="user.avatar">
                    <img class="avatar" :src="user.avatar + '?imageView2/1/w/100'" alt="头像">
                    <div slot="content">
                        <img :src="user.avatar + '?imageView2/1/w/500'" alt="" style="width:200px;height:200px;">
                    </div>
                </Poptip>
                <img v-if="!user.avatar" class="avatar" src="../assets/placeholder.png" alt="默认头像" title="默认头像">
                <Dropdown placement="bottom-end" @on-click="userInfoDropdown">
                    <a href="javascript:void(0)">
                        {{user.nickname}}
                        <Icon type="arrow-down-b"></Icon>
                    </a>
                    <Dropdown-menu slot="list">
                        <Dropdown-item name="info">个人信息</Dropdown-item>
                        <!--<Dropdown-item name="contactus" divided>联系我们</Dropdown-item>-->
                        <Dropdown-item name="logout" divided>登出</Dropdown-item>
                    </Dropdown-menu>
                </Dropdown>
            </div>
        </div>
    </section>
</template>
<script>
    import {mapState} from 'vuex'
    export default {
        name: "navbar",
        props: {
            show: Boolean
        },
        methods: {
            selectTable(index) {
                let table = this.tables[index];
                this.$store.dispatch('SetTable', table);
                this.$router.push({name: 'table', params:{tableId: table.id}});
            },
            userInfoDropdown(name) {
                if(name=="logout") {
                    this.$store.dispatch('UserLogout');
                    this.$router.push({name: 'index'});
                } else if(name=="info"){
                    this.$router.push({name:'userInfoSetting'})
                }
            }
        },
        computed: mapState({
            user: state => state.user,
            tables: state => state.tables,
            table: state => state.table,
            isLoading: state => state.isLoading,
            showNavTables: state => state.showNavTables,
        }),

    }
</script>
<style scoped>
    a {
        color: #fff;
    }
    a:hover {
        color: #aaa;
    }
    .item {
        padding: 5px 0;
    }
    .layout {
        position: relative;
        min-height: 40px;
        border-bottom:1px solid #111;
        box-shadow: 0 0 1px rgba(0,0,0,0.25);
        background:#555;
        color: #fff;
    }
    .layout:after {
        content: '';
        display: block;
        clear: both;
    }    
    .layout-spin {
        float: left;
        margin-left: 10px;
        margin-top: 10px;
        
    }
    .layout-spin i {
        color: white !important;
    }
    .layout-tables {
        position: absolute;
        left: 50%;
        margin-left: 20px;
        margin-top: 10px;
    }
    .layout-logo {
        float:left;
        margin-left: 20px;
        margin-top: 10px;
    }
    .layout-logo img {
    }
    .layout-ceiling{
        padding: 5px 0;
        overflow: hidden;
    }
    .layout-ceiling-main{
        float: right;
        margin-right: 25px;
    }
    .avatar {
        width:25px;
        height:25px;
        cursor: pointer;
        border: 1px solid #fff;
        border-radius: 15px;
        display:inline-block;
        vertical-align:middle;
        margin-right:5px;
    }
</style>