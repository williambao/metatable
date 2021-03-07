<template>
    <div class="container">
        <div span="24" class="content-box" >
            <toolbar v-if="permission" :table="table" v-on:refresh="getTable" v-on:refreshData="getRecords"></toolbar>
            <table-data v-if="permission" :columns="columns" :table="table" v-on:refresh="getTable" v-on:refreshData="getRecords"></table-data>
            <div v-cloak v-if="permission" class="warning-container">
                <h1 >您无权限查看此表格内容</h1>
            </div>
        </div>
    <!--
        <div class="sheet">
            <sheet ></sheet>
        </div>
    -->
    </div>
</template>
<script>
import {mapState} from 'vuex'
import * as types from '../../store/types'
import api from '../../axios'
import Sheet from './Sheet'
import Toolbar from './Toolbar'
import TableData from './TableData'
export default {
    components: {
        Sheet,
        Toolbar,
        TableData,
    },
    data() {
        return {
            searchKey: "",
        }
    },
    methods: {
        getUserTables() {
            let self = this;
            api.getUserTables(self.user.id).then(function(res) {
                console.log(res)
                self.$store.dispatch('SetTables', res.data)
            })
        },
        getTable(searchText) {
            let tableId = this.$route.params.tableId;
            let self = this;
            api.getTable(tableId).then(response => {
                self.$store.commit(types.SELECT_TABLE, response.data);
                self.$store.commit(types.VIEWS, response.data.views);
                self.$store.commit(types.COLUMNS, response.data.columns);
                self.$store.commit(types.TABLE_PERMISSION, response.data.permission);
                if(self.table.organization_id) {
                    self.$store.dispatch('GetOrganizationMembers', response.data.organization_id);
                }
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        getRecords() {
            let self = this;
            if(!self.table.id){
                return;
            }
            self.$store.dispatch('IsLoading', true);

            var query = {};
            query.searchText = this.tableSearchText;
            
            // debugger;
            api.getViewRecords(self.table.id, self.view.id, query).then(res => {
                // 处理下数据
                var data = [];
                for(var i = 0; i <1; i++) {
                    res.data.forEach(function(item) {
                        item.is_selected = false
                        data.push(item);
                    })
                }
                self.$store.dispatch('SetRecords', data);
                self.$store.dispatch('IsLoading', false);
            }, err => {
                self.$store.dispatch('IsLoading', false);
                this.$Message.error(err.response.data.message);
            });
        },
        getConditions() {
            let self = this;
            api.getConditions().then(res => {
                self.$store.commit(types.CONDITIONS, res.data);
            });
        },       
    },
    computed: {
        permission() {
            return this.$store.state.permission;
        },
        ...mapState([
            'user', 'tables', 'table', 'columns', 'view', 'views', 'tableSearchText', 'isLoading'
        ]),
    },
    watch: {
        '$route': 'getTable', 
        user() {
            this.getTable();
            if(!this.tables || !this.tables.length){
                this.getUserTables();
            }
        },
        views(newVal, oldVal) {
            let viewId = this.$route.params.viewId;
            var view = this.views[0];
            this.views.forEach(function(item) {
                if(viewId) {
                    if(item.id == viewId) {
                        view = item;
                    }
                } else {
                    if(item.is_default) {
                        view = item;
                    }
                }
            });
            this.$store.commit(types.TABLE_VIEW, view);
            this.getRecords();
        },
        tableSearchText(oldVal, newVal) {
            this.getRecords();
        }
    },
    created() {
        let token = localStorage.getItem("token");
        this.$store.commit(types.SHOW_NAV_TABLES, true);
        this.getConditions();
        this.getTable();
        if(!this.tables || !this.tables.length && token){
            this.getUserTables();
        }
    }
}    
</script>
<style scoped>
    .container {
        height: 100%;
        min-width:1024px;
    }
    .warning-container {
        text-align: center;
        padding: 80px;
    }    
    .sheet {
        position: fixed;
        bottom: 0px;
        width:100%;
        height: 36px;
        border-top: 1px solid #d7dde4;
    }
    .ivu-tabs-nav .ivu-tabs-tab:hover {
        color: #333;
    }
    .left-nav {
        height: 100%;
    }
    .content-box {
        background:#fff;
        height: 100%;
        padding-bottom: 50px;
    }
    .nav-menu {
        /*background: #f5f7f9;*/
        height: 100%;
    }

    .table-container {
        background: #fff;
        width:100%;
        height:100%;
    }
    .ivu-menu {
        background:#f2f0f1;
    }
    .ivu-menu-item-active {
        background: #fff;
    }
    [v-cloak] { display:none; }
</style>