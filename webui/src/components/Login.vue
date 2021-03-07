<template lang="html">
  <div style="text-align:center;margin-top:50px;">
    <img src="../assets/logo.png" style="margin-top:20px;width:60px;height:60px;border-radius: 5px;cursor:pointer;" @click="goHome">
    <h3 style='margin-bottom:50px;'>{{$config.APP_NAME}}</h3>
    <Row>
      <Col span="6" offset="9">
        <Form :model="dynamicValidateForm" ref="dynamicValidateForm">
            <Input v-model="dynamicValidateForm.username" placeholder="用户名" ></Input>
            <Input type="password" v-model="dynamicValidateForm.password" placeholder="密码"></Input>
            <Button type="primary" long @click.stop="submitForm('dynamicValidateForm')">登录</Button>
            <!--<el-button @click="resetForm('dynamicValidateForm')">重置</el-button>-->
        </Form>
        <div class="m-t-md" v-if="true">
            还无帐号? 
            <router-link to="/register">点击注册新帐号</router-link>
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
            dynamicValidateForm: {
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
        submitForm(formName) {
            
            let self = this;
            let opt = this.dynamicValidateForm;
            
            api.login(opt).then(function(response) {
                localStorage.setItem("username", self.dynamicValidateForm.username);
                self.$store.dispatch('UserLogin', response.data);
                // self.$store.dispatch('SetUser', response.data);
                let redirect = decodeURIComponent(self.$route.query.redirect || '/console/tables');
                self.$router.push({
                    path: redirect
                })
            }).catch(function(error) {
                self.$Message.error(error.response.data.message);
            });
        
        }
    },
    created() {
        if(localStorage.getItem('token')) {
            this.$router.push({name:'tables'});
        }
        if(localStorage.getItem('username')) {
            this.dynamicValidateForm.username = localStorage.getItem('username');
        }
    }
}
</script>
<style scoped>

.ivu-input-wrapper {
    margin-bottom: 20px;
}

</style>