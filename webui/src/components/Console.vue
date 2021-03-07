<template>
  <div id="app" :class="{'show-left-nav': showNav}">
    <nav-left-bar v-show="showNav" :show="true"></nav-left-bar>
    <app-main></app-main>
  </div>
</template>

<script>
// import Navbar from './Navbar'
import NavLeftBar from './NavLeftBar'
import AppMain from './AppMain'
import * as types from '../store/types'
import {mapState} from 'vuex'
import api from '../axios'

let token = localStorage.getItem("token");

export default {
  name: 'app',
  components: {
    NavLeftBar,
    AppMain
  },
  methods: {
    getUser() {
      let self = this;
      api.mySelf().then(res => {
        // console.log(res)
        self.$store.dispatch('SetUser', res.data);
        self.getOrganizations();
      }, err => {
        
      });
    },
    getOrganizations() {
      this.$store.dispatch('GetOrganizations');
    },
    toggleNav() {
      // this.$store.commit(types.SHOW_NAV, !this.showNav)
    },
    getFileUploadToken() {
      let self = this;
      api.getFileUploadToken().then(function(res) {
          self.$store.dispatch('SetFileUploadToken', res.data);
      });
    }
  },
  computed: mapState({
    showNav: state => state.showNav,
    user: state => state.user,
  }),
  mounted() {
    // 只有用户登录的情况下才去fetch相关数据
    if(token) {
      this.getUser();
      this.getFileUploadToken();
    }
  }
}
</script>

<style>
html{height:100%;}
body{
  margin:0;
  height:100%;
  overflow:hidden;
}
#app {
  /*font-family: 'Avenir', Helvetica, Arial, sans-serif;*/
  font-family: "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei","微软雅黑",Arial,sans-serif;
  height: 100%;
  background: #f8f8f9;  
}
</style>
