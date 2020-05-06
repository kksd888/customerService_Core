<template>
    <el-row class="AccessMain" id="AccessMain">
        <div style="margin-top: 20px;">
            <el-button type="primary" size="small" :loading="loading" plain @click="search">查询</el-button>
            <el-date-picker
                    v-model="pickerValue"
                    :picker-options="pickerOptions"
                    size="small"
                    type="datetimerange"
                    style="margin-left: 10px;"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期">
            </el-date-picker>
        </div>

        <el-table
                v-loading="loading"
                element-loading-text="拼命加载中"
                element-loading-spinner="el-icon-loading"
                element-loading-background="rgba(255, 255, 255, 0.8)"

                height="570"
                :data="tableData"
                ref="dTable"
                stripe
                show-summary
                size="medium"
                style="width: 100%">
            <el-table-column
                    type="index">
            </el-table-column>
            <el-table-column
                    prop="JobNum"
                    label="工号"
                    width="180">
            </el-table-column>
            <el-table-column
                    prop="NickName"
                    label="姓名">
            </el-table-column>
            <el-table-column
                    prop="GroupName"
                    label="组名"
                    :filters="[{ text: '咨询组', value: '咨询组' }, { text: '投诉组', value: '投诉组' }]"
                    :filter-method="filterTag"
                    filter-placement="bottom-end">
                <template slot-scope="scope">
                    <el-tag
                            :type="scope.row.GroupName === '投诉组' ? 'primary' : 'success'"
                            disable-transitions>{{scope.row.GroupName}}
                    </el-tag>
                </template>
            </el-table-column>
            <el-table-column
                    prop="SendCount"
                    sortable
                    label="发送消息总量">
            </el-table-column>
            <el-table-column
                    prop="ReceptionCount"
                    sortable
                    label="服务客户总量">
            </el-table-column>
        </el-table>
    </el-row>
</template>

<script>
    export default {
        name: "StatisticsPanel",
        data() {
            return {
                loading: false,
                pickerValue: [],
                pickerOptions: {
                    shortcuts: [{
                        text: '最近一周',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
                            picker.$emit('pick', [start, end]);
                        }
                    }, {
                        text: '最近一个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
                            picker.$emit('pick', [start, end]);
                        }
                    }, {
                        text: '最近三个月',
                        onClick(picker) {
                            const end = new Date();
                            const start = new Date();
                            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
                            picker.$emit('pick', [start, end]);
                        }
                    }]
                },
                tableData: []
            }
        },
        methods: {
            filterTag(value, row) {
                return row.GroupName === value;
            },
            search() {
                if (this.pickerValue.length !== 2) {
                    this.$message.error("查询时间不正确");
                    return
                }

                this.loading = true;
                this.axios.post('/admin/statistics/', {
                    "StartTime": this.pickerValue[0],
                    "EndTime": this.pickerValue[1],
                }).then(resp => {
                    this.tableData = resp.data;
                }).finally(() => {
                    this.loading = false;
                });
            },
        },
    }
</script>

<style scoped>

</style>