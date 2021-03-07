<template>
<div class="nav-left-bar">
    <div class="logo-container" @click="go('index')">
        <img src="../assets/logo.png" class="logo" >
        <!-- <Icon type="cube" size="30" class="logo"></Icon> -->
    </div>
    <div class="nav-item" @click="go('tables')" :class="{'selected': this.$route.name.includes('tables')}">
        <Icon type="logo-buffer" />
        <div class="text">表格</div>
    </div>
    <div class="nav-item" @click="go('organizationsMember')" 
        :class="{'selected': this.$route.name.includes('organizations')}" v-if="user.organization_id">
        <Icon type="md-people" />
        <div class="text">团队</div>
    </div>
    <div class="nav-item" @click="go('templates')" 
        v-if="user && user.role_ids && user.role_ids.includes('offical')"
        :class="{'selected': this.$route.name.includes('templates')}">
        <Icon type="md-clipboard" />
        <div class="text">模板</div>
    </div>
    <!-- <div class="nav-item" @click="showNotImplementWarning()">
        <Icon type="android-contacts"></Icon>
        <div class="text">好友</div>
    </div>
    <div class="nav-item" @click="showNotImplementWarning()">
        <Icon type="ios-color-filter"></Icon>
        <div class="text">团队</div>
    </div>
    <div class="nav-item" @click="showNotImplementWarning()">
        <Icon type="ios-bell"></Icon>
        <div class="text">通知</div>
    </div> -->
    <!--
    <div class="nav-item" @click="go('dashboard')">
        <Icon type="ios-folder"></Icon>
        <div class="text">模板</div>
    </div>
    <div class="nav-item">
        <Icon type="ios-paper"></Icon>
        <div class="text">门户</div>
    </div>
    
    <div class="nav-item" @click="go('users')">
        <Icon type="ios-gear"></Icon>
        <div class="text">账户</div>
    </div>
    -->
    <div class="nav-item-bottom-container">
        <div class="nav-item" @click="go('userInfoSetting')">
            <img v-if="user.avatar" class="avatar" :src="user.avatar + '?imageView2/1/w/100'" alt="头像">
            <Avatar v-else style="background-color: #87d068" icon="ios-person" size="small"/>
            <div class="text">{{user.nickname}}</div>
        </div>
        <div class="nav-item help" @click="changeNav()">
            <Poptip placement="right-end" style="z-index:99999" width="300">
                <Icon type="md-help-circle" />
                <div class="text">帮助</div>
                <div  slot="content" class="content">
                    <div class="title">
                        联系我们
                    </div>
                    <div class="phone">
                        William: 186 7564 2414
                    </div>
                    <!-- <div class="item">
                        地址: <br>深圳市南山区留仙大道塘朗城3座17A（塘朗站B出口附近）
                    </div> -->
                    <div class="item">
                        工作时间：<br>周一至周五，9:30-18:30，节假日休息
                    </div>
                    <div class="item">&nbsp;</div>
                </div>
            </Poptip>
        </div>
        <div class="nav-item" @click="logout">
            <Icon type="md-exit" />
            <div class="text">登出</div>
        </div>
    </div>
</div>
</template>
<script>
    import {mapState} from 'vuex'
    export default {
        name: "",
        props: {
            show: Boolean
        },
        computed: {
            ...mapState(['showNav', 'user'])
        },        
        data() {
            return {

            }
        },
        methods: {
            logout() {
                this.$store.dispatch('UserLogout');
                this.$router.push({name: 'index'});
            },
            go(path) {
                if(path == 'organizationsMember') {
                    this.$router.push({name: path, params: {'organizationId': this.user.organization_id}});
                    return;
                }
                this.$router.push({name:path});
            },
            showNotImplementWarning() {
                this.$Message.warning("此功能正加紧制作中:)");
            },
            changeNav() {

            }
        },
        computed: mapState({
            user: state => state.user,
        }),
    }
</script>
<style scoped>
.nav-left-bar {
    line-height: 1;
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 2;
    background-color: #353644;
    width: 50px;
    color: #fff;
    font-size: 12px;
}
.logo-container {
    text-align:center;
    padding: 10px 0;
    background: #3F3F4F;
}
.logo {
    width:30px;
    height:30px;
    cursor: pointer;
    /*border: 1px solid #fff;*/
    color: #FE4D66;
    border-radius: 15px;
}  
.nav-item {
    padding: 10px 0;
    text-align:center;
    cursor: pointer;
}
.nav-item:hover {
    background: #222;
}
.nav-item.selected {
    background: #223;
}
.nav-item i {
    font-size: 20px;
}
.nav-item .text {
    padding: 5px;
}
.nav-item-bottom-container {
    position: absolute;
    bottom: 0;
    width: 100%;
    /*background-color: #172a3d;*/
    border-top: 1px solid #3f3f3f;
}
.avatar {
    width:25px;
    height:25px;
    cursor: pointer;
    border: 1px solid #fff;
    border-radius: 15px;
    display:inline-block;
    vertical-align:middle;
}
.help .content{
    color: #222;
    width:100%;
}
.title {
    background: #0f88e5;
    color: white;
    font-size: 1.3em;
    padding: 5px 0;
    border-radius: 1px;
}
.phone {
    font-size:1.5em;
    font-weight: bold;
    color: #0f88e5;
    padding: 8px 0;
}
.item {
    padding: 5px 0;
    white-space: normal;
    text-align: left;
}
</style>