<template>
<div>

        <Breadcrumb class="breadcrumb">
            <BreadcrumbItem href="/console/templates">
                <Icon type="ios-home"></Icon> 模板列表
            </BreadcrumbItem>
            <BreadcrumbItem>
                {{template.name}}
            </BreadcrumbItem>
        </Breadcrumb>
  <div  class="btn">
        <Menu mode="horizontal" :active-name="activeName" @on-select="menuSelected">
            <MenuItem name="templatesDetail">
                <Icon type="information-circled"></Icon>
                基础设置
            </MenuItem>
            <MenuItem name="templatesColumn">
                <Icon type="film-marker"></Icon>
                列/字段
            </MenuItem>
            
        </Menu> 
    </div>
</div>
</template>
<script>
import {mapState} from 'vuex'
export default {
    props: ['template'],
    data () {
        return {
            activeName: this.$route.name
        }
    },
    methods: {
        go(store, routerName) {
            this.$router.push({name:routerName, params:{templateId: this.template.id}});
        },
        menuSelected(routerName) {
            this.$router.push({name:routerName, params:{templateId: this.template.id}});
        },
    },
    computed: {

        ...mapState({
            user: state => state.user,
            users: state => state.users,
        })
    }      
}
</script>
<style scoped>
.breadcrumb {
    font-size: 1.4em;
    margin-bottom: 20px;
}
</style>
