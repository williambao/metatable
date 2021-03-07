<template>
      <div class="container">
        <div class="content" v-if="false">
            <div class="yue">
                <div class="header">
                    {{organizations && organizations.length || 0}}团队
                    <Button class="right" type="success" @click.native="openAdd=true">添加新团队</Button>
                </div>
                <div >
                    <div class="plan-item" v-for="organization in organizations" :key="organization.id">
                        <img :src="organization.icon" class="icon" v-if="organization.icon" />
                        <img src="../../assets/placeholder.png" class="icon" v-if="!organization.icon"/>
                        <div class="data">
                            <div class="item-header">
                                <h3 class="title">{{organization.name}}</h3>
                                <div class="right">
                                    <Button type="primary" size="small" @click.native="goOrganization(organization)">编辑</Button>
                                </div>
                            </div>
                            
                            <Row class="line-item">
                                <Col span="12" class="col">
                                    <span class="label">
                                        <Icon type="ios-people"></Icon>
                                        成员数:</span>
                                    <span>{{organization.member_count}}/{{organization.member_limit}}</span>
                                </Col>
                                <Col span="12"  class="col">
                                    <span class="label">
                                        <Icon type="ios-list-outline"></Icon>
                                        表格数:</span>
                                    <span>{{organization.table_count}}/{{organization.table_limit}}</span>
                                </Col>
                            </Row>
                            <Row class="line-item">
                                <Col span="12"  class="col">
                                    <span class="label">
                                        <Icon type="ribbon-b"></Icon>
                                        分组:</span>
                                    <span>{{organization.teams.length}}</span>
                                </Col>
                            </Row>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    <Modal
        v-model="openAdd"
        title="新增团队"
        :loading="loading"
        @on-ok="save">
        <div>
            <Form :model="newOrgnization" :label-width="140">
                <FormItem label="团队名称">
                    <Input v-model="newOrgnization.name" placeholder="比如: 艾米数据" ></Input>
                </FormItem>
                <FormItem label="团队介绍">
                    <Input v-model="newOrgnization.original_price" placeholder="" type="textarea" :rows="4"></Input>
                </FormItem>
                <!-- <FormItem label="可选: 原价的补充描述">
                    <Input v-model="newPlan.oringal_price_tail" placeholder="比如: 1288. 则预约界面此套餐价格将显示: 888~1288"></Input>
                </FormItem> -->
            </Form>
        </div>
    </Modal>        
    </div>

</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'

let _organization = {name: "", description: ""};
export default {
    data () {
        return {
            openAdd: false,
            loading: false,
            newOrgnization: Object.assign({}, _organization),
        }
    },
    methods: {
        save(){
            api.addOrganization(this.newOrgnization).then(res => {
                this.$Message.success("创建成功")   
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        goOrganization(org) {
            this.$router.push({name: 'organizationsBasic', params: {organizationId: org.id}});
        },
    },
    computed: {
        ...mapState({
            user: state => state.user,
            organizations: state => state.organizations,
        }),
    },
    watch: {
        user() {
            this.$store.dispatch('GetOrganizations');
        }
    },    
    created() {
        if(!this.organizations || this.organizations.length) {
            this.$store.dispatch('GetOrganizations');
        }
    }
}
</script>
<style scoped>
.container {
    padding: 40px;
    /* width: 500px; */
}
.content {
    
}
.right{
    float: right;
}
.yue {
    padding: 30px;
    background: white;
}
.header {
    color: black;
    min-height: 45px;
    font-size: 1.3em;
    padding: 3px 0;
    border-bottom: 3px solid #000;
    margin-bottom: 40px;
}
.sub-header {
    color: gray;
    margin-bottom: 20px;
}
.plans {
    margin-bottom: 15px;
}
.plan {
    float: left;
    margin-right: 40px;
    margin-bottom: 10px;
}
.plan-item {
    width: 100%;
    margin-bottom: 70px;
    min-height: 160px;
}
.icon {
    width: 140px;
    height: 140px;
    float: left;
}
.data {
    padding-left: 180px;
}
.item-header {
    padding-bottom: 8px;
    border-bottom: 1px solid #f8f8f9;
}
.line-item {
    height: 40px;
    line-height: 30px;
    font-size: 1.2em;
}
.col {
    border-bottom: 1px solid #f8f8f9;
    padding: 10px 0;
}
.label {
    display: inline-block;
    width: 120px;
    color:#ccc;
}
.title {
    display: inline-block;
}
.right {
    float: right;
}
</style>
