<template>
    <div class="ChatWrapItem" id="ChatWrapReply" style="display: block;">
        <div class="ContactsList">
            <el-row v-if="ContactsItems.length<1" class="none_prompt">
                ~~空空如也~~
            </el-row>
            <ul v-else class="ContactsList_ul" id="dialogList">
                <li v-for="item in ContactsItems" :type="item.room_customer.customer_id"
                    @click="openChatBody(item.room_customer)"
                    :class="[CustomerId === item.room_customer.customer_id ? 'dialog_select':'']">
                    <dl class="ContactsDl">
                        <dt>
                            <img v-if="item.room_customer.customer_head_img_url"
                                 :src="item.room_customer.customer_head_img_url"
                                 :alt="item.room_customer.customer_nick_name"
                                 class="ContactsImg">
                            <img v-else src="../assets/images/default.png" alt="暂无头像" class="ContactsImg">
                        </dt>
                        <dd>
                            <span class="ContactsName">
                            {{ item.room_customer.customer_nick_name }}
                            </span>
                            <div class="ChatCon">
                                {{ item.room_messages[item.room_messages.length - 1].msg }}
                            </div>
                        </dd>
                    </dl>
                    <p class="ChatTime">
                        {{ item.room_messages[item.room_messages.length - 1].create_time.substring(0, 19).replace('T', ' ') }}
                    </p>
                    <em class="UnreadEm"
                        v-if="item.room_messages[item.room_messages.length - 1].oper_code === 2002 && item.room_messages[item.room_messages.length - 1].ack === false">
                        未读
                    </em>
                </li>
            </ul>
        </div>
        <ChatBody ref="chatB"></ChatBody>
    </div>
</template>

<script>
    import ChatBody from "@/components/ChatBody";

    export default {
        name: "ReplyPanel",
        components: {ChatBody},
        data() {
            return {
                CustomerId: "",
            }
        },
        computed: {
            ContactsItems() {
                return this.$store.state.roomData;
            },
        },
        methods: {
            /**
             * 打开聊天窗口
             * @param customerInfo 聊天的用户信息
             */
            openChatBody: function (customerInfo) {
                this.CustomerId = customerInfo.customer_id;
                // 初始化当前聊天人信息
                this.$store.commit('setCurrentCustomer', customerInfo);
                // 确认已读
                this.axios.put('/admin/dialog/ack', {customer_ids: [customerInfo.customer_id]});
                // 移除已读消息
                this.$store.commit('setWaitAck', customerInfo.customer_id);
            },
        },
    }
</script>

<style scoped>
    .dialog_select {
        background-color: whitesmoke;
    }

    .none_prompt {
        text-align: center;
        font-size: 12px;
        color: #B4B4B4;
        line-height: 30px;
    }

    .ChatWrapItem {
        display: none;
    }

    .ContactsList {
        width: 189px;
        background: #fff;
        border-right: 1px solid #DEDEDE;
        overflow: auto;
        -webkit-overflow-scrolling: touch;
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        z-index: 0;
    }

    /*滚动条样式*/
    .ContactsList::-webkit-scrollbar {
        display: none;
    }

    .ContactsList::-webkit-scrollbar {
        width: 4px;
        height: 4px;
    }

    /*滚动条整体样式*/
    .ContactsList::-webkit-scrollbar-thumb {
        border-radius: 5px;
        -webkit-box-shadow: inset 0 0 5px #ccc9c7;
        background: #ccc9c7;
    }

    /*滚动条里面小方块*/
    .ContactsList::-webkit-scrollbar-track {
        -webkit-box-shadow: inset 0 0 5px #ccc9c7;
        border-radius: 0;
        background: none;
    }

    /*滚动条里面轨道*/
    .ContactsList_ul li {
        padding: 10px 15px;
        clear: both;
        cursor: pointer;
        position: relative;
    }

    .ContactsList_ul li:hover {
        background: #f9f9f9;
    }

    .ContactsList_ul .ContactsListCur {
        background: #f1f1f1;
    }

    .NoCon {
        font-size: 12px;
        color: #C6C6C6;
        text-align: center;
        padding-top: 200px;
    }

    .ContactsDl {
        margin-right: 20px;
        position: relative;
    }

    .ContactsDl dt {
        position: absolute;
        top: 0;
        left: 0;
    }

    .ContactsImg {
        width: 30px;
        height: 30px;
    }

    .ContactsDl dd {
        height: 30px;
        margin-left: 40px;
    }

    .ContactsName {
        text-align: left;
        display: block;
        line-height: 1.2;
        font-size: 14px;
        color: #000;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }


    .ChatTime {
        padding-top: 10px;
        font-size: 12px;
        color: #C6C6C6;
        text-align: right;
    }

    .UnreadEm {
        font-size: 12px;
        color: #fff;
        background: #ED5103;
        padding: 2px 3px;
        position: absolute;
        top: 15px;
        right: 15px;
    }

    .ChatCon {
        display: block;
        height: 13px;
        font-size: 12px;
        text-align: left;
        color: #B4B4B4;
        margin-top: 3px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
</style>