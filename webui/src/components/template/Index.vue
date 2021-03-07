<template>
<section class="container">
    <div class="header">
        <Icon type="ios-home"></Icon>
        <span>模板列表</span>
        <span class="right">
            <Button type="success" @click.native="isShowAddTableDialog=true">新增模板</Button>
        </span>
    </div>
    <div class="content">
        <div class="store" v-for="template in templates">
            <div class="store-title">{{template.name}}</div>
            <div class="store-item">
                <Row>
                    <Col span="8">
                        <span class="store-item-label">
                        <Icon type="person"></Icon>
                            创建人: 
                        </span>
                        {{template.username}}  
                    </Col>
                    <Col span="8">
                        <span class="store-item-label">
                            <Icon type="ios-calendar-outline"></Icon>
                            创建时间: 
                        </span>
                        {{template.created_at | moment('YYYY-MM-DD HH:mm')}}                
                    </Col>
                    <Col span="8">
                        <span class="store-item-label">
                        <Icon type="checkmark-circled"></Icon>
                            状态: 
                        </span>
                        <span v-if="template.is_active" class="text-success">可用</span>
                        <span v-if="!template.is_active" class="text-error">不可用</span>
                    </Col>
                </Row>
            </div>
            <div class="store-item">
                <Row>
                    <Col span="8">
                        <span class="store-item-label">
                            <Icon type="android-menu"></Icon>
                            列/字段数: 
                        </span>
                        {{template.column_count}}                
                    </Col>
                    <Col span="8">
                        <span class="store-item-label">
                        <Icon type="checkmark-circled"></Icon>
                            初始化给用户: 
                        </span>
                        <span v-if="template.is_init" class="text-success">是</span>
                        <span v-if="!template.is_init" class="text-error">否</span>
                    </Col>
                </Row>
                
            </div>
            <div  class="store-item">
                <Button @click.native="go(template, 'templatesDetail')"><Icon type="information-circled"></Icon> 模板信息</Button>
                <Button @click.native="go(template,'templatesColumn')" ><Icon type="film-marker"></Icon> 模板列/字段</Button>
            </div>

        </div>
    </div>    

        <Modal v-model="isShowAddTableDialog" title="新增模板" width="300" @on-ok="addTable">
            <div class="">
                <div>
                    <Input v-model="newTable.name" placeholder="请输入名称" ></Input>
                </div>
                <div class="m-t-md">
                    <Input v-model="newTable.description" type="textarea" :rows="4" placeholder="模板描述" ></Input>
                </div>
            </div>
        </Modal>    
</section>

</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'
import * as types from '../../store/types'

export default {
    data() {
        return {
            templates: [],
            isShowAddTableDialog: false,
            page: 1, 
            pagesize: 100,
            newTable: {
                name: "",
                description: "",
                is_private: true,
            },            
        }
    },
    methods: {
        getTemplates() {
            let self = this;
            api.getTemplates(this.page, this.pagesize).then(function(res) {
                console.log(res)
                self.templates = res.data;
            }, err => {
                this.$Notice.error({title:'保存出错', desc: err.response.data.message});
            });
        },
        go(template, routerName) {
            this.$router.push({name:routerName, params:{templateId: template.id}});
        },        
        selectTemplate(template) {
            this.$router.push({name: 'templatesDetail', params:{templateId: template.id}});
        },
        addTable() {
            if(!this.newTable.name) {
                this.$Message.info("请输入模板名称")
                return;
            }
            let self = this;
            api.createTemplate(this.newTable).then(res => {
                self.getTemplates();
            });
        },
        removeTemplate(template) {
            let self = this;
            this.$Modal.confirm({
                title: '确认',
                content: '<p>确认要继续删除吗?</p>',
                width: 320,
                onOk: () => {
                    api.deleteTemplate(template.id).then(function(response) {
                        self.$Message.success("删除成功");
                        self.getTemplates();
                    });        
                },
                onCancel: () => {
                    
                }
            });
        }
    },
    created(){
        this.getTemplates();
    }    
}
</script>
<style scoped>
.container {
    padding: 40px;
    height: 100%;
    overflow-y: auto;
}
.header {
    font-size:18px;
    height: 60px;
}
.right {
    float: right;
}
.store {
    background: white;
    padding: 30px;
    margin-bottom: 15px;
}
.store-title {
    font-size: 1.5em;
    border-bottom: 2px solid black;
    padding: 3px 0;
    margin-bottom: 15px;
}
.store-item {
    padding: 15px 0;
    border-bottom: 1px solid #f8f8f9;
}
.store-item-label {
    display: inline-block;
    width: 120px;
}
.store-item-label i {
    margin-right: 15px;
}
.plan {
    display: inline-block;
    padding:3px 5px;
    background: #d6f3fc;
    margin: 0 2px;
    border-radius: 2px;
}
.add-button:hover {
    box-shadow:0 1px 6px rgba(0, 0, 0, 0.2);
}
.ivu-tabs{
    height: 100%;
}
.right {
    float: right;
}
</style>