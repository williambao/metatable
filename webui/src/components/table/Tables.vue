<template>
    <section class="clearfix container">
        <Tabs value="name1">
            <Tab-pane label="所有表格" name="name1">
                <Card class="card pull-left"  v-for="table in tables" :key="table.id" @click.native="selectTable(table)">
                    <p slot="title" >
                        {{table.name}}
                        <Icon type="android-close" size="18" class="right" @click.native.stop="removeTable(table)"></Icon>
                    </p>
                    <p >
                        <span class='card-item-label'>团队:</span>
                        <span class="card-item-value">
                            <span v-if="table.organization_id">{{table.organization_name}}</span>
                            <span v-else>无</span>
                        </span>
                    </p>
                    <p>
                        <span class='card-item-label'>创建人:</span>
                        <span class="card-item-value">{{table.username}}</span>
                    </p>
                    <p>
                        <span class='card-item-label'>创建时间:</span>
                        <span class="card-item-value">{{table.created_at | moment('YYYY-MM-DD HH:mm')}}</span>
                    </p>
                </Card>
                <div class="card pull-left add-button" @click="showAddTable">
                    <Icon type="ios-plus" size="30"></Icon>
                    <h4>新增表格</h4>
                </div>
            </Tab-pane>
        </Tabs>
        <Modal v-model="isShowAddTableDialog"  :closable="false"  @on-ok="addTable">
            <div class="">
                <Form :model="newTable" >
                    <!-- <FormItem label="所有者:">
                        <Select v-model="newTable.organization_id" placeholder="">
                            <Option value="">{{user.nickname}}</Option>
                            <Option v-for="org in organizations" :value="org.id">{{org.name}}(团队)</Option>
                        </Select>
                    </FormItem> -->
                    <Row v-if="organizations && organizations.length">
                        <Col span="8" style="padding-right:5px;">
                            <FormItem label="团队:">
                                <Select v-model="newTable.organization_id" placeholder="" clearable>
                                    <Option v-for="item in organizations" :value="item.id" :key="item.id">{{ item.name }}</Option>
                                </Select>
                            </FormItem>
                        </Col>
                        <Col span="16">
                            <FormItem label="表格名称:">
                                <Input v-model="newTable.name" placeholder=""></Input>
                            </FormItem>
                        </Col>
                    </Row>
                    <FormItem label="表格名称:" v-else>
                        <Input v-model="newTable.name" placeholder=""></Input>
                    </FormItem>
                    <FormItem label="从模板初始化:">
                        <Select v-model="newTable.template_id" placeholder="" clearable>
                            <Option v-for="item in templates" :value="item.id" :key="item.id">{{ item.name }}</Option>
                        </Select>
                    </FormItem>
                </Form>
            </div>
        </Modal>
    </section>
</template>

<script>
import {mapState} from 'vuex'
import * as types from '../../store/types'
import api from '../../axios'
export default {
    name: 'Index',
    data() {
        return {
            isShowAddTableDialog: false,
            isShowTemplate: false,
            newTable: {
                name: "",
                description: "",
                organization_id: "",
                template_id: "",
                is_private: true,
            },
            templates: [],
        }
    },
    methods: {
        showAddTable() {
            if(this.organizations.length) {
                this.newTable.organization_id = this.organizations[0].id;
            }
            this.isShowAddTableDialog = true;
            this.getTemplates();
        },
        getUserTables() {
            let self = this;
            api.getUserTables(self.user.id).then(function(res) {
                console.log(res)
                self.$store.dispatch('SetTables', res.data)
            })
        },
        selectTable(table) {
            this.$store.dispatch('SetTable', table);
            this.$router.push({name: 'tablesDetail', params:{tableId: table.id}});
        },
        removeTable(table) {
            let self = this;
            this.$Modal.confirm({
                title: '确认',
                content: '<p>确认要继续删除吗?</p>',
                width: 320,
                onOk: () => {
                    api.deleteTable(table.id).then(function(response) {
                        self.$Message.success("删除成功");
                        self.getUserTables();
                    });        
                },
                onCancel: () => {
                    
                }
            });
        },
        getTemplates() {
            let self = this;
            api.getTemplates(1, 9999).then(function(res) {
                console.log(res)
                self.templates = res.data;
            }, err => {
                this.$Notice.error({title:'保存出错', desc: err.response.data.message});
            });
        },        
        addTable() {
            if(!this.newTable.name) {
                this.$Message.info("请输入表格名称")
                return;
            }
            let self = this;
            api.createTable(this.newTable).then(res => {
                self.getUserTables();
                this.$Message.success("创建表格成功!");
            }, err => {
                this.$Notice.error({title:'保存出错', desc: err.response.data.message});
            });
        }
    },
    computed: mapState({
        user: state => state.user,
        tables: state => state.tables,
        organizations: state => state.organizations,
    }),
    watch: {
        user() {
            this.getUserTables();
        }
    },
    beforeCreate() {
        this.$store.commit(types.SHOW_NAV_TABLES, false);
    },
    created(){
        this.$store.dispatch('SetTableSearchText', "");
        this.getUserTables();
    }


}
</script>

<style scoped>
.container{
    padding: 20px;
    height: 100%;
    overflow-y: auto;
}
.card {
    margin-left:15px;
    margin-top: 20px;
    width:260px;
}
.card:hover {
    cursor:pointer;
}
.card-item-label {
    display: inline-block;
    width: 90px;
}
.add-button {
    border: 1px solid #e3e8ee;
    border-radius: 4px;
    padding-top: 40px;
    text-align: center;
    height: 147px;

}
.add-button:hover {
    box-shadow:0 1px 6px rgba(0, 0, 0, 0.2);
}
.ivu-tabs{
    /* height: 100%; */
}
.right {
    float: right;
}
</style>