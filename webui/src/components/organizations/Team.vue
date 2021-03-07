<template>
    <div class="container">
        <organization-header :organization="organization"></organization-header>
        <div class="content" v-if="organization">
            <div class="header" v-if="currentMember && currentMember.is_owner">
                团队分组({{organization.teams.length}}) 
                <Button class="right" type="success" @click.native="openAdd=true">新增分组</Button>
            </div>
            <div>
                <table>
                    <thead>
                        <tr>
                            <th>名称</th>
                            <!-- <th>备注</th> -->
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(team, index) in organization.teams" :key="team.id">
                            <td style="vertical-align:middle">
                                <span class="team-name">{{team.name}}</span>
                            </td>
                            <td style="text-align:right"> 
                                <Poptip
                                    confirm
                                    title="您确认删除此成员吗？"
                                    style="text-align:left;"
                                    @on-ok="removeTeam(index)"
                                    v-if="currentMember.is_admin"
                                    >
                                <Button size="small" >删除</Button>
                                </Poptip>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <Alert type="warning" style="margin-top:10px;" v-if="currentMember && currentMember.is_owner">删除分组时, 用户分组信息等都将被删除</Alert>
        </div>
        <Modal v-model="openAdd" width="360" @on-ok="addTeam">
            <p>
                <Input v-model="username" placeholder="新分组名字, 如: 销售部门, 运营部门."/>
            </p>
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
            username: "",            
            members: [],
            currentMember: {},
        }
    },
    methods: {
        goImages() {
            this.$router.push({'name':'images'});
        },
        getMembers() {
            if(!this.organization || !this.organization.id) return;
            api.getOrganizationUsers(this.organization.id).then(res => {
                var current = res.data.filter(it => it.user_id == this.user.id);
                if(current && current.length) {
                    this.currentMember = current[0];
                }

                this.members = res.data;
            });
        },
        addTeam() {
            if(!this.username) {
                this.$Message.warning("请输入分组名");
                return;
            }
            var data = Object.assign([], this.organization.teams);
            data.push({name: this.username});
            api.updateOrganization(this.organization.id, {"teams":data}).then(res => {
                this.$store.dispatch('GetOrganizations');
                this.$Message.success("保存成功");
            }, err => {
                this.$Notice.error({title:'保存出错', desc: err.response.data.message});
            });
        },
        removeTeam(index) {
            var data = this.organization.teams;
            data.splice(index, 1);
            api.updateOrganization(this.organization.id, {teams: data}).then(res => {
                this.$store.dispatch('GetOrganizations');
                this.$Message.success('删除成功');
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        }
    },
    computed: {
        organization() {
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
    /* width: 500px; */
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
    padding: 5px;
    padding: 10px 15px 10px 5px;
}
tbody tr:hover {
    background-color: #f8f8f9;
}
.team-name {
    font-size: 1.3em;
}
</style>
