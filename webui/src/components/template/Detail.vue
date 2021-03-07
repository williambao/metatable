<template>
    <div class="container">
        <template-header :template="template"></template-header>
        <div class="content">
            <Form :model="template"  v-if="template">
                <FormItem label="模板名称">
                    <Input v-model="template.name" placeholder=""></Input>
                </FormItem>
                <FormItem label="门店名">
                    <Input v-model="template.description" type="textarea" :rows="5" placeholder=""></Input>
                </FormItem>

                <FormItem label="是否正常开放">
                    <i-switch v-model="template.is_active" ></i-switch>
                </FormItem>
                <FormItem label="注册时默认初始化给用户">
                    <i-switch v-model="template.is_init" ></i-switch>
                </FormItem>
                <FormItem label="">
                    <Button type="primary" @click.native="save">保存并生效</Button>
                    <Poptip
                        confirm
                        title="您确定删除此模板吗？"
                        placement="left"
                        class="delete-button"
                        @on-ok="removeTemplate"
                        >
                        <span class="text-error" >删除本模板</span>
                    </Poptip>
                </FormItem>
            </Form>
        </div>
    </div>  
</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'
import TemplateHeader from './TemplateHeader'
export default {
    components: {
        TemplateHeader,
    },
    data () {
        return {
            template: {},
        }
    },
    methods: {
        save() {
            var data = {'name': this.template.name, description: this.template.description, is_active: this.template.is_active
                , is_init: this.template.is_init};
            api.updateTemplate(this.template.id, data).then(res => {
                this.getTemplate();
                this.$Message.success("保存成功");
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        getTemplate() {
            api.getTemplate(this.$route.params.templateId).then(res => {
                this.template = res.data;
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });
        },
        removeTemplate() {
            api.deleteTemplate(this.template.id).then(res => {
                this.$Message.success("删除成功");
                this.$router.push({name:'templates'});
            },err => {
                this.$Notice.error({title:'出错', desc: err.response.data.message});
            });        
        }        
    },
    created() {
        this.getTemplate();
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
.delete-button {
    float: right;
    cursor: pointer;
    
}
</style>