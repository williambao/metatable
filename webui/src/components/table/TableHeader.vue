<template>
    <div class="header">
    <div class="freezed-data" style="width:85px;" >
       <div class="col col-fixed" style="width:45px;">&nbsp;</div><!--
        --><div class="col col-fixed" style="width:40px" @click="selectAll">
            <span class="ivu-checkbox">
                <Icon type="android-checkbox" class="text-primary" size="20" v-if="isSelectedAll"></Icon>
                <span class="ivu-checkbox-inner"  v-if="!isSelectedAll"></span>
            </span>
        </div>
    </div>
    <div class="common-data" style="margin-left:85px;">
        <div :id="'col_'+column.id" class="col"  v-for="(column, index) in visiableColumns" :key="column.id" 
                v-bind:style="{width:column.width+'px'}" :class="{resizing:selectedColumn != null &&column.id==selectedColumn.id}">
                <div class="resize-box" @mousedown.left="startCalcWidth(column, index)"></div>
                <div class="col-name">
                    <span :class="{primary:column.is_primary}">
                        <Icon type="images" v-if="column.type=='image'"></Icon>
                        <Icon type="ios-grid-view-outline" v-if="column.type=='number'"></Icon>
                        <Icon type="android-create" v-if="column.type=='string'"></Icon>
                        <Icon type="navicon" v-if="column.type=='text'"></Icon>
                        <Icon type="ios-calendar-outline" v-if="column.type=='datetime'"></Icon>
                        <Icon type="ios-checkmark-outline" v-if="column.type=='boolean'"></Icon>
                        <Icon type="chevron-down" v-if="column.type=='selection'"></Icon>
                        <Icon type="ios-information-outline" v-if="column.type=='system'"></Icon>
                        <Icon type="link" v-if="column.type=='relation'"></Icon>
                        <Icon type="person" v-if="column.type=='member'"></Icon>
                        
                        {{column.name}}
                    </span>
                    <Dropdown class="dropdown-icon" trigger="click" @on-click="columnAction">
                        <Icon type="md-arrow-dropdown" />
                        <Dropdown-menu slot="list">
                        <!--
                            <Dropdown-item :name="'edit'">编辑数据</Dropdown-item>
                            <Dropdown-item divided>排序</Dropdown-item>
                            <Dropdown-item :name="'hidden_'+column.id">隐藏此列</Dropdown-item>
                            -->
                            <Dropdown-item :name="'sort_'+column.id+'_ase'" divided>排序(1,2,3)</Dropdown-item>
                            <Dropdown-item :name="'sort_'+column.id+'_desc'" >倒序(3,2,1)</Dropdown-item>
                            <Dropdown-item :name="'setting_'+column.id" divided v-if="permission.is_table_admin">设置列</Dropdown-item>
                            <Dropdown-item :name="'delete_'+column.id" v-if="permission.is_table_admin">删除列</Dropdown-item>
                        </Dropdown-menu>
                    </Dropdown>
                </div>
                
            </div><!--
            --><div class="col col-add-col" style="width:60px;text-align:center;cursor:pointer;" v-if="permission.is_table_admin" @click="openEditColumn(null)">
                <Tooltip content="新增一列" >
                    <Icon type="md-add" />
                </Tooltip>
            </div>
            <!-- 编辑字段弹框 -->
            <Modal title="" v-if="editingColumn!=null" v-model="isEditingColumn"  
                width="300" :styles="{top: '160px'}" @on-ok="saveColumnChange">
                <div class="container">
                    <div class="title">名称</div>
                    <div>
                        <Input v-model="editingColumn.name" placeholder="" ></Input>
                    </div>
                    <div class="title">类型</div>
                    <div>
                        <Select v-model="editingColumn.type" >
                            <Option v-for="item in types" :value="item.value" :key="item.value">{{ item.label }}</Option>
                        </Select>
                    </div>

                    <div class="config-bar" v-if="editingColumn.type=='member'">
                        <div style="margin-top: 15px">
                            <Checkbox v-model="editingColumn.config.is_multiple"> 允许多选</Checkbox>
                        </div>
                    </div>
                    <div class="config-bar" v-if="editingColumn.type=='selection'">
                        <div class="title">选项</div>     
                        <div>
                            <div class="option" v-for="(option, index) in editingColumn.options" :key="option.id">
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
                            <Checkbox v-model="editingColumn.config.is_multiple"> 允许多选</Checkbox>
                        </div>
                    </div>
                    <div v-if="editingColumn.type=='datetime'" style="margin-top: 15px">
                        <Checkbox v-model="editingColumn.config.is_contain_time"> 包含时间</Checkbox>
                    </div>
                    <div v-if="editingColumn.type=='relation'" style="margin-top: 15px">
                        <div class="title">表格</div>
                        <div>
                            <Select v-model="editingColumn.config.table_id" >
                                <Option v-for="item in tables" :value="item.id" :key="item">{{ item.name }}</Option>
                            </Select>
                        </div>
                        <div style="margin-top: 15px">
                            <Checkbox v-model="editingColumn.config.is_multiple"> 允许多选</Checkbox>
                        </div>
                    </div>
                    
                    <div class="config-bar" v-if="editingColumn.type=='number'">
                        <div class="title">显示精度</div>
                        <div>
                            <Select v-model="editingColumn.config.precision" >
                                <Option v-for="item in precisions" :value="item.value" :key="item">{{ item.label }}</Option>
                            </Select>
                        </div>
                    </div>
                    <div class="title" v-if="editingColumn.type=='system'">
                        <Radio-group v-model="editingColumn.config.type" vertical>
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
                    </div>
                    <!--<div class="title">显示宽度</div>
                    <div>
                        <Slider v-model="editingColumn.width" :min="90" :max="600" @on-input="updateColumnWidth"></Slider>
                    </div>-->
                </div>
            </Modal>

        </div>
    </div> 
</template>
<script>
import api from '../../axios'
import {mapState} from 'vuex'
import * as types from '../../store/types'
export default {
    props: ["columns", "table"],
    components: {
        // "resizer": require "vue-resize-handle/bidirectional"
    },
    data() {
        return {
            editingColumn: null,
            editingColumnOriginal: null,
            isEditingColumn: false,
            isSelectedAll: false,
            selectedColumn: null,
            selectedColumnStartPosition: 0,

            types: [
                { value: 'string', label: '单行文本'},
                { value: 'text', label: '多行文本'},
                { value: 'number', label: '数值'},
                { value: 'datetime', label: '日期时间'},
                { value: 'selection', label: '下拉选项'},
                { value: 'boolean', label: '复选框'},
                { value: 'member', label: '团队成员'},
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
        startCalcWidth(col, index) {
            // 若有编辑过的内容, 则先做保存
            this.$emit('updateCol');

            this.selectedColumn = col;
            var el = document.getElementById('col_'+col.id);
            this.selectedColumnStartPosition = getElementLeft(el);
            document.body.style.cursor = 'col-resize';
        },
        caclColumnWidth(e) {
            if(!this.selectedColumn) return;
            var width = e.clientX - this.selectedColumnStartPosition;
            if(width < 100) width = 100;
            console.log(width);
            this.selectedColumn.width = width;
        },
        stopResize() {
            if(!this.selectedColumn) return;
            document.body.style.cursor = 'default';
            if(this.permission.is_table_admin) {
                this.editingColumn = this.selectedColumn;
                this.saveColumnChange();
            }
            this.selectedColumn = null;
        },
        openEditColumn(column) {
            if(!column) {
                column = {name: "", type:"string", options: [], config: {type:'created_at'}, width: 180}
            }
            this.editingColumnOriginal = column;
            this.editingColumn = Object.assign({}, column);
            this.isEditingColumn = true;
        },
        addOption() {
            if(!this.editingColumn.options) {
                this.editingColumn.options = [];
            }
            var _options = Object.assign([], this.editingColumn.options);
            _options.push({name: "", 'background_color': this.randomBackgroundColor()});
            this.editingColumn.options = _options;
        },    
        removeOption(index) {
            this.editingColumn.options.splice(index,1);
        },            
        columnAction(key) {
            
            let self = this;
            let id = key.split("_")[1];
            // console.log(id);
            let column = this.getColumnById(id);
            // console.log(column);
            if(/^delete/.test(key)) {
                this.openDeleteDialog(id);
            } else if (/^setting/.test(key)) {
                this.openEditColumn(column);    
            } else if (/^hidden/.test(key)) {
                var cols =  [...this.hiddenColumns];
                cols.push(id);
                this.$store.commit(types.TABLE_HIDDEN_COLUMNS, cols);
                this.updateHiddenColumns();
            } else if(/^sort/.test(key)) {
                this.updateSorts(id, /desc/.test(key));
            }
        },
        updateSorts(id, is_desc) {
            let _sorts = this.sorts;
            _sorts.push(is_desc ? '-'+id : id);
            let self = this;
            let data = {'sorts':_sorts};
            api.updateTableView(this.table.id, this.view.id, data)
            .then(function(res){
                self.$store.dispatch('SetView', res.data);
                self.$emit('refreshData');
            });
        },
        selectAll() {
            this.isSelectedAll = !this.isSelectedAll;
            this.$emit('selectAll', this.isSelectedAll);
        },
        openDeleteDialog(id) {
            let self = this;
            this.$Modal.confirm({
                title: '确认',
                content: '<p>删除列后无法再恢复, 确认要继续删除吗?</p>',
                width: 320,
                onOk: () => {
                    api.deleteColumn(self.table.id, id).then(function(response) {
                        self.$Message.success("删除列成功");
                        self.$emit("refresh");
                    });        
                },
                onCancel: () => {
                    
                }
            });
        },
        getColumnById(id) {
            var column = null;
            this.columns.forEach(function(col) {
                if(col.id == id) {
                    column = col;
                }
            });
            return column;
        }, 
        updateColumnWidth(value) {
            if(!this.editingColumnOriginal) return;
             this.editingColumnOriginal.width = value;
        },        
        saveColumnChange() {
            if(!this.editingColumn.name) {
                this.$Message.error("请输入列名称");
                this.isLoading = false;
                return false;
            }
            this.isLoading = true;

            if(this.editingColumn.type!='system') {
                delete this.editingColumn.config.type;
            }

            var column = {"name": this.editingColumn.name, "type": this.editingColumn.type, 
                "options": this.editingColumn.options, "config": this.editingColumn.config,
                "tableId": this.table.id, "userId": this.user.id, "width": this.editingColumn.width};
            if(this.editingColumn.id) {
                column["id"] = this.editingColumn.id;
            }

            let self = this;

            var callback = function(res){
                // self.close();
                self.$Message.success(self.editingColumn.id ? "修改成功" :"创建列成功!");
                self.$emit("refresh");
            }

            var errs = function(err) {
                self.$Message.error(err.response.data.message || '未知错误');
            }

            if(this.editingColumn.id) {
                api.updateColumn(column).then(callback, errs);
            } else {
                api.createColumn(column).then(callback, errs);
            }

        },
        randomBackgroundColor() {
            return '#'+(0x1000000+(Math.random())*0xffffff).toString(16).substr(1,6);
        },
        updateHiddenColumns() {
            let self = this;
            let data = {'hide_columns':this.tableHiddenColumns};
            api.updateTableView(this.table.id, this.selectedView.id, data)
                .then(res => {
                    self.getTable();
                }, err => {
                    self.$Message.error(err.message || '未知错误');
                });
        }
    },
    computed: {
        tables() {
            return this.$store.state.tables;
        },
        user() {
            return this.$store.state.user;
        },
        visiableColumns() {
            return this.$store.getters.visiableColumns;
        },
        sorts() {
            return this.$store.getters.sorts;
        },
        permission() {
            return this.$store.state.permission || {};
        }
    },
    created() {
        // this.$store.dispatch('ShowNav', false);
        window.addEventListener('mousemove',this.caclColumnWidth);
        window.addEventListener('mouseup', () => {
            this.stopResize();
            
        });    
    },
    destroyed: function() {
        window.removeEventListener('mousemove', this.caclColumnWidth);
    } 
}
function getElementLeft(ele){
    var viewportOffset = ele.getBoundingClientRect();
    // these are relative to the viewport, i.e. the window
    var top = viewportOffset.top;
    var left = viewportOffset.left;
    return left;
　　}
</script>
<style scoped>
.freezed-data {
    position: absolute;
    background: #fff;
    /*z-index: 10;*/
}
.header {
    font-size: 1.2em;
    white-space:nowrap;
    color: #828282;
}
.header .col {
    text-align:center;
}
.header .ivu-dropdown-item {
    text-align: left;
}

.header i {
    vertical-align:middle;
}
.col {
    position: relative;
    display: inline-block;
    height: 100%;
    /*border: 1px solid #d7dde4;*/
    border-right: 1px solid #d7dde4;
    border-bottom: 1px solid #d7dde4;
    vertical-align:middle;
    padding: 0 10px;
    text-align:center;
    -webkit-user-select: none;  
    -moz-user-select: none;    
    -ms-user-select: none;      
    user-select: none; 
}

.dropdown-icon {
    cursor: pointer;
    display: none;
    position: absolute;
    right: 10px;
}
.col-name {
    padding: 10px 0;
}
.col-name:hover .dropdown-icon {
    display: inline-block;
}
.col-add-col, .col-fixed {
    padding: 10px 0;
}
.col .primary {
    font-weight: 900;
    color: #dd985f;
}
.ivu-checkbox-wrapper {
    margin: 0;
}
.title {
    margin-top: 10px;
}
.container {
    text-align: left;
}
.add-button {
    margin-top: 10px;
    color: #555;
    display: block;
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
.config-bar- {
    background-color:#efefef;
    border-radius: 6px;
    margin:10px 0;
    padding: 10px;
}
.resize-box {
    width: 4px;
    height: 100%;
    background: transparent;
    position: absolute;
    right: -2px;
    top: 0;
    cursor: col-resize; 
}
.resize-box:hover, .resizing .resize-box {
    display:block;
    background: #dd985f;
}
</style>