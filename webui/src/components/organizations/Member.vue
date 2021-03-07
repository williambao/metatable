<template>
    <div class="container">
        <organization-header :organization="organization"></organization-header>
        <div class="content">
            <div class="header" v-if="currentMember && currentMember.is_owner">
                团队成员({{members.length}}) 
                <Button class="right" type="success" @click.native="openAdd=true" >新增成员</Button>
            </div>
            <div>
                <table>
                    <thead>
                        <tr>
                            <th>成员</th>
                            <!-- <th>备注</th> -->
                            <th style="width:200px">分组</th>
                            <th>团队角色</th>
                            <th>加入时间</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="member in members" :key="member.id">
                            <td style="vertical-align:middle">
                                <img v-if="member.user_avatar" :src="member.user_avatar" class="avatar"/>
                                <img v-if="!member.user_avatar" src="../../assets/placeholder.png" class="avatar"/>
                                <span>{{member.nickname}}</span>
                            </td>
                            <!-- <td>
                                <Icon type="edit" v-if="currentMember.is_owner || currentMember.is_admin"></Icon>
                            </td> -->
                            <td>
                                <span v-if="!currentMember.is_owner && !currentMember.is_admin" :class="{'gray':!member.team_id}">
                                    {{member.team_name}}
                                </span>
                                <Select placeholder="未分组" v-model="member.team_id" @on-change="updateMember(member)" v-if="currentMember.is_owner || currentMember.is_admin">
                                    <Option value="">未分组</Option>
                                    <Option v-for="item in organization.teams" :value="item.id" :key="item.id">{{ item.name }}</Option>
                                </Select>
                            </td>
                            <td>
                                <span class="text-success" v-if="member.is_owner">所有者</span>
                                <span class="text-success" v-else-if="member.is_admin">管理员</span>
                                <span class="" v-else="member.is_admin">成员</span>
                            </td>
                            <td>{{member.created_at|moment("YYYY-MM-DD HH:mm")}}</td>
                            <td>
                                <Poptip
                                    confirm
                                    title="您确认删除此成员吗？"
                                    @on-ok="removeMember(member)"
                                    v-if="!member.is_owner && currentMember.is_admin"
                                    >
                                    <Button size="small"  >删除</Button>
                                </Poptip>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="content" style="margin-top:15px;" v-if="false">
            <Button type="error">转让团队</Button>
        </div>
    <Modal
        v-model="openAdd"
        title="新增成员"
        :loading="loading"
        @on-ok="addMember">
        <div>
            <Form :model="newMember" :label-width="140">
                <FormItem label="分组:">
                    <Select v-model="newMember.team_id" placeholder="未分组">
                        <Option value="">未分组</Option>
                        <Option :value="st.id" v-for="st in organization.teams" :key="st.id">{{st.name}}</Option>
                    </Select>
                </FormItem>
                <FormItem label="手机号:">
                    <Input v-model="newMember.phone" placeholder="限中国大陆11位手机号" ></Input>
                </FormItem>
                <FormItem label="名称:">
                    <Input v-model="newMember.nickname" placeholder="比如: 张小丽" ></Input>
                </FormItem>
                <FormItem label="账户Email(用于登陆):">
                    <Input v-model="newMember.email" placeholder="比如: zhang@qq.com"></Input>
                </FormItem>
                <FormItem label="登陆密码:">
                    <Input v-model="newMember.password" placeholder="至少6位以上, 不要用太过简单的密码"></Input>
                </FormItem>
                <!-- <FormItem label="可选: 原价的补充描述">
                    <Input v-model="newMember.oringal_price_tail" placeholder="比如: 1288. 则预约界面此套餐价格将显示: 888~1288"></Input>
                </FormItem> -->
            </Form>
        </div>
    </Modal>        
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
            openAdd: false,
            loading: false,
            members: [],
            currentMember: {},
            newMember: {},
        }
    },
    methods: {
        goImages() {
            this.$router.push({'name':'images'});
        },
        addMember() {
            if(!this.newMember.nickname) {
                this.$Message.warning("请输入名称");
                return;
            }
            if(!this.newMember.email) {
                this.$Message.warning("请输入Email");
                return;
            }
            if(!this.newMember.password) {
                this.$Message.warning("请输入登陆密码");
                return;
            }
            api.addOrganizationUser(this.organization.id, this.newMember).then(res => {
                this.getMembers();
                this.username = "";
                this.$Message.success("增加成功");
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        getMembers() {
            if(!this.organization.id) return;
            api.getOrganizationUsers(this.organization.id).then(res => {
                var current = res.data.filter(it => it.user_id == this.user.id);
                if(current && current.length) {
                    this.currentMember = current[0];
                }

                this.members = res.data;
            });
        },
        updateMember(member) {
            console.log(member.team_id);
            var data = {"team_id": member.team_id};
            api.updateOrganizationUser(this.organization.id, member.id, data).then(res => {
                this.getMembers();
                this.$Message.success("更新成功");
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });

        },
        removeMember(member) {
            api.deleteOrganizationUser(this.organization.id, member.id).then(res => {
                this.getMembers();
                this.$Message.success("删除成功");
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        
    },
    computed: {
        organization() {
            if(!this.organizations || !this.organizations.length) return {};
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
        },
        organizations() {
            this.getMembers();
        }
    },
    created() {
        if(!this.organizations && this.user && this.user.id){
            this.$store.dispatch('GetOrganizations');
        }
        if(this.organizations && !this.members.length) {
            this.getMembers();
        }
    }          
}
</script>
<style scoped>
.container {
    padding: 40px;
    height: 100%;
    overflow-y: auto;
}
.content {
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
.plan {
    min-height: 200px;
    padding: 20px 40px;
    background: white;
    margin-bottom: 40px;
    border-bottom: 1px solid #f8f8f9;
    clear: both;
}
.avatar {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    vertical-align: middle;
    margin-right: 8px;
}
.data {
    padding-left: 180px;
}
.delete-button {
    float: right;
    cursor: pointer;
    
}
table {
    width: 100%;
    border-collapse: collapse;
}
thead tr {
    border-bottom: 1px solid #ccc;
}
th {
    color: #999;
    border-bottom: 1px solid #ddd;
    padding: 0.5em;
    text-align: left;
}
tbody td {
    padding: 8px 15px 8px 5px;
}
tbody tr:hover {
    background-color: #f8f8f9;
}
.gray {
    color: #bbbec4;
}
</style>
