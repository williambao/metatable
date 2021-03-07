<template >
  <div class="container" style="">
    <img src="../assets/logo.png" style="margin-top:20px;width:60px;height:60px;border-radius: 5px;cursor:pointer;" @click="goHome">
    <h3 style='margin-bottom:50px;'>{{$config.APP_NAME}}</h3>
    <Row>
      <Col span="6" offset="9">
        <Form :model="newItem">
                <FormItem label="公司名:">
                    <Input v-model="newItem.organization_name" placeholder="比如: 如:深圳市柔然科技有限公司" ></Input>
                </FormItem>
                <FormItem label="名称:">
                    <Input v-model="newItem.nickname" placeholder="比如: William" ></Input>
                </FormItem>
                <FormItem label="账户Email(用于登陆):">
                    <Input v-model="newItem.email" placeholder="比如: zhang@qq.com"></Input>
                </FormItem>
                <FormItem label="登陆密码:">
                    <Input v-model="newItem.password" placeholder="至少6位以上, 不要用太过简单的密码"></Input>
                </FormItem>
                <FormItem label="手机号:">
                    <Input v-model="newItem.phone" placeholder="限中国大陆11位手机号" ></Input>
                </FormItem>
                <FormItem label="邀请码:">
                    <Input v-model="newItem.invite_code" placeholder="请联系客服人员索取" ></Input>
                </FormItem>
                <FormItem label="">
                    <Button type="primary" long @click.native="submitForm">提交</Button>
                </FormItem>
        </Form>

        <div class="redirect">
            <router-link to="/login">点击返回登陆页面</router-link>
        </div>
    </Col>
    </Row>
  </div>
</template>
<script>
import * as types from '../store/types'
import api from '../axios'
export default {
    name: '',
    data() {
        return {
            newItem: {
                email: '',
                username: '',
                password: ''
            },
            activeName: this.$store.state.activeName,
        }
    },
    methods: {
        handleClick(tab, event) {},
        // 重置
        resetForm(formName) {
            this.$refs[formName].resetFields();
        },
        goHome() {
            this.$router.push({name:'index'});
        },
        // 登录
        submitForm() {
            let self = this;
            let opt = this.newItem;

            // 检查下用户名邮箱号

            api.register(opt).then(function(response) {
                self.$store.dispatch('UserLogin', response.data);
                let redirect = decodeURIComponent(self.$route.query.redirect || '/console/tables');
                self.$router.push({
                    path: redirect
                })
            }).catch(function(error) {
                self.$Message.error(error.response.data.message);
            });
        
        }
    }
}
</script>
<style scoped>
.container {
    text-align:center;
    margin-top:50px;
    height: 100%;
    padding-bottom: 100px;
    overflow-y: auto;
}
.ivu-input-wrapper {
    margin-bottom: 20px;
}
.remark {
    
}
.redirect {
    margin-top: 30px;
}

</style>