<template>
    <div >
        <Modal title="" :loading="isLoading" v-model="isOpen"  width="300" @on-ok="ok" @on-cancel="close">
            <div class="container">
                <div class="title">名称</div>
                <div>
                    <Input v-model="name" placeholder="" ></Input>
                </div>
                <div class="title">类型</div>
                <div>
                    <Select v-model="type" >
                        <Option v-for="item in types" :value="item.value" :key="item">{{ item.label }}</Option>
                    </Select>
                </div>      
                <div v-show="type=='selection'">
                    <div class="title">选项</div>     
                    <div>
                        <div class="option" v-for="(option, index) in options" :key="option.id">
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
                            <Icon type="ios-minus-outline" size="24" style="vertical-align:middle;" @click.native="removeColor(index)"></Icon>
                        </div>
                    </div>
                    <a type="text" class="add-button" @click="addOption">增加选项</a>
                </div>
            </div>
        </Modal>
    </div>
</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'

var randomColor = function(){ return '#'+(0x1000000+(Math.random())*0xffffff).toString(16).substr(1,6); }

export default {
    props: ['column', 'open'],
    data() {
        return {
            isOpen: false,
            isLoading: true,
            name: "",
            type: "string",
            options: [],
            types: [
                { value: 'string', label: '单行文本'},
                { value: 'text', label: '多行文本'},
                { value: 'number', label: '数值'},
                { value: 'datetime', label: '日期时间'},
                { value: 'selection', label: '下拉选项'},
                { value: 'boolean', label: '复选框'},
                { value: 'image', label: '图片'},
            ],
            colors: [
                "#f2c61f","#bcdfc3","#3b83c0","#e07b53","#564f8a","#d95c5c","#00b5ad",
                "#dacef2", "#d6d6d6", "#464c5b", "#657180", "#FE84CD"
            ]
        }
    }, 
    methods: {
        ok() {
            if(!this.name) {
                this.$Message.error("请输入列名称");
                this.isLoading = false;
                return false;
            }

            this.isLoading = true;

            var column = {"name": this.name, "type": this.type, "options": this.options,
                "tableId": this.table.id, "userId": this.user.id};
            let self = this;
            api.createColumn(column).then(function(response) {
                self.close();
                self.$Message.success("创建列成功!");
                self.$emit("refresh");
            });
            
        },
        close() {
            this.$emit("close");
        },
        removeColor(index) {
            this.options.splice(index,1);
        },
        addOption() {
            this.options.push({name: "", 'background_color': randomColor()});
        },
        randomBackgroundColor() {
            return randomColor();
        }
    },
    computed: mapState({
        user: state => state.user,
        table: state => state.table,
    }),
    watch: {
        open: function(oldVal, newVal) {
            this.isOpen = newVal;
            console.log(oldVal, newVal);
        }
    }
}    
</script>
<style scoped>
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
</style>