<template>
<div class="layout">
    <header class="layout-header">
        <Menu class="layout-menu" mode="horizontal" theme="light" active-name="" @on-select="go">
            <div class="layout-logo">
                <a :href="$config.WEB_URL" class="wrapper-header-nav-logo router-link-active">
                    <img src="../assets/logo.png" style="vertical-align:middle;width:30px;height:30px;border-radius:3px;">
                    <strong>{{$config.APP_NAME}} </strong>
                    <span style="font-size:10px;" v-if="false">beta</span>
                </a>
            </div>
            <div class="layout-nav">
                <MenuItem name="product" v-if="false">
                    产品与功能
                </MenuItem>
                <!--<MenuItem name="fangan">
                    解决方案
                </MenuItem>-->
                <MenuItem name="anli" v-if="false">
                    案例
                </MenuItem>
                <MenuItem name="contactus" v-if="false">
                    联系我们
                </MenuItem>
            </div>
            <div class="layout-signin">
                <MenuItem name="tables" class="login-button" v-if="isLogined">
                    <Button type="primary">控制台</Button>
                </MenuItem>
                <MenuItem name="signin" class="login-button" v-if="!isLogined">
                    登录
                </MenuItem>
                <MenuItem name="signup" class="register-button"  v-if="!isLogined">
                    注册
                </MenuItem>
            </div>
        </Menu>
    </header>
    <div class="layout-content">
        <div class="carousel">
            <Row>
                <Col>
                    <div class="title"><h1>记录点滴, 参知伴行</h1></div>
                    <div class="sub-title">工作,生活从此井然有序. 随时随地查询管理.</div>
                </Col>
            </Row>
        </div>
        <div class="product">
            <h1>产品功能</h1>
            <Row class="images" :gutter="30">
                <Col span="6" offset="3">
                    <img src="../assets/product1.png" style="width:100%;height:100%;"/>
                    <div class="title">动态表格</div>
                    <div class="sub-title">就如Excel一样熟悉</div>
                </Col>
                <Col span="6">
                    <img src="../assets/product2.png" style="width:100%;height:100%;"/>
                    <div class="title">多端同步</div>
                    <div class="sub-title">可在电脑及手机上同步编辑</div>
                </Col>
                <Col span="6">
                    <img src="../assets/product3.png" style="width:100%;height:100%;"/>
                    <div class="title">导入导出</div>
                    <div class="sub-title">可导入导出多种文件格式</div>
                </Col>
            </Row>
        </div>
    </div>
    <footer class="footer">
        <Row>
            <Col span="10" >
                <span>&copy; 2021 <strong>&middot;</strong> </span>
                <span>{{$config.APP_NAME}}</span>
                <span style="margin-left:30px;">
                    <a href="http://www.miitbeian.gov.cn/">粤ICP备17157833号 </a>
                </span>
            </Col>
            <Col>
                <span style="float:right;">电话: 186 7564 2414</span>
            </Col>
        </Row>
    </footer>
</div>
</template>
<script>
export default {
    data () {
        return {
            carousel: 1,
            isLogined: false,
        }
    },
    methods: {
        go(name) {
            if(name=='signin'){
                this.$router.push({name:'login'});
            } else if(name=='signup'){
                this.$router.push({name:'register'});
            } else if(name=='tables'){
                this.$router.push({name:'tables'});
            }
        }
    },
    created() {
        // let token = localStorage.getItem("token");
        // // 已经登录的用户直接跳转到表格列表页面
        // if(token) {
        //     this.$router.push({name:'tables'});
        // }
        if(localStorage.getItem('token')) {
            this.isLogined = true;
        }
    }
}
</script>
<style scoped>
    .layout{
        background: #f5f7f9;
        height: 100%;
        overflow-y: auto;
    }
    .layout-header {
        background-image: linear-gradient(312deg, rgb(55, 58, 146) 0%, rgb(12, 92, 149) 100%);
    }
    .layout-logo{
        height: 40px;
        line-height: 40px;
        border-radius: 3px;
        float: left;
        position: relative;
        top: 10px;
        left: 20px;
    }
    .layout-logo a {
        color:#333;
    }
    .layout-signin {
        float: right;
        margin-right: 15px;
    }
    .layout-menu {
        height: 62px;
    }
    .layout-nav {
        margin-left: 20px;
        padding-left: 150px;
        
    }
    .layout-content{
        min-height: 500px;
        margin: 0;
        overflow: hidden;
        background: #fff;
    }
    .carousel {
        width: 100%;
        /* background: url(../../src/assets/home-bg.jpg) top center no-repeat #11214e; */
        /* background-image: linear-gradient(312deg, rgb(55, 58, 146) 0%, rgb(12, 92, 149) 100%); */
        background-image: linear-gradient(312deg, rgb(30, 150, 229) 0%, rgb(12, 92, 149) 100%);
        min-height: 600px;
        text-align:center;
    }
    .carousel .title {
        font-size: 36px;
        color: white;
        margin-top: 200px;
    }
    .carousel .sub-title {
        font-size: 24px;
        color: #d8d8d8;
        margin-top: 10px;
    }
    .product {
        padding: 80px 0;

        text-align:center;
    }
    .product h1 {
        font-size: 36px;
    }
    .product .images {
        margin-top: 30px;
    }
    .product .title {
        margin-top: 10px;
        font-size: 24px;
        font-weight: 400;
    }
    .product .sub-title {
        margin-top: 5px;
        font-size: 16px;
        color: gray;
    }
    .footer {
        color: #9ea7b4;
        background: #f2f4f6;
        border-top: #d8dde4;
        padding: 20px 40px;
        font-size:14px;
    }
    .footer a {
        color: #9ea7b4;
    }
    .login-button {
        padding: 0 8px !important;
    }
    .register-button {
        padding: 0 8px !important;
    }
</style>