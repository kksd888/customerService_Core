<template>
    <div v-if="customerId" class="ChatWindowCon" id="ChatWindowCon">
        <h2 class="ChatName">
            {{ customerName }}
            <span>ID: {{ customerId }}</span>
        </h2>

        <!--聊天面板-->
        <ul class="ChatWindowDiv" id="ChatWindow">
            <li v-for="item in ContactItem"
                :class="[item.oper_code === 2003 ? 'MeCon' : 'FriendsCon']"
                :ids="item.id">
                <p class="ChatWindowTime">{{ item.create_time.substring(0, 19).replace('T', ' ') }}</p>
                <div class="ChatWindowIn">
                    <img src="../assets/images/logo.png" alt="img" class="touxiang"
                         v-if="item.oper_code === 2003">
                    <img v-else :src="customerImg? customerImg : '../assets/images/default.png' " alt="img"
                         class="touxiang">
                    <div class="ChatWindowText" v-html="item.msg.replace(/\n/g, '<br>')"></div>
                </div>
            </li>
        </ul>

        <!--输入框-->
        <div class="InputBox clear">
            <div class="InputBoxIcon">
                <!-- 转接 -->
                <el-tooltip class="item" effect="dark" content="将此用户转接给其他客服" placement="top">
                    <el-popover
                            placement="right"
                            width="300"
                            trigger="click"
                            v-model="TransferPopover"
                            @show="TransferShow"
                            @after-leave="TransferValue=''">
                        <el-select v-model="TransferValue" placeholder="请选择">
                            <el-option-group
                                    v-for="group in TransferOptions"
                                    :key="group.label"
                                    :label="group.label">
                                <el-option
                                        v-for="item in group.options"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                </el-option>
                            </el-option-group>
                        </el-select>
                        <el-button type="danger" size="mini" style="margin-left: 20px;" @click="Transfer">
                            确定转接
                        </el-button>
                        <i class="el-icon-share icon_Upload" slot="reference"></i>
                    </el-popover>
                </el-tooltip>
            </div>
            <textarea
                    class="textarea"
                    placeholder="请输入内容
按Enter发送 Ctrl+Enter换行"
                    @keyup.enter.stop="sendMesage"
                    v-model="Msg">
            </textarea>
        </div>

        <!--右侧智能面板-->
        <div class="rightCon">
            <div class="ChatOffLineDiv">
                <el-tabs v-model="ActiveName">
                    <el-tab-pane label="智能回复" name="SmartReply">
                        <span slot="label"><i class="el-icon-date"></i> 智能回复</span>
                        <div v-if="AnswerText" class="ChatOffLineHide">
                            <div class="ChatOffLineDivCon">
                                <div class="ChatOffLineDivConWen">
                                    <span>问：</span>
                                    <p class="AnswerTitle" v-html="AnswerTitle.replace(/\n/g, '<br >')">
                                    </p>
                                </div>
                                <div class="ChatOffLineDivConWen">
                                    <span>答：</span>
                                    <pre class="AnswerText" v-html="AnswerText.replace(/\n/g, '<br >')">
                                    </pre>
                                </div>
                            </div>
                            <div class="ChatOffLineInput">
                                <el-button type="primary" size="small" plain @click="Msg = AnswerText">将回答放入输入框
                                </el-button>
                            </div>
                        </div>
                        <div style="text-align: center;" v-else>
                            <small>小金机器人未侦测到可回复问题</small>
                        </div>
                    </el-tab-pane>
                    <el-tab-pane label="用户画像" name="UserPortrait">
                        <span slot="label"><i class="el-icon-picture"></i> 用户画像</span>
                        <div style="text-align: center;">
                            <small>程序猿小哥哥加紧开发中...</small>
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </div>

    </div>
</template>

<script>
    import {formatDate} from '@/util/date'

    export default {
        name: "ChatBody",
        data() {
            return {
                customerName: "",
                customerImg: "",
                Msg: "",
                ActiveName: "SmartReply",
                AnswerTitle: "",
                AnswerText: "",
                TransferValue: '',
                TransferOptions: [],
                TransferPopover: false,
            };
        },
        computed: {
            customerId() {
                this.Msg = "";
                this.customerName = this.$store.state.currentCustomer.customer_nick_name;
                this.customerImg = this.$store.state.currentCustomer.customer_head_img_url;
                return this.$store.state.currentCustomer.customer_id;
            },
            ContactItem() {
                let megs = this.$store.state.currentRoomMessages;
                this.AnswerTitle = megs[megs.length - 1].msg;
                this.AnswerText = megs[megs.length - 1].ai_msg;
                return megs;
            }
        },
        methods: {
            clearStatus() {
                this.Msg = "";
                this.AnswerTitle = "";
                this.AnswerText = "";
            },
            /**
             * 发送消息
             */
            sendMesage(e) {
                if (e.ctrlKey) {
                    this.Msg += '\n';
                    return
                }

                let sendMsg = this.Msg;
                this.clearStatus();
                if (!sendMsg || sendMsg.replace(/\n/g, '').replace(/\s/g, '').length === 0) {
                    this.$message.error('不能发送空的内容');
                    return
                }
                this.axios.post("/admin/dialog", {
                    msg_type: "text",
                    customer_id: this.customerId,
                    msg: sendMsg,
                }).then(resp => {
                    this.ContactItem.push({
                        create_time: formatDate(new Date(), 'yyyy-MM-dd hh:mm:ss'),
                        msg: sendMsg,
                        oper_code: 2003,
                    });
                }, err => {
                    this.$message.error('发送失败，微信限制给固定用户发送的频率和时间，需要等待用户再次主动发起聊天才能够发送');
                    console.log(err)
                });
            }
            ,
            /**
             * 显示转接select
             * @constructor
             */
            TransferShow() {
                this.axios.get('admin/kf/online')
                    .then(resp => {
                        this.TransferOptions = resp.data;
                    });
            }
            ,
            /**
             * 将当前用户转接给其他客服
             * @constructor
             */
            Transfer() {
                let transferKfId = this.TransferValue;
                let customerId = this.customerId;

                this.axios.post('/admin/room/transfer', {
                    customer_id: customerId,
                    transfer_kf_id: transferKfId,
                }).then(resp => {
                    this.TransferValue = "";
                    this.TransferPopover = false;
                    this.$store.commit('transferRoom', this.customerId);
                    this.$store.commit('cleanCurrentCustomer');
                    this.$notify.success({
                        title: "客户转接",
                        message: `转接成功`,
                    });
                });
            }
            ,
        }
        ,
        updated: function () {
            var container = document.getElementById("ChatWindow");
            if (container) {
                container.scrollTop = container.scrollHeight;
            }
        }
    }
</script>

<style scoped>
    .textarea {
        font-size: 12px;
    }

    .ChatName span {
        font-size: small;
        color: #DADADA;
    }
</style>