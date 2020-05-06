import Vue from 'vue'
import Vuex from 'vuex'
import {Notification} from 'element-ui'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        /**
         * 当前的用户
         */
        currentCustomer: {
            customer_id: "",
            customer_nick_name: "",
            customer_head_img_url: "",
        },
        /**
         * 当前会话
         */
        currentRoomMessages: [
            {
                id: "",
                kf_id: "",
                type: "",
                media_url: "",
                msg: "",
                ai_msg: "",
                ack: false,
                oper_code: 0,
                create_time: "",
            }
        ],
        /**
         * 进行中的聊天
         */
        roomData: [
            {
                room_customer: {
                    customer_id: "",
                    customer_nick_name: "",
                    customer_head_img_url: "",
                },
                room_messages: [
                    {
                        id: "",
                        kf_id: "",
                        type: "",
                        media_url: "",
                        msg: "",
                        ai_msg: "",
                        ack: false,
                        oper_code: 0,
                        create_time: "",
                    }
                ],
            }
        ],
        /**
         * 待接入的聊天
         */
        waitData: [
            {
                room_customer: {
                    customer_id: "",
                    customer_nick_name: "",
                    customer_head_img_url: "",
                },
                room_messages: [
                    {
                        id: "",
                        kf_id: "",
                        type: "",
                        media_url: "",
                        msg: "",
                        ai_msg: "",
                        ack: false,
                        oper_code: 0,
                        create_time: "",
                    }
                ],
            }
        ],
        /**
         * 后台WebSocket连接
         */
        wSocket: null,
    },
    mutations: {
        setCurrentCustomer(state, payload) {
            state.currentCustomer = payload;
            for (let i = 0; i < state.roomData.length; i++) {
                if (payload.customer_id === state.roomData[i].room_customer.customer_id) {
                    state.currentRoomMessages = state.roomData[i].room_messages;
                    return
                }
            }
        },
        setRoomData(state, payload) {
            state.roomData = payload;
        },
        setWaitData(state, payload) {
            state.waitData = payload;
        },
        updateRoomData(state, payload) {
            // 更新roomData
            let isFind = false;
            for (let i = 0; i < state.roomData.length; i++) {
                if (payload.room_customer.customer_id === state.roomData[i].room_customer.customer_id) {
                    state.roomData[i] = payload;
                    isFind = true;
                }
            }
            if (!isFind) {
                // 置于聊天列表第一位
                state.roomData = [payload].concat(state.roomData);
                Notification.warning({
                    title: "新客户接入",
                    message: "聊天列表加入一位新客户，请注意查看",
                });
            }

            if (payload.room_customer.customer_id === state.currentCustomer.customer_id) {
                state.currentRoomMessages = payload.room_messages;
            }

            // 更新
            state.roomData = [].concat(state.roomData);
        },
        updateWaitData(state, payload) {
            let isFind = false;
            state.waitData.forEach(r => {
                if (payload.room_customer.customer_id === r.room_customer.customer_id) {
                    r = payload;
                }
            });
        },
        cleanWaitData(state, ids) {
            if (ids) {
                state.waitData = state.waitData.filter(x => !ids.includes(x.room_customer.customer_id))
            } else {
                state.waitData = [];
            }
        },
        cleanCurrentCustomer(state) {
            state.currentCustomer = {};
            state.currentRoomMessages = [];
        },
        setWaitAck(state, customer_id) {
            state.roomData.forEach(d => {
                if (d.room_customer.customer_id === customer_id) {
                    d.room_messages[d.room_messages.length - 1].ack = true;
                }
            });

        },
        transferRoom(state, customer_id) {
            for (let i = 0; i < state.roomData.length; i++) {
                if (state.roomData[i].room_customer.customer_id === customer_id) {
                    state.roomData.splice(i, 1);
                    return
                }
            }
        },
        initwSocket(state, ws) {
            if (ws) {
                state.wSocket = ws;
            }
        },
        cleanwSocket(state) {
            if (state.wSocket && state.wSocket.readyState <= 1) {
                state.wSocket.onclose = function (evn) {
                };
                state.wSocket.close();
            }
            state.wSocket = null;
        }
    },
    actions: {}
})
