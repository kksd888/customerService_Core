<template>
    <div class="login">
        <el-row>
            <img src="../assets/images/login_logo.png" alt="金色世纪" class="login_logo"/>
        </el-row>
        <el-row>
            <el-input
                    class="input_text"
                    placeholder="请输入"
                    v-model="login_data.job_num"
                    prefix-icon="el-icon-service"
                    clearable>
                <template slot="prepend">工号</template>
            </el-input>
        </el-row>
        <el-row>
            <el-input
                    class="input_text"
                    placeholder="请输入密码"
                    v-model="login_data.pass_word"
                    prefix-icon="el-icon-bell"
                    show-password>
                <template slot="prepend">密码</template>
            </el-input>
        </el-row>
        <el-row>
            <el-radio class="input_radio" v-model="login_data.group_name" label="咨询组">咨询组</el-radio>
            <el-radio class="input_radio" v-model="login_data.group_name" label="投诉组">投诉组</el-radio>
        </el-row>
        <el-row>
            <el-button class="input_button"
                       icon="el-icon-arrow-left"
                       ref="l_button"
                       type="primary"
                       @click="login"
                       :loading="loading"
                       plain>
                登录
                <i class="el-icon-arrow-right el-icon--right"></i>
            </el-button>
        </el-row>
    </div>
</template>

<script>
    export default {
        name: "Login",
        data() {
            return {
                login_data: {
                    job_num: "",
                    pass_word: "",
                    group_name: "咨询组",
                },
                loading: false
            };
        },
        methods: {
            login() {
                let that = this;
                that.loading = true;
                that.axios.post('/admin/login', that.login_data)
                    .then(resp => {
                        this.$cookies.set('token', resp.data['Authentication'], '7d');
                        this.$cookies.set('group_name', resp.data['GroupName'], '7d');
                        this.$cookies.set('job_num', resp.data['JobNum'], '7d');
                        this.$cookies.set('nick_name', resp.data['NickName'], '7d');
                        this.$router.push("/home");
                    }, err => {
                        console.log("登录发生错误", err);
                        this.$notify.error({
                            title: '登录错误',
                            message: `${err.message}, 账号或密码错误！`
                        });
                    })
                    .finally(() => {
                        that.loading = false;
                    });
            }
        },
        created() {
            this.$store.commit('cleanwSocket');
        }
    }
</script>

<style scoped>
    .login {
        position: fixed;
        top: 50%;
        left: 50%;
        width: 483px;
        background: #F5F5F5;
        padding: 70px 0;
        margin: -220px 0 0 -241px;
        border-radius: 20px;
        box-shadow: 0 6px 22px 0 rgba(0, 0, 0, 0.1);
        text-align: center;
    }

    .login_logo {
        display: block;
        width: 240px;
        margin: 0 auto 50px;
    }

    .input_text {
        width: 300px;
        margin-bottom: 5px;
    }

    .input_radio {
        margin: 10px;
    }

    .input_button {
        margin-top: 10px;
        width: 200px;
        letter-spacing: 4px;
    }
</style>