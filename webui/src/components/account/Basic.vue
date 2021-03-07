<template>
    <div class="container">
        <Row>
            <Col span="8" offset="8">
                <div class="avatar-container">
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
                            :before-upload="handleBeforeUpload"
                            :action="fileUploadToken && fileUploadToken.upload_url"
                            :data="{'token':fileUploadToken.token}"
                            style="display: inline-block;width:100px;">
                            <div class="avatar" style="width: 100px;height:100px;line-height: 100px;">
                                <img :src="userInfo.avatar + '?imageView2/1/w/200'" style="width: 100px;height:100px;"/>
                                <div class="upload-icon">
                                    <Icon type="ios-cloud-upload" size="25"></Icon>
                                </div>
                            </div>
                        </Upload>
                        
                    </span>
                </div>
                <Form :model="userInfo" :label-width="80">
                    <Form-item label="用户名">
                        <Input :value="userInfo.username" placeholder="用户名" disabled></Input>
                    </Form-item>
                    <Form-item label="密码">
                        <Row>
                            <Col span="18">
                                <Input type="password" value="*******" placeholder="请" disabled></Input>
                            </Col>
                            <Col span="4" offset="1">
                                <Button type="ghost" @click="showPasswordDialog=true">修改</Button>
                            </Col>
                        </Row>
                    </Form-item>
                    <Form-item label="昵称">
                        <Input v-model="userInfo.nickname" placeholder="请输入昵称"></Input>
                    </Form-item>

                    <Form-item label="性别">
                        <Radio-group v-model="userInfo.sex">
                            <Radio label="1">男</Radio>
                            <Radio label="2">女</Radio>
                        </Radio-group>
                    </Form-item>
                    <Form-item label="手机">
                        <Input v-model="userInfo.phone" placeholder=""></Input>
                    </Form-item>
                    <Form-item label="邮箱">
                        <Input v-model="userInfo.email" placeholder=""></Input>
                    </Form-item>

                    <Form-item>
                        <Button type="primary" long @click="submit">提交</Button>
                        <Button style="margin-top:10px;float:right;" size="small" @click="logout">登出帐号</Button>
                    </Form-item>
                </Form>
            </Col>
        </Row>
     <Modal
        v-model="showPasswordDialog"
        title="修改密码"
        width="300"
        @on-ok="updatePassword">
        <Form :model="password" label-position="left" :label-width="80">
            <Form-item label="旧密码">
                <Input v-model="password.old_password"></Input>
            </Form-item>
            <Form-item label="新密码">
                <Input v-model="password.new_password"></Input>
            </Form-item>
        </Form>
    </Modal>
    </div>
</template>
<script>
import {mapState} from 'vuex'
import api from '../../axios'
import * as types from '../../store/types'

export default {
    data() {
        return {
            isOpen: false,
            showPasswordDialog: false,
            password: {
                old_password: "",
                new_password: ""
            }
        }
    },
    computed: {
        userInfo() {
            return this.$store.state.user;
        },
        ...mapState([
            'user', 'fileUploadToken'
        ]),
    },  
    methods: {
        submit() {
            if(!this.userInfo.nickname) {
                this.$Message.warning("请输入昵称");
                return;
            }
            var data = {
                avatar: this.userInfo.avatar,
                nickname: this.userInfo.nickname, email:this.userInfo.email,
                phone: this.userInfo.phone, sex: parseInt(this.userInfo.sex)};
            api.updateUserInfo(this.userInfo.id, data).then(res => {
                this.$Message.success("更新个人资料成功!");
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        updatePassword() {
            api.updatePassword(this.user.id, this.password).then(res => {
                this.$Message.success("更新密码成功!");
            }, err => {
                this.$Message.error(err.response.data.message);
            });
        },
        uploadSuccess (res, file) {
            debugger;
            this.$store.dispatch('IsLoading', false);

            let imgURL = this.fileUploadToken.url_prefix + "/" + res.key;
            this.userInfo.avatar = imgURL;
            this.$Message.success("上传头像成功, 点击[提交]后将更新到系统.");
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
        handleBeforeUpload () {
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
        logout() {
            this.$store.dispatch('UserLogout');
            this.$router.push({name: 'index'});
        }
    },  
    beforeCreate() {
        this.$store.commit(types.SHOW_NAV_TABLES, false);
    },
}
</script>
<style scoped>
.container { padding:30px;}
.ivu-upload {
    background: #fff;
    border-radius: 4px;
    text-align: center;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: border-color .2s ease;
}
.avatar-container {
    text-align:center;
    padding-bottom: 30px;
}
.upload-icon {
    width: 100px;
    height:100px;
    background: rgba(0,0,0,0.6);
    position:absolute;
    top: 0;
    left: 0;
    z-index:99;
    text-align:center;
    display:none;
}
.avatar {
    text-align:center;
}

.avatar:hover > .upload-icon {
    display:block;
}
.avatar i {
    z-index: 99;
    color: white;
    font-size: 20px;
    vertical-align:middle;
}

</style>