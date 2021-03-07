<template>
    <div class="table-container">
        <!-- 打印表格头部 -->
        <table-header :table="table" :columns="visiableColumns" 
            v-on:refresh="refresh"
            v-on:selectAll="selectAll"
            v-on:updateCol="updateEditingCol"
            ></table-header>
        <div id="table-body" class="table-body" >
            <!--<Spin size="large" fix v-if="isLoading && !records.length"></Spin>-->
            <div class="table-data">
                <div class="table-content clearfix"  :style="{'height':(contentHeight + 400) + 'px'}">
                    <table-row v-for="(record, rowIndex) in records" :key="record.id" 
                        :active-row="activeRow" 
                        :active-col="activeCol" 
                        :editing-col="editingCol"
                        :columns="visiableColumns" :record="record" 
                        :row-index="rowIndex" 
                        :width="rowWidth"
                        v-on:selectRow="selectRow"
                        v-on:selectCol="selectCol"
                        v-on:editCol="editCol" 
                        ></table-row>                
                    <!-- 增加一行空的 -->
                    <div @click="addRow()" class="add-row" v-if="permission.edit_data!='none'">
                            <Tooltip content="新增一行"  placement="bottom-start">
                                <Icon size="20" type="ios-plus-empty"></Icon>
                            </Tooltip>
                    </div>
                </div>
            </div>
            <div class="col-editor " :type="activeCol.type" :class="{'editing': editingCol && editingCol.id}"
                :style="{top:offset.top+'px',left:offset.left+'px',width:offset.width+'px'}" v-if="activeCol">
                <template v-if="isHasEditPermission(activeRow) && editingCol">
                    <textarea autofocus="autofocus" id="col_editor_textarea"  class="col-editor-text" 
                        v-if="['string','text','number'].includes(editingCol.type)" 
                        v-model="editingCol.cellValue" ></textarea>
                </template>
            </div>
            <div class="popup" :type="activeCol.type" v-if="activeCol && isHasEditPermission(activeRow)" :style="{top:(offset.top)+'px',left:offset.left+'px'}">
                <div class="datetime" v-if="editingCol && editingCol.type=='datetime'">
                    <Date-picker
                        :open="editingDialog"
                        :value="editingValue"
                        confirm
                        :type="editingCol.config && editingCol.config.is_contain_time ? 'datetime' : 'date'"
                        @on-change="handleDateChange"
                        @on-clear="handleDateClear"
                        @on-ok="handleDateOk"><span></span>
                    </Date-picker>
                </div>
                <div class="selection" v-if="editingCol && editingCol.type=='selection'">
                    <div class="selection-item float-container" v-for="option in editingCol.options" @click="toggleSelectionColItem(editingCol, option)">
                        <span>
                            <Icon size="20" style="verticle-align:middle" type="ios-circle-outline" v-if="!activeRow.cells[editingCol.id].includes(option.id)"></Icon>
                            <Icon size="20" style="verticle-align:middle" type="ios-checkmark" v-if="activeRow.cells[editingCol.id].includes(option.id)"></Icon>
                        </span>
                        <span>{{option.name}}</span>
                        <div class="pull-right" style="height:18px;width:18px;margin-top:11px;" :style="{'background-color': option.background_color}"></div>
                    </div>
                </div>
                <div class="selection" v-if="editingCol && editingCol.type=='member'">
                    <div class="selection-item float-container" v-for="member in members" @click="toggleSelectionColItem(editingCol, member)">
                        <span v-if="true">
                            <Icon size="20" style="verticle-align:middle" type="ios-circle-outline" v-if="!activeRow.cells[editingCol.id].includes(member.id)"></Icon>
                            <Icon size="20" style="verticle-align:middle" type="ios-checkmark" v-if="activeRow.cells[editingCol.id].includes(member.id)"></Icon>
                        </span>
                        <span>{{member.nickname}}</span>
                        <div class="pull-right" style="height:18px;width:18px;margin-top:11px;" :class="{'men': member.user_sex==1, 'women': member.user_sex==2, 'unknown': member.user_sex==0}"></div>
                    </div>
                </div>
                <div class="image" :style="{width:(offset.width)+'px',height:'40px'}" v-if="activeCol && activeCol.type=='image'">
                    <span class="file-add-button"  >
                        <Upload
                            ref="upload"
                            :show-upload-list="false"
                            :default-file-list="[]"
                            :on-success="uploadSuccess"
                            :format="['jpg','jpeg','png']"
                            :max-size="1024"
                            :on-format-error="handleFileFormatError"
                            :on-exceeded-size="handleFileMaxSize"
                            :on-progress="handleProgress"
                            :before-upload="handleBeforeUpload"
                            :action="fileUploadToken.upload_url"
                            :data="{'token':fileUploadToken.token}"
                            style="">
                            <Icon type="plus" size="20"  style="margin-top:10px;"></Icon>
                        </Upload>
                        
                    </span>
                    <!--<div style="margin-right:50px;">
                        <Progress :percent="uploadPercent" :hide-info="true" :stroke-width="3" ></Progress>
                    </div>-->
                </div>
            </div>
            
        </div>
        <!--<edit-column :open="isEditColumn" :column="editingColumn" v-on:complete="refresh"></edit-column>-->
    </div>
    
</template>
<script>
import Vue from 'vue'
import TableHeader from './TableHeader'
import TableRow from './TableRow'
import TableRowReadonly from './TableRowReadonly'
import TableColumn from './TableColumn'
import {mapState} from 'vuex'
import api from '../../axios'
export default {
    props: ['table', 'columns'],
    components: {
        TableHeader,
        TableRow,
        TableRowReadonly,
        TableColumn,
    },
    data() {
        return {
            page: 1,
            pagesize: 100,
            freezedWidth: 85,
            uploadPercent: 0,
            isEditing: false,
            rowIds: {},
            colIds: {},
            colValue: null, // 当编辑前的数据, 用于做对比, 如果未修改col值, 则不调用更新步骤
            editingValue: "",    
            editingDialog: true,   // 非文本框, 有弹出框的模式        
            activeRow: null,
            activeCol: null,
            editingCol: null,
        }
    },
    computed: {
        isLoading() {
            return this.$store.state.isLoading;
        },
            // return this.$store.state.records;
        records: {
            get() {
                return this.$store.state.records
            },
            set(value) {
                this.$store.commit('table_records', value)
            }            
        },
        fileUploadToken() {
            return this.$store.state.fileUploadToken;
        },
        rowWidth() {
            var width = 80;
            this.visiableColumns.forEach(col => width += col.width);
            width += 200; // 多加点空间
            return width;
        },
        bodyHeight() {
            let obj = document.getElementById('#table-body');
            return obj.offsetTop;
        },
        contentHeight() {
            var height = 80 + this.records.length * 40;
            return height;
        },
        permission() {
            return this.$store.state.permission || {};
        },
        members() {
            return this.$store.state.members || [];
        },
        visiableColumns() {
            return this.$store.getters.visiableColumns;
        },        
        offset() {
            var offset = {top: 0, left: 0, width: 0, height: 0};
            if (this.activeRow && this.activeCol) {
                offset.width = this.activeCol.width+1;
                offset.top = this.rowIds[this.activeRow.id] * 40 - 1;
                
                var left = this.freezedWidth*1-1;
                for(var i = 0; i < this.visiableColumns.length;i++){
                    if(this.visiableColumns[i].id==this.activeCol.id) break;
                    left += this.visiableColumns[i].width;
                }
                offset.left = left;
            }
            
            return offset;
        }        
    },
    methods: {
        isHasEditPermission(record) {
            if(!this.permission.edit_data || this.permission.edit_data == 'none') return false; // 未设置权限
            if(this.permission.edit_data =='user' && this.permission.user_id != record.user_id) return false; // 只有本人数据权限时, 当前行不是自己创建的
            return true;
        },
        openEditColumn(column) {
            if(!column) {
                column = {"name":"", "type":"string"}
            }
            // this.editingColumn = column;
            this.isEditColumn = true
        },
        refresh() {
            this.$emit("refresh");
        },
        keydown(e) {
            // if(!e) return;
            // console.log(e.keyCode);
            if(e.keyCode == 27) {
                // 如果编辑状态, 则回滚修改的内容
                this.rollebackChanges();
            } else if(e.keyCode == 9) {
                this.updateEditingCol(this.activeRow, this.activeCol);
                e.preventDefault();
            } else if(e.keyCode == 13) {
                if(this.editingCol && /string|number/.test(this.editingCol.type)) {
                    this.updateEditingCol(this.activeRow, this.activeCol);
                    e.preventDefault();
                    return;
                }
                if(this.activeRow && this.activeCol && !this.editingCol) {
                    this.editCol(this.activeRow, this.activeCol)
                }
                // e.preventDefault();
            } 

            // 编辑状态下不允许跳格
            if((!this.editingCol || !this.editingCol.id) && this.activeCol) {
                if(e.keyCode == 37) { // left
                    let idx37 =  this.colIds[this.activeCol.id];
                    if(idx37 > 0) { 
                        this.selectCol(this.activeRow, this.visiableColumns[idx37-1]);
                    }
                    e.preventDefault();
                } else if(e.keyCode ==39) { // right
                    let idx39 =  this.colIds[this.activeCol.id];
                    if(idx39 < this.visiableColumns.length-1){ 
                        this.selectCol(this.activeRow, this.visiableColumns[idx39+1]);
                    }
                    e.preventDefault();
                } else if(e.keyCode == 38) { // up
                    let idx38 = this.rowIds[this.activeRow.id];
                    if(idx38 > 0) { 
                        this.selectCol(this.records[idx38-1], this.activeCol);
                    }
                    e.preventDefault();
                } else if(e.keyCode == 40) { // down
                
                    let idx40 =  this.rowIds[this.activeRow.id];
                    if(idx40 < this.records.length-1) { 
                        this.selectCol(this.records[idx40+1], this.activeCol);
                    }
                    e.preventDefault();
                }            
            }
        },
        updateEditingCol(record, column, callback) {
            if(!this.editingCol) return;

            // 检查当前编辑字段是否真改变了. 若未做修改, 则不去更新后台数据库
            if(this.colValue == this.editingCol.cellValue) {
                this.editingCol = null;
                callback && callback();
                return;
            }

            if(!record) record = this.activeRow;
            if(!column) column = this.editingCol;

            this.isEditing = true;

            var item = {};
            item[this.editingCol.id] = this.editingCol.cellValue;
            let cells = {"cells": item};
            this.saveToServer(this.activeRow, cells, (ok) => {
                callback && callback();
                this.isEditing = false;
            });
            this.activeCol = this.editingCol;

        },
        // 把row/col的id单独做成一个map, 方便表格里跳转
        setIdIndex() {
            var rowIds = {};
            this.records.forEach((row, index) => rowIds[row.id] = index);
            this.rowIds = rowIds;

            var colIds = {};
            this.visiableColumns.forEach((col, index) => colIds[col.id] = index);
            this.colIds = colIds;
        },
        selectAll(is_selected) {
            this.records = this.records.map(function(row) {
                row.is_selected = is_selected;
                return row;
            });
        },
        selectRow(record, index) {
            // debugger;
            //this.records[index].is_selected = !this.records[index].is_selected;
        },
        selectCol(record, column) {
            // 如果有更新的值, 先做更新
            if(this.editingCol) {
                this.updateEditingCol(this.activeRow, this.editingCol, () => {
                    this.activeRow = Object.assign({}, record);
                    this.activeCol = column;        
                });
                return;
            }
            console.log(record, column);
            // return;
            this.activeRow = Object.assign({}, record);
            this.activeCol = column;
        },
        rollebackChanges() {
            if(!this.editingCol) return;
            this.editingCol.cellValue = this.colValue;
            this.activeCol = this.editingCol;
            this.editingCol = null;
        },
        editCol(record, column) {
            // debugger;
            if(!record || !column) return;
            if(!this.isHasEditPermission(record)) return;

            this.colValue = this.activeRow.cells[column.id];
            column.cellValue = this.activeRow.cells[column.id];
            this.editingCol =  Object.assign({}, column);
            if(column && column.type=='datetime') {
                this.editingDialog = true;
            } else if(column && /string|number|text/.test(column.type)) {
                setTimeout(() => document.getElementById("col_editor_textarea").focus(), 50);
            }
        },
        addRow() {
            // console.log(this.table);
            let self = this;

            // 如果有更新的值, 先做更新
            if(this.editingCol) {
                this.updateEditingCol(this.activeRow, this.editingCol, () => {
                    api.addRecord({tableId: self.table.id}).then(res => {
                        self.records.push(res.data);
                    }, err => {
                        this.$Message.error(err.response.data.message);
                    });
                });
                return;
            }
            api.addRecord({tableId: self.table.id}).then(res => {
                self.records.push(res.data);
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        saveToServer(record, data, callback) {
            let self = this;
            self.$store.dispatch('IsLoading', true);
            api.editRecordCell(record.table_id, record.id, data).then(res => {
                self.$store.dispatch('IsLoading', false);
                let rowIndex = self.rowIds[record.id];
                // Object.assign(self.records[rowIndex].cells[self.editingCol.id] = ""
                let cells = Object.assign(self.records[rowIndex].cells, data.cells);
                self.records[rowIndex].cells = cells;
                // self.selectCol(record, this.editingCol);
                self.editingCol = null;
                self.colValue = null;
                callback && callback(true);
            }, err => {
                self.$store.dispatch('IsLoading', false);
                this.$Message.error(err.response.data.message);
                callback && callback(false);
            });
        },
        handleDateChange(date) {
            this.editingValue = date;
        },        
        handleDateClear() {
            let item = {};
            item[this.editingCol.id] = "";
            let cells = {"cells": item};
            this.saveToServer(this.activeRow, cells);
            // let rowIndex = this.rowIds[this.activeRow.id]
            // this.records[rowIndex].cells[this.editingCol.id] = ""
            // this.setActive(this.activeRow, null);
        },
        handleDateOk() {
            let date = new Date(this.editingValue);
            let item = {};
            item[this.editingCol.id] = date.toISOString();
            let cells = {"cells": item};
            this.saveToServer(this.activeRow, cells);

            // let rowIndex = this.rowIds[this.activeRow.id]
            // this.records[rowIndex].cells[this.editingCol.id] = date.toISOString();
            // this.setActive(this.activeRow, null);
        },
        uploadSuccess (res, file) {
            this.$store.dispatch('IsLoading', false);

            let imgURL = this.fileUploadToken.url_prefix + "/" + res.key;
            let item = {};
            item["url"] = imgURL;
            var data = (this.activeRow.cells[this.activeCol.id] || []);
            data.push(item);
            let cells = {"cells": data};


            let patchData = {};
            patchData[this.activeCol.id] = data;

            this.saveToServer(this.activeRow, {"cells":patchData});
            // let rowIndex = this.rowIds[this.activeRow.id]
            // this.records[rowIndex].cells[this.editingCol.id] = item;
            // this.setActive(this.activeRow, null);
            this.$Message.success("上传图片成功");
        },
        handleFileFormatError (file) {
            this.$Notice.warning({
                title: '文件格式不正确',
                desc: '文件 ' + file.name + ' 格式不正确，请上传 jpg 或 png 格式的图片。'
            });
        },
        handleFileMaxSize (file) {
            this.$Notice.warning({
                title: '超出文件大小限制',
                desc: '文件 ' + file.name + ' 太大，不能超过 1M。'
            });
        },
        handleProgress(e, file, list) {
            console.log(file.percentage);
            this.uploadProgress = file.percentage;
        },
        handleBeforeUpload () {
            const check = (this.activeRow.cells[this.activeCol.id] || []).length < 9;
            if (!check) {
                this.$Notice.warning({
                    title: '最多只能上传 9 张图片。'
                });
            }
            // 检查下token是否过期
            if((new Date(this.fileUploadToken.expires_at)) < (new Date)) {
                this.$Notice.warning({
                    title: '正在刷新上传token, 请重新上传!'
                });
                let self = this;
                api.getFileUploadToken().then(function(res) {
                    self.$store.dispatch('SetFileUploadToken', res.data);
                });
                return false;
            }

            this.$store.dispatch('IsLoading', true);
            return true;
        },
        toggleSelectionColItem(column, option) {
            
            let isAllowMultiSelect = column.config && column.config.is_multiple || false;
            let isSelected = this.activeRow.cells[column.id].includes(option.id);
            var val = this.activeRow.cells[column.id] || [];

            // 不允许多选时, 首先清空已选
            if(!isAllowMultiSelect) {
                val = [];
            }
            if(isSelected) {
                val = val.filter(item => item != option.id)
            } else {
                val.push(option.id);
            }
            let item = {};
            item[column.id] = val;
            let cells = {"cells": item}
            this.saveToServer(this.activeRow, cells);

            // let rowIndex = this.rowIds[this.activeRow.id]
            // this.records[rowIndex].cells[column.id] = val;
            
            //console.log(this.record.cells[column.id]);
        },
        clear() {
            this.activeRow = null;
            this.activeCol = null;
            this.editingCol = null;
        }   
    },
    watch: {
        table(oldVal, newVal) {
            // console.log(oldVal,newVal);
        },
        records(newVal, oldVal) {
            this.setIdIndex();
            this.clear();
            
        }
    },
    created() {
        window.addEventListener('keydown', this.keydown);
    },
}
</script>

<style scoped>


.table-container {
    overflow-y: hidden;
    overflow-x: scroll;
    width: 100%;
    height: 100%;
    position: relative;
}
.table-body {
    min-width:100%;
    overflow-y: auto;
    /* overflow-x: visible; */
    position:absolute;
    top: 42px;
    bottom:0;
}
.table-data {
    overflow-x: visible;
}
.add-row {
    border-bottom:1px solid #d7dde4;
    border-right:1px solid #d7dde4;
    width:45px;
    height:35px;
    line-height:40px;
    text-align:center;
    cursor:pointer;
}

.col-editor {
    top: 0;
    position: absolute;
    z-index: 10;
    border: 1px solid #dd985f;
    padding: 0;
    margin: 0;
    height: 41px;
    /*box-shadow: 0 0 0.1rem 0 #dd985f;*/
    border-color: #dd985f;
    overflow:hidden;
    resize: none;
}
.col-editor:not(.editing) {
    pointer-events: none;
}
.col-editor[type=string] textarea, .col-editor[type=number] textarea{
    white-space: pre;
    word-wrap: normal;
    overflow-y:  hidden;
}
.col-editor[type=text].editing {
    height: 95px;
}
.col-editor:not(.editing) .col-editor-text{
    opacity: 0.0001;
}
.col-editor .col-editor-text {
    height: 100%;
    width: 100%;
    padding: 10px;
    font-size: 14px;
    border: 0;
    resize: none;
    outline: 0;
    background: white;
}
.file-add-button {
    position: absolute;
    right: 0;
    height: 100%;
    width: 50px;
    background: rgba(255,255,255, 0.9);
    text-align:center;
    cursor: pointer;
    color: #ccc;
    border-left: 1px solid #ededed;
}
.popup {
    position: absolute;
    z-index: 100;

}
.popup .datetime {
    margin-top: 40px;
}
.popup .selection {
    width: 200px;
    margin-top:40px;
    padding: 3px 0;
    border: 1px solid #dd985f;
    background:white;
}
.popup .image {
    /*background: rgba(0,0,0,0.3);*/
}
.ivu-select-dropdown {
    margin-top: 0 !important;
}
.file-add-button {
    height:38px;
    top:1px;
    right:1px;
}
.selection-item {
    height: 40px;
    line-height: 40px;
    cursor: pointer;
    padding: 0 20px;
}
.selection-item:hover {
    background-color: #f3f3f3;
}
.selection-item span {
    margin-right: 5px;
    display: inline-block;
}
.selection-item i {
    vertical-align:middle;
}
.selection-value {
    padding: 2px 6px;
    border-radius: 1px;
    color: #222;
    margin-right: 5px;

}
.ivu-upload, .ivu-upload-select {
    height: 100%;
}
.ivu-upload-select i {
    vertical-align:middle;
}
</style>