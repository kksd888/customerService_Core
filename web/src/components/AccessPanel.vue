<template>
    <el-row class="AccessMain" id="AccessMain">
        <div style="margin-top: 20px">
            <el-button size="small" type="info" plain @click="toggleAllSelection">全选/取消选择</el-button>
            <el-button size="small" type="danger" plain @click="handleClick()">接入选中客户</el-button>
        </div>

        <el-table
                :data="AccessList"
                ref="dTable"
                stripe
                size="medium"
                style="width: 100%;"
                @selection-change="handleSelectionChange">
            <el-table-column
                    type="selection"
                    width="55">
            </el-table-column>
            <el-table-column
                    type="index"
                    :index="indexMethod">
            </el-table-column>
            <el-table-column
                    prop="customer_nick_name"
                    label="用户昵称"
                    width="280">
                <template slot-scope="scope">
                    <img :src="scope.row.customer_head_img_url" alt="img" class="touxiang"
                         v-if="scope.row.customer_head_img_url">
                    <img alt="img" class="touxiang"
                         v-else src="../assets/images/logo.png">
                    <span style="margin-left: 10px;">{{ scope.row.customer_nick_name }}</span>
                </template>
            </el-table-column>
            <el-table-column
                    prop="create_time"
                    label="日期"
                    width="180">
                <template slot-scope="scope">
                    <i class="el-icon-time"></i>
                    <span style="margin-left: 10px">{{ scope.row.create_time.substring(0, 19).replace('T', ' ') }}</span>
                </template>
            </el-table-column>
            <el-table-column
                    prop="msg"
                    label="消息"
                    width="400"
                    :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column
                    fixed="right"
                    label="操作"
                    width="100">
                <template slot-scope="scope">
                    <el-button
                            size="mini"
                            type="danger"
                            @click="handleClick(scope.row)">接入
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-row>
</template>

<script>
    export default {
        name: "AccessPanel",
        data() {
            return {
                indexMethod: 1,
                multipleSelection: [],
            }
        },
        computed: {
            AccessList() {
                let tDate = [];
                this.$store.state.waitData.forEach(w => {
                    tDate.push({
                        customer_id: w.room_customer.customer_id,
                        customer_nick_name: w.room_customer.customer_nick_name,
                        customer_head_img_url: w.room_customer.customer_head_img_url,
                        msg: w.room_messages[0].msg,
                        create_time: w.room_messages[0].create_time,
                    });
                });
                return tDate;
            }
        },
        methods: {
            handleSelectionChange(val) {
                this.multipleSelection = val;
            },
            toggleAllSelection() {
                this.$refs.dTable.toggleAllSelection();
            },
            handleClick(val) {
                if (val) {
                    this.accessCustomer([val.customer_id]);
                } else {
                    let customerIds = [];
                    this.multipleSelection.forEach(x => {
                        customerIds.push(x.customer_id);
                    });
                    this.accessCustomer(customerIds);
                }
            },
            accessCustomer(customer_ids) {
                this.axios.post('/admin/wait_queue/access', {
                    customer_ids: customer_ids,
                }).then(resp => {
                    this.$store.commit('cleanWaitData', customer_ids);
                });
            }
        }
    }
</script>

<style scoped>

</style>