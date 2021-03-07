<template>
    <div class="row " 
        :class="{'selected': record.is_selected, 'active': activeRow && activeRow.id==record.id}" 
        :style="{'width':width + 'px'}"
        >
        <div class="freezed-data" style="width:85px;">
            <div class="col center" style="width:45px;">{{rowIndex+1}}</div><!--
            --><div class="col center" style="width:40px;cursor:pointer;" @click="record.is_selected=!record.is_selected">
                <span class="ivu-checkbox">
                    <Icon type="android-checkbox" class="text-primary" size="20" v-if="record.is_selected"></Icon>
                    <span class="ivu-checkbox-inner"  v-if="!record.is_selected"></span>
                </span>
                <!--<Checkbox v-model="record.is_selected" @change.native="selectRow"></Checkbox>-->
            </div>
            </div>
        <div class="common-data" style="margin-left:85px;">
            <div class="col"  v-for="(column, colIndex) in visiableColumns" :key="column.id" 
                            :style="{width:column.width+'px', 'text-align': column.text_align}" 
                            :class="{selected: isActiveRow && column.id == activeCol.id }"
                            @click="selectCell(record, column)"
                            @dblclick="editCell(record, column)"
                            >
                            <!-- @contextmenu.prevent="selectCell(record, column, true)" -->
                <div class="container">

                    <span v-if="column.type=='boolean'" class="ivu-checkbox" @click="updateBooleanColumn(column)">
                        <Icon type="android-checkbox" class="text-primary" size="20" v-if="record.cells[column.id]"></Icon>
                        <span class="ivu-checkbox-inner"  v-if="!record.cells[column.id]"></span>
                    </span>
                    
                    <span v-if="['string','text'].includes(column.type)" :title="record.cells[column.id]">
                        {{record.cells[column.id]}}
                    </span>
                    <span v-if="['number'].includes(column.type)" :title="record.cells[column.id]">
                        {{record.cells[column.id] && parseFloat(record.cells[column.id]).toFixed(column.config.precision)}}
                    </span>
                    <span v-if="['datetime'].includes(column.type)&&record.cells[column.id]">
                        <span v-if="column.config && column.config.is_contain_time">{{record.cells[column.id] | moment('YYYY-MM-DD HH:mm')}}</span>
                        <span v-else>{{record.cells[column.id] | moment('YYYY-MM-DD')}}</span>
                    </span>

                    <span v-if="['system'].includes(column.type)" class="text-gray">
                        <span v-if="column.config.type=='created_at'">{{record.created_at | moment('YYYY-MM-DD HH:mm')}}</span>
                        <span v-else-if="column.config.type=='updated_at'">{{record.updated_at | moment('YYYY-MM-DD HH:mm')}}</span>
                        <span v-else-if="column.config.type=='created_by'">{{record.user_name}}</span>
                        <span v-else-if="column.config.type=='updated_by'">{{record.update_user_name}}</span>
                    </span>
                    
                    <template v-if="column.type=='image'">
                        <span class="image-item" v-for="(item, index) in (record.cells[column.id]||[])" @click="openImage(item, index)">
                            <img :src="item.url + '?imageView2/1/w/100'">
                        </span>
                        <Modal title="查看图片" v-model="isPreviewImage">
                            <img :src="currentImage.url + '?imageView2/2/w/600'" style="width:100%;">
                            <div slot="footer" >
                                <a @click="deleteImage" v-if="isHasEditPermission(record)">
                                    <Icon type="trash-a" size="20"></Icon>
                                    删除此图片
                                </a>
                            </div>
                        </Modal>
                    </template>

                    <div class="selection-wrapper" v-if="column.type=='selection'">
                        <span class="selection-value" v-for="option in column.options" v-if="record.cells[column.id] && record.cells[column.id].includes(option.id)" :style="{'background-color': option.background_color}">{{option.name}}</span>
                    </div>
                    <div class="selection-wrapper" v-if="column.type=='member'">
                        <span class="selection-value" v-for="member in members" v-if="record.cells[column.id] && record.cells[column.id].includes(member.id)" :class="{'men': member.user_sex==1, 'women': member.user_sex==2, 'unknown': member.user_sex==0}">{{member.nickname}}</span>
                    </div>
                    
                    
                    <!-- 当右击鼠标时弹出菜单 
                    <Dropdown v-if="isActiveRow && column.id == activeCol" class="context-menu" trigger="custom" :visible="isOpenContextMenu" @on-click="">
                        <Dropdown-menu slot="list">
                            <Dropdown-item >编辑数据</Dropdown-item>
                            <Dropdown-item divided>排序</Dropdown-item>
                            <Dropdown-item>隐藏此列</Dropdown-item>
                            <Dropdown-item  divided>设置列</Dropdown-item>
                            <Dropdown-item  >删除列</Dropdown-item>
                        </Dropdown-menu>
                    </Dropdown>-->
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import api from '../../axios'
export default {
    props: ["record", "columns", "width", "rowIndex", "selectedRows", "activeRow", "activeCol", "editingCol"],
    data() {
        return {
            editingValue: "",
            isActive: false,
            isOpenContextMenu: false,
            currentImage: {},
            isPreviewImage: false,
            isTextEditOpened: false,
        }
    },
    computed: {
        isActiveRow() {
            return this.activeRow && this.activeRow.id == this.record.id
        },
        editingColumn() {
            return this.editingCol;
        },
        fileUploadToken() {
            return this.$store.state.fileUploadToken;
        },
        visiableColumns() {
            return this.$store.getters.visiableColumns;
        },
        permission() {
            return this.$store.state.permission || {};
        },
        members() {
            return this.$store.state.members || [];
        },              
    },    
    methods: {
        isHasEditPermission(record) {
            if(!this.permission.edit_data || this.permission.edit_data == 'none') return false; // 未设置权限
            if(this.permission.edit_data =='user' && this.permission.user_id != record.user_id) return false; // 只有本人数据权限时, 当前行不是自己创建的
            return true;
        },        
        openImage(item, index) {
            this.currentImage = item;
            this.isPreviewImage = true;
        },
        deleteImage() {
            let item = {};
            var data = (this.record.cells[this.activeCol.id] || []);
            data.splice(this.openImageIndex, 1);
            item[this.activeCol.id] = data;
            let cells = {"cells": item};

            // console.log(cells);
            this.saveToServer(cells);
            this.record.cells[this.activeCol] = data;
            // this.$emit("editCol", null);
            this.$Message.success("删除图片成功");
        },
        selectRow() {
            this.$emit("selectRow", this.record, this.rowIndex);
        },
        selectCell(record, column, isContextClick) {
            // 区分单击, 双击操作
            clearTimeout(window.selectCellTimer);
            let self = this;
            // window.selectCellTimer = setTimeout(function() {
                if(self.editingColumn != column.id) {
                    self.editEndCell(record, column);
                }

                self.$emit("selectCol", record, column);
                
                if(isContextClick) {
                    self.isOpenContextMenu = true;
                }
            // }, 200);

        }, 
        editEndCell(record, column) {
            this.$emit("setActive", record, null);
            this.editingValue = "";
        },
        editCell(record, column) {
            clearTimeout(window.selectCellTimer);

            this.$emit("editCol", record, column);
        },   

        openTextEditPage(column) {
            this.$emit("editCol", column);
            this.isTextEditOpened = true;
        },

        updateBooleanColumn(column) {
            if(!this.isHasEditPermission(record)) return;
            this.record.cells[column.id] = !this.record.cells[column.id];
            this.updateColumn(column);
        },
        updateColumn(column) {
            let item = {};
            item[column.id] = this.record.cells[column.id];
            let cells = {"cells": item};
            this.saveToServer(cells);
        },
        saveToServer(data) {
            let self = this;
            self.$store.dispatch('IsLoading', true);
            api.editRecordCell(this.record.table_id, this.record.id, data).then(res => {
                self.$store.dispatch('IsLoading', false);
            }, err => {
                self.$store.dispatch('IsLoading', false);
                this.$Message.error(err.response.data.message);
            });
        },
        checkNumberInput(evt) {
            if(evt.which < 48 || evt.which > 57) {
                return false;
            }
            return true;
        },

    },

}
</script>
<style scoped>

.row {
    white-space:nowrap;
}
.row i {
    vertical-align:middle;
}
.row.selected .col{
    background: #FBF5EE;
}
.row.active .col{
    background: #fdfaf6;
}
.freezed-data {
    position: absolute;
    background: #fff;
    /*z-index: 10;*/
}
.small .col {
    height: 30px;
    line-height:30px;
}
.col {
    position: relative;
    display: inline-block;
    height: 100%;
    /*border: 1px solid #d7dde4;*/
    border-right: 1px solid #d7dde4;
    border-bottom: 1px solid #d7dde4;
    vertical-align:middle;
    height: 40px;
    line-height:40px;
    font-size: 14px;
    color: #222;
    
}
.col textarea, .col input {
    width:100%;
    height:100%;
    
    border: 0;
    resize: none;
    background-color: transparent;
    outline-color: #dd985f;
    outline-width: 0;
    white-space: nowrap;
    overflow-x: auto;
}
.col .editing  {
    background: #fff;
}

.col.row-selected {
    background: #fdfaf6;
    
}
.col.selected {
    background: #fbf5ed;
}

.col .container {
    padding: 0 10px;
    height: 41px;
    margin: -1px;
    border-width: 1px;
    border-style: solid;
    border-color: transparent;
    text-overflow:ellipsis;
    overflow: hidden;
}
.col.selected .container {
    /* box-shadow: 0 0 0.1rem 0 #dd985f; */
    /*border-color: #dd985f;*/
}
.ivu-checkbox-wrapper {
    margin: 0;
}
.col .ivu-poptip {
    position: absolute;
    right: 5px;
    top: 0px;
    
}
.col .ivu-poptip i {
    color: #888;
}
.col .ivu-poptip .ivu-poptip-body {
    padding: 0 !important;
}
.selection-wrapper {
    width: 100%;
    margin-right: 20px;
    overflow: hidden;
    user-select: none;
    -webkit-user-select: none;
}
.selection-pop-wrapper {
    padding: 5px 0;
}
.selection-item {
    height: 40px;
    line-height: 40px;
    cursor: pointer;
    
}
.selection-item:hover {
    background-color: #f3f3f3;
}
.selection-item span {
    margin-right: 5px;
    display: inline-block;
}
.selection-value {
    padding: 2px 6px;
    border-radius: 1px;
    color: #222;
    margin-right: 5px;

}

.image-item {
    display: inline-block;
    width: 30px;
    height: 30px;
    text-align: center;
    line-height: 30px;
    border: 1px solid transparent;
    border-radius: 4px;
    overflow: hidden;
    background: #fff;
    position: relative;
    box-shadow: 0 1px 1px rgba(0,0,0,.2);
    margin-right: 4px;
    vertical-align: middle;
    cursor: pointer;
}
.image-item img {
    width: 100%;
    height: 100%;
}
.ivu-picker-confirm > span {
    color: #657180;
}
.ivu-dropdown-item {
    line-height: 15px;
    text-align: left;
}


</style>
