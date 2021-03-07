<template>
    <div class="container">
        <template-header :template="template"></template-header>
        <div class="content">
            <div class="yue">
                <div class="header">
                    已有{{columns.length}}列
                    <Button class="right" type="success" @click.native="addColumn">添加新列</Button>
                </div>
                <div>
                    <Row>
                        <Col span="8">
                            <div class="dropdown-pop-wrapper">
                                <draggable v-model="columns" :options="{group:'people'}" @end="updateColumnOrders">
                                    <div v-for="column in columns" class="draggable-item" :key="'sort_'+column.id" @click="selectColumn(column)"
                                        :class="{'selected':selectedColumn.id==column.id}">
                                        <Icon type="images" v-if="column.type=='image'"></Icon>
                                        <Icon type="ios-grid-view-outline" v-if="column.type=='number'"></Icon>
                                        <Icon type="android-create" v-if="column.type=='string'"></Icon>
                                        <Icon type="navicon" v-if="column.type=='text'"></Icon>
                                        <Icon type="ios-calendar-outline" v-if="column.type=='datetime'"></Icon>
                                        <Icon type="ios-checkmark-outline" v-if="column.type=='boolean'"></Icon>
                                        <Icon type="chevron-down" v-if="column.type=='selection'"></Icon>
                                        <Icon type="ios-information-outline" v-if="column.type=='system'"></Icon>                    
                                        
                                        <strong v-if="column.is_primary">{{column.name}}</strong>
                                        <span v-if="!column.is_primary">{{column.name}}</span>

                                        <Icon type="navicon" class="right text-gray" size="18" style="margin-top:8px;"></Icon>
                                    </div>
                                </draggable>
                            </div>
                            <Alert type="warning" style="margin-top:10px;" v-if="columns.length">
                                点击上面列可编辑内容
                            </Alert>
                        </Col>
                        <Col span="10" offset="1" class="form" v-if="selectedColumn.id">

                            <Form :model="selectedColumn"  >
                                <FormItem label="列名">
                                    <Input v-model="selectedColumn.name" placeholder=""></Input>
                                </FormItem>
                                <FormItem label="类型">
                                    <Select v-model="selectedColumn.type" @on-change="typeChanged">
                                        <Option v-for="item in types" :value="item.value" :key="item.value">{{ item.label }}</Option>
                                    </Select>
                                </FormItem>
                                
                                <FormItem label="" v-if="selectedColumn.type=='datetime'">
                                    <Checkbox v-model="selectedColumn.config.is_contain_time"> 包含时间</Checkbox>
                                </FormItem>
                                <FormItem label="" v-if="selectedColumn.type=='number'">
                                    <Select v-model="selectedColumn.config.precision" >
                                        <Option v-for="item in precisions" :value="item.value" :key="item.value">{{ item.label }}</Option>
                                    </Select>
                                </FormItem>
                                <FormItem label="" v-if="selectedColumn.type=='selection'">
                                    <div>
                                        <div class="option" v-for="(option, index) in selectedColumn.options" :key="option.id">
                                            <Input v-model="option.name" style="width:200px;" placeholder="名称"></Input>
                                            <Poptip placement="bottom" width="180">
                                                <span class="color-picker" :style="{'background-color':option.background_color}"></span>
                                                <div slot="content">
                                                    <div>
                                                        <Input v-model="option.background_color" size="small" style="width:100px;" placeholder="颜色hex代码"></Input>
                                                        <Icon type="loop" size="18" style="vertical-align:middle;" @click.native="option.background_color=randomBackgroundColor()"></Icon>
                                                    </div>
                                                    <div class="float-container" style="margin-top:5px;">
                                                        <div v-for="color in colors" :key="color" class="pull-left" style="width:30px;height:30px;" :style="{'background-color':color}" @click="option.background_color=color"></div>
                                                    </div>
                                                </div>
                                            </Poptip>
                                            <Icon type="ios-minus-outline" size="24" style="vertical-align:middle;" @click.native="removeOption(index)"></Icon>
                                        </div>
                                    </div>
                                    <a type="text" class="add-button" @click.stop="addOption">增加选项</a>
                                    <div style="margin-top: 15px">
                                        <Checkbox v-model="selectedColumn.config.is_multiple"> 允许多选</Checkbox>
                                    </div>
                                </FormItem>
                                <FormItem label="" v-if="selectedColumn.type=='system'">
                                    <Radio-group v-model="selectedColumn.config.type" vertical @on-change="radioChanged">
                                        <Radio label="created_at">
                                            <span>创建时间</span>
                                        </Radio>
                                        <Radio label="updated_at">
                                            <span>更新时间</span>
                                        </Radio>
                                        <Radio label="created_by">
                                            <span>创建人</span>
                                        </Radio>
                                        <Radio label="updated_by">
                                            <span>更新人</span>
                                        </Radio>
                                    </Radio-group>
                                </FormItem>

                                <FormItem label="">
                                    <Button type="primary" @click="saveColumnChange">保存</Button>
                                    <Poptip
                                        confirm
                                        title="您确定删除此列吗？"
                                        placement="left"
                                        class="delete-button"
                                        @on-ok="deleteColumn"
                                        >
                                        <a class="right" style="color:#ed3f14;" >删除本列</a>
                                    </Poptip>
                                </FormItem>
                            </Form>
                        </Col>
                    </Row>
                </div>
            </div>
        </div>
    </div>  
</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'
import TemplateHeader from './TemplateHeader'
import draggable from 'vuedraggable'
import Vue from 'vue'
export default {
    components: {
        TemplateHeader,
        draggable,
    },
    data () {
        return {
            template: {},
            selectedColumn: {},
            columns: [],
            types: [
                { value: 'string', label: '单行文本'},
                { value: 'text', label: '多行文本'},
                { value: 'number', label: '数值'},
                { value: 'datetime', label: '日期时间'},
                { value: 'selection', label: '下拉选项'},
                { value: 'boolean', label: '复选框'},
                { value: 'image', label: '图片'},
                // { value: 'relation', label: '关联'},
                { value: 'system', label: '系统字段'},
            ],
            precisions: [
                { value: 0, label: "0"},
                { value: 1, label: "1"},
                { value: 2, label: "2"},
                { value: 3, label: "3"},
                { value: 4, label: "4"},
                // { value: 5, label: "5"},
                // { value: 6, label: "6"},
            ],
            colors: [
                "#f2c61f","#bcdfc3","#3b83c0","#e07b53","#a59a3b","#d95c5c","#7cb329",
                "#dacef2", "#d6d6d6", "#464c5b", "#657180", "#FE84CD"
            ]            
        }
    },
    methods: {
        typeChanged() {
            var col = Object.assign({}, this.selectedColumn);

            if(col.type=='system' && !col.config.type) {
                col.config.type = 'created_at';
            }
            if(col.type=='selection' && !col.options) {
                col.options = [];
            }
            if(col.type=='number' && !col.config.precision) {
                col.config.precision = 0;
            }
            console.log(col);
            
            this.selectedColumn = col;
        },
        radioChanged() {
            console.log(this.selectedColumn);
            this.selectColumn = Object.assign({}, this.selectedColumn, true);
        },
        selectColumn(col) {
            var newCol = Object.assign({}, col);
            console.log(newCol);
            this.selectedColumn = col;
        },
        addColumn() {
            var data = {name: "未命名", "type": "string"};
            api.addTemplateColumn(this.template.id, data).then(res => {
                this.getTemplateColumns();
                this.$Message.success("新增成功");
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        updateColumnOrders() {
            var ids = this.columns.map(item => item.id);
            api.updateTemplateColumnOrders(this.template.id, {"ids":ids}).then(res => {
                this.$Message.success('修改列显示顺序成功!');
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },        
        save() {
            var data = {'name': this.template.name, description: this.template.description, is_active: this.template.is_active};
            api.updateTemplate(this.template.id, data).then(res => {
                this.getTemplate();
                this.$Message.success("保存成功");
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        deleteColumn() {
            api.deleteTemplateColumn(this.template.id, this.selectedColumn.id).then(res => {
                this.selectedColumn = {};
                this.getTemplateColumns();
                this.$Message.success("删除成功");
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        getTemplateColumns() {
            api.getTemplateColumns(this.template.id).then(res => {
                this.columns = res.data;
            });
        },
        getTemplate() {
            api.getTemplate(this.$route.params.templateId).then(res => {
                this.template = res.data;
                this.getTemplateColumns();
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        addOption() {
            if(!this.selectedColumn.options) {
                this.selectedColumn.options = [];
            }
            this.selectedColumn.options.push({name: "", 'background_color': this.randomBackgroundColor()});
        },    
        removeOption(index) {
            this.selectedColumn.options.splice(index,1);
        },                    
        randomBackgroundColor() {
            return '#'+(0x1000000+(Math.random())*0xffffff).toString(16).substr(1,6);
        },        
        saveColumnChange() {
            if(!this.selectedColumn.name) {
                this.$Message.error("请输入列名称");
                this.isLoading = false;
                return false;
            }

            if(this.selectedColumn.type!='system') {
                delete this.selectedColumn.config.type;
            }

            var column = {"name": this.selectedColumn.name, "type": this.selectedColumn.type, 
                "options": this.selectedColumn.options, "config": this.selectedColumn.config, 
                "width": this.selectedColumn.width};
            
            api.updateTemplateColumn(this.template.id, this.selectedColumn.id, column).then(res => {
                this.getTemplateColumns();
                this.$Message.success("保存成功");
            }, err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
            

        },        
    },
    created() {
        this.getTemplate();
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
    min-height: 500px;
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
.right{
    float:right;
}
.dropdown-pop-wrapper {
    border: 1px solid #ccc;
}
.draggable-item {
    height: 35px;
    line-height: 35px;
    cursor: pointer;
    padding: 0 10px;
    border-bottom: 1px solid #dedede;
}
.draggable-item:hover {
    background: #eee;
}
.draggable-item.selected {
    background: #fbf5ed;
}
.form {
    padding: 0 30px;
    border: 1px solid #ccc;
}
.delete-button {
    float: right;
    cursor: pointer;
    float: right;
}
.option {
    margin-top: 5px;
}
.color-picker {
    display: inline-block;
    border-radius: 3px;
    width: 30px;
    height: 30px;
    vertical-align:middle;
    margin: 0 5px;
}
</style>