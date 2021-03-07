<template>
    <div class="row " 
        :class="{'selected': record.is_selected, 'active': activeRow==record.id}" 
        :style="{'width':width + 'px'}"
        >

        <div class="freezed-data" style="width:85px;">
            <div class="col center" style="width:40px;cursor:pointer;" @click="record.is_selected=!record.is_selected">
                <span class="ivu-checkbox">
                    <Icon type="android-checkbox" class="text-primary" size="20" v-if="record.is_selected"></Icon>
                    <span class="ivu-checkbox-inner" size="20" v-if="!record.is_selected"></span>
                </span>
            </div><!--
            --><div class="col center" style="width:45px;">{{rowIndex+1}}</div>
        </div>
        <div class="common-data" style="margin-left:85px;">
            <div class="col" v-for="(column, colIndex) in visiableColumns" :key="column.id" 
                            :style="{width:column.width+'px', 'text-align': column.text_align}" 
                            :class="{selected: isActiveRow && column.id == activeCol }"
                            @click="selectCell(record, column)"
                            >
                            <!-- @dblclick="editCell(record, column)" -->
                            <!-- @contextmenu.prevent="selectCell(record, column, true)" -->
                <div class="container" >
                    <template v-if="column.type=='boolean'" class="ivu-checkbox" >
                        <Icon type="android-done" class="text-primary" size="20" v-if="record.cells[column.id]"></Icon>
                    </template>
                    <template v-if="['string','text','number'].includes(column.type)" :title="record.cells[column.id]">
                        {{record.cells[column.id]}}
                    </template>      
                    <template v-if="['datetime'].includes(column.type)&&record.cells[column.id]">
                        <span v-if="column.config && column.config.is_contain_time">{{record.cells[column.id] | moment('YYYY-MM-DD HH:mm')}}</span>
                        <span v-else>{{record.cells[column.id] | moment('YYYY-MM-DD')}}</span>
                    </template>   
                    <template class="text-gray" v-if="['system'].includes(column.type)">
                        <span v-if="column.config.type=='updated_at'">{{record.updated_at | moment('YYYY-MM-DD HH:mm')}}</span>
                        <span v-else>{{record.created_at | moment('YYYY-MM-DD HH:mm')}}</span>
                    </template> 
                    <div class="selection-wrapper" v-if="column.type=='selection'">
                        <span class="selection-value" v-for="option in column.options" v-if="record.cells[column.id].includes(option.id)" :style="{'background-color': option.background_color}">{{option.name}}</span>
                    </div>                       
                    <template v-if="column.type=='image'">
                        <span class="image-item" v-for="(item, index) in (record.cells[column.id]||[])" @click="openImage(item, index)">
                            <img :src="item.url + '?imageView2/1/w/100'">
                        </span>                           
                    </template>
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
            return this.$store.state.permission;
        }        
    },   
    methods: {
        selectCell(record, column, isContextClick) {
            this.$emit("selectCol", record, column);
            if(isContextClick) {
                this.isOpenContextMenu = true;
            }

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
    z-index: 10;
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
    border-color: #dd985f;
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
.image-item img{
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