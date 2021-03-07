<template>
<!--
    <Menu mode="vertical" theme="light"  width="auto" class="nav-menu" :active-name="selectedView.id" @on-select="openView">
        <Menu-group title="视图" >
            <Menu-item v-for="view in views" :name="view.id" :key="view.id">
                <Icon type="document-text"></Icon>
                {{view.name}}
            </Menu-item>
            <Menu-item name="add-new">
            </Menu-item>
        </Menu-group>
        
    </Menu>
    -->
    <Tabs @on-click="openView" class="my-tabs">
        <Tab-pane class="my-tab-pane" v-for="view in views" :name="view.id" :key="view.id" :label="view.name" ></Tab-pane>
    </Tabs>
    
</template>
<script>
import {mapState} from 'vuex'
export default {
    methods: {
        openView: function(key) {
            this.$router.push({name: 'tableView', params:{tableId: this.table.id, viewId: key}});
        }
    },
    computed: mapState([
        'table', 'views', 'view'
    ])    
}
</script>

<style scoped>
    .nav-menu {
        /*background: #f5f7f9;*/
        height: 100%;
    }
    .ivu-tabs {
        background:#f2f0f1;
    }
    .ivu-tabs-nav .ivu-tabs-tab-active {
        background: #fff !important;
    }
    
</style>