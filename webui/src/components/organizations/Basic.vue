<template>
    <div class="container">
        <organization-header :organization="organization"></organization-header>
        <div class="content">
            <div class="organization" v-if="false">
                <img :src="organization.icon" class="icon" v-if="organization.icon" />
                <img src="../../assets/placeholder.png" class="icon" v-if="!organization.icon"/>
                <div class="data">
                    <div style="padding-bottom:10px;font-size:1.2em;">
                        封面图的网址: 
                        <small style="margin-left:10px;">
                            <a target="_blank" href="/console/setting/images">点此上传图片</a>
                        </small>
                    </div>
                    <div style="padding-bottom:15px;">
                        <Input v-model="organization.icon" />
                    </div>
                    <div>
                        <Button type="primary" @click.native="save(true)">保存封面图</Button>
                    </div>
                </div>
            </div>

            <Form :model="organization"  label-position="top" v-if="organization">
                <FormItem label="团队名称">
                    <Input v-model="organization.name" placeholder=""></Input>
                </FormItem>
                <Row v-if="false">
                    <Col span="12" style="padding:0 10px 0 0;">
                        <FormItem label="官方原价（元）">
                            <Input v-model="organization.original_price" placeholder=""></Input>
                        </FormItem>
                    </Col>
                    <!-- <Col span="6" style="padding:0 10px 0 0;">
                        <FormItem label="可选：原价的补充描述">
                            <Input v-model="organization.oringal_price_tail" placeholder=""></Input>
                        </FormItem>
                        <Alert type="warning">警告提示文案</Alert>
                    </Col> -->
                    <Col span="12">
                        <FormItem label="预约时需支付订金（元）">
                            <Input v-model="organization.deposit_price" placeholder=""></Input>
                        </FormItem>
                    </Col>
                </Row>

                <FormItem label="团队简介">
                    <Input v-model="organization.description" :rows="5" type="textarea" placeholder=""></Input>
                </FormItem>

                <FormItem label="">
                    <Button type="primary" @click.native="save()">保存并生效</Button>
                    <Poptip
                        confirm
                        title="您确定删除此服务吗？"
                        placement="left"
                        class="delete-button"
                        @on-ok="deleteOrganization"
                        >
                        <span class="text-error" >删除本团队</span>
                    </Poptip>
                </FormItem>
            </Form>
            
        </div>
    </div>
</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'
import OrganizationHeader from './OrganizationHeader'
export default {
    components: {
        OrganizationHeader,
    },
    data () {
        return {
        }
    },
    methods: {
        goImages() {
            this.$router.push({'name':'images'});
        },
        save(onlyAvatar) {
            if(!this.organization.name) {
                this.$Message.warning('请输入团队名');
                return;
            }
            console.log(this.plan);
            var data = {name:this.organization.name,description:this.organization.description,}
            if(onlyAvatar) {
                data = {"icon": this.plan.icon}
            }
            api.updateOrganization(this.organization.id, data).then(res => {
                this.$store.dispatch('GetOrganizations');
                this.$Message.success("保存成功");
            }, err => {
                this.$Notice.error({title:'保存出错', desc: err.response.data.message});
            });
        },
        deleteOrganization() {
            api.deleteOrganization(this.organization.id).then(res => {
                this.$Message.success('删除成功');
                this.$router.push({name:'organizations'})
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
     
    },
    computed: {
        organization() {
            // if(!this.organizations || !this.organizations.length) return {};
            return this.organizations.filter(it => it.id == this.$route.params.organizationId)[0];
        },
        ...mapState({
            user: state => state.user,
            users: state => state.users,
            organizations: state => state.organizations,
        })
    },
    watch: {
        user() {
            this.$store.dispatch('GetOrganizations');
        }
    },
    created() {
        if(!this.organizations && this.user && this.user.id){
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
    padding: 30px;
    background: white;
}
.plan {
    min-height: 200px;
    padding: 20px 40px;
    background: white;
    margin-bottom: 40px;
    border-bottom: 1px solid #f8f8f9;
    clear: both;
}
.icon {
    width: 140px;
    height: 140px;
    float: left;
}
.data {
    padding-left: 180px;
}
.delete-button {
    float: right;
    cursor: pointer;
    
}
</style>
