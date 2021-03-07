<template>
    <div>
        <Breadcrumb class="breadcrumb" v-if="false">
            <BreadcrumbItem href="/console/organizations">
                <Icon type="ios-people"></Icon> 团队
            </BreadcrumbItem>
            <BreadcrumbItem v-if="organization && organization.id">
                {{organization.name}} 
            </BreadcrumbItem>
        </Breadcrumb>
        <div class="organization" v-if="organization && organization.id">
            <img :src="organization.icon" class="icon" v-if="organization.icon" />
            <img src="../../assets/placeholder.png" class="icon" v-if="!organization.icon"/>
            <div class="data">
                <div class="title">{{organization.name}} </div>

                <div  class="sub-title">
                    成员: {{organization.member_count}}/{{organization.member_limit}}人
                </div>
                <div  class="sub-title">
                    表格: {{organization.table_count}}/{{organization.table_limit}}个
                </div>
                
            </div>
        </div>
        <div class="btn" v-if="organization && organization.id">
            <Menu mode="horizontal" :active-name="activeName" @on-select="menuSelected">
                <MenuItem name="organizationsBasic" v-if="false">
                    <Icon type="ios-information"></Icon>
                    基础设置
                </MenuItem>
                <MenuItem name="organizationsMember">
                    <Icon type="ios-people"></Icon>
                    团队成员
                </MenuItem>
                <MenuItem name="organizationsTeam">
                    <Icon type="ribbon-b"></Icon>
                    团队分组
                </MenuItem>
            </Menu>
        </div>
    </div>
</template>
<script>
import {mapState} from 'vuex'
export default {
    props: ['organization'],
    data () {
        return {
            activeName: this.$route.name
        }
    },
    methods: {
        go(organization, routerName) {
            this.$router.push({name:routerName, params:{organizationId: this.organization.id}});
        },
        menuSelected(routerName) {
            this.$router.push({name:routerName, params:{organizationId: this.organization.id}});
        },
    },
    computed: {
        ...mapState({
            user: state => state.user,

        })
    },

    
}
</script>
<style scoped>
.breadcrumb {
    font-size: 1.4em;
    margin-bottom: 15px;
}
.btn {
    /* margin-bottom: 10px; */
}
.ivu-btn {
    margin-right:10px;
}
.ivu-btn-ghost {
    background-color: white;
}
.organization {
    min-height: 150px;
    padding: 20px 40px;
    background: white;
    margin-bottom: 20px;
}
.organization img {
    vertical-align: middle;
}
.icon {
    width: 100px;
    height: 100px;
    float: left;
}
.data {
    padding-left: 130px;
}
.title {
    color: #000;
    font-size: 1.3em;
}
.sub-title {
    font-size: 1.1em;
    padding: 10px 0 5px 0;
}
</style>