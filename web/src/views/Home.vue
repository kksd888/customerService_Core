<template>
    <div class="ChatWindow">
        <div class="ChatHead">
            <ul class="ChatHead_lt">
                <li :class="[ currentTab === 'ReplyPanel' ? 'ChatHeadCur' : '' ]"
                    @click="currentTab='ReplyPanel'">
                    <span class="access">
                        待回复
                        <em v-if="replyNum>0">{{ replyNum }}</em>
                    </span>
                </li>
                <li :class="[ currentTab === 'AccessPanel' ? 'ChatHeadCur' : '' ]"
                    @click="currentTab='AccessPanel'">
                    <span class="access">
                        待接入
                        <em v-if="accessNum>0">{{ accessNum }}</em>
                    </span>
                </li>
                <li :class="[ currentTab === 'StatisticsPanel' ? 'ChatHeadCur' : '' ]"
                    @click="currentTab='StatisticsPanel'">
                    <span class="access">
                        统计数据
                    </span>
                </li>
            </ul>
            <div class="ChatHead_rt">
                <p class="ChatHead_btn">
                    <a class="OffLineBtn" @click="loginOut" href="javascript:">退出</a>
                </p>
            </div>
        </div>

        <div class="ChatWindowMain">
            <div class="ChatWrap clear">
                <keep-alive>
                    <component :is="currentTabComputed">
                    </component>
                </keep-alive>
            </div>
        </div>

        <div class="ChatHead_status">
            <div :class="['ChatHead_status_icon', isOnline ? 'ChatHead_status_icon_green' : 'ChatHead_status_icon_red']"></div>
            <span>{{ isOnline ? "在线" : "离线" }}</span>
        </div>

    </div>
</template>

<script>
    import ReplyPanel from "@/components/ReplyPanel";
    import AccessPanel from "@/components/AccessPanel";
    import StatisticsPanel from "@/components/StatisticsPanel";

    export default {
        components: {AccessPanel, ReplyPanel, StatisticsPanel},
        data() {
            return {
                currentTab: 'ReplyPanel',
                wsHeartBeat: null,
                isOnline: false,
            }
        },
        computed: {
            currentTabComputed() {
                return this.currentTab
            },
            replyNum() {
                return this.$store.state.roomData.filter(x => x.room_messages[x.room_messages.length - 1].oper_code === 2002 && x.room_messages[x.room_messages.length - 1].ack === false).length;
            },
            accessNum() {
                return this.$store.state.waitData.length;
            }
        },
        methods: {
            initWSocket() {
                let wsUrl = `${this.axios.wsHost}/admin/ws?token=${encodeURIComponent(this.$cookies.get("token"))}`;

                var openWs = function (that) {
                    var ws = new WebSocket(wsUrl);
                    ws.onopen = function (evn) {
                        that.$notify.success({
                            title: "WS协议连接成功",
                            message: `【${that.$cookies.get('group_name')}】客服代表【${that.$cookies.get('nick_name')}】欢迎登录`,
                        });
                        that.isOnline = true;
                    };
                    ws.onmessage = function (evn) {
                        that.isOnline = true;
                        let responseData = JSON.parse(evn.data);
                        // console.log('收到websocket消息', responseData);

                        switch (responseData.Type) {
                            case 0:
                                // 心跳
                                break;
                            case 1:
                                // 新消息
                                that.axios.get(`/admin/room/${responseData.Body}`)
                                    .then(resp => {
                                        let roomData = resp.data;
                                        if (roomData.room_kf.kf_id) {
                                            // 新消息
                                            console.log("新消息", roomData);
                                            that.$store.commit("updateRoomData", roomData);
                                        } else {
                                            // 待接入用户
                                            console.log("待接入用户", roomData);
                                            that.$store.commit("updateWaitData", roomData);
                                        }
                                    });
                                break;
                            case 2:
                                // 广播关闭已经被接入的客户
                                console.log("收到广播", responseData.Body);
                                that.$store.commit('cleanWaitData', [responseData.Body]);
                                break;
                            default:
                                console.log("未知的Socket指令", responseData);
                                break;
                        }
                    };

                    ws.onclose = function (evn) {
                        that.$notify.error({
                            title: "WS连接异常",
                            message: "服务器已断开连接，正在尝试重连"
                        });
                        that.isOnline = false;
                    };
                    return ws;
                };
                this.$store.commit('initwSocket', openWs(this));

                this.wsHeartBeat = setInterval(() => {
                    let ws = this.$store.state.wSocket;
                    switch (ws.readyState) {
                        case 0:
                            break;
                        case 1:
                            this.$store.state.wSocket.send("+");
                            break;
                        case 2 | 3:
                            this.$store.commit('cleanwSocket');
                            this.$store.commit('initwSocket', openWs(this));
                            break;
                    }
                }, 5000);
            },
            loginOut
                () {
                this.$cookies.remove('token');
                this.$notify.success({
                    title: '安全退出',
                    message: `
                            您已经安全退出！`
                });
                this.$router.push("/login");
            }
            ,
        },
        created() {
            this.$notify.info({
                title: "连接服务器",
                message: "数据初始化中...",
            });

            // 初始化WS协议连接
            this.initWSocket();

            // 初始化主面板数据
            this.axios.get('/admin/init')
                .then(resp => {
                    this.$store.commit('setRoomData', resp.data.online_customer);
                    this.$store.commit('setWaitData', resp.data.wait_customer);
                });
        },
        destroyed() {
            clearInterval(this.wsHeartBeat);
        },
    }
</script>

<style scoped>
    @import "../assets/home.css";

    .ChatHead_status {
        position: fixed;
        top: 10px;
        left: 10px;
        line-height: 12px;
    }

    .ChatHead_status span {
        margin-left: 5px;
        font-size: 12px;
        color: #666666;
    }

    .ChatHead_status_icon {
        width: 10px;
        height: 10px;
        -moz-border-radius: 5px;
        -webkit-border-radius: 5px;
        border-radius: 5px;
        float: left;
    }

    .ChatHead_status_icon_red {
        background: red;
    }

    .ChatHead_status_icon_green {
        background: #6BD66A;
    }
</style>