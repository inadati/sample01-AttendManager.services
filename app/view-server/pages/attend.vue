<template>
    <div class="attend">
        <h1>出勤管理</h1>
        <v-container>
            <v-row class="px-3 mb-5">
                <v-btn v-if="!isThisWeek" class="pl-3" color="primary" @click="isThisWeek = true">
                    <v-icon small>keyboard_arrow_left</v-icon>前の週
                </v-btn>
                <v-spacer></v-spacer>
                <v-btn v-if="isThisWeek" class="pr-3" color="primary" @click="isThisWeek = false">
                    次の週
                    <v-icon small>keyboard_arrow_right</v-icon>
                </v-btn>
            </v-row>
            <v-card class="mb-5" v-for="(staff, istaff) in staffs" :key="istaff">
                <v-row>
                    <v-col class="pr-0" md="2">
                        <div class="pa-4">
                            <img src="/img/staff.svg" />
                        </div>
                        <v-divider></v-divider>
                        <div class="pa-3 pb-0">{{staff.name}}</div>
                    </v-col>
                    <v-divider vertical></v-divider>
                    <v-col class="pl-0 py-0">
                        <v-row class="container--fluid fill-height ml-0">
                            <v-col
                                class="px-0 py-1 attend_list"
                                v-for="(attend, iattend) in staff.attends"
                                :key="iattend"
                                v-if="iattend >= startDateIndex && iattend <= endDateIndex"
                            >
                                <div>{{dateStringJp(attend.date_start_time)}}</div>
                                <v-divider></v-divider>
                                <v-container fluid class="fill-height mt-n5">
                                    <v-row class="ml-0 pt-5" justify="center" align="center">
                                        <v-btn
                                            v-if="!attend.is_attend"
                                            fab
                                            color="primary"
                                            @click="setAttend(attend)"
                                        >
                                            <v-icon>add</v-icon>
                                        </v-btn>
                                        <div v-else>
                                            <v-select
                                                solo
                                                menu-props="auto"
                                                :items="timeTable"
                                                v-model="attend.in_time"
                                                @change="updateAttend(attend)"
                                            ></v-select>
                                            <div class="mt-n5 mb-2">~</div>
                                            <v-select
                                                solo
                                                menu-props="auto"
                                                :items="timeTable"
                                                v-model="attend.out_time"
                                                @change="updateAttend(attend)"
                                            ></v-select>
                                            <v-row justify="center">
                                                <v-btn dark color="red" @click="reset(attend)">削除</v-btn>
                                            </v-row>
                                        </div>
                                    </v-row>
                                </v-container>
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-card>
        </v-container>
    </div>
</template>

<script>
import { TimeTableDataProvider } from '~/servant/TimeTableDataProvider'
import staffsGql from '~/apollo/queries/staffsForAttend'
import updateAttendGql from '~/apollo/mutations/updateAttend'

export default {
    middleware: 'authenticated',
    data() {
        const ttdp = new TimeTableDataProvider.Summon()
        return {
            isThisWeek: true,
            startDateIndex: 0,
            endDateIndex: 6,
            staffs: [],
            timeTable: ttdp.serveData(),
        }
    },
    apollo: {
        staffs: {
            query: staffsGql,
        },
    },
    methods: {
        async reset(attend) {
            await this.$apollo.mutate({
                mutation: updateAttendGql,
                variables: {
                    id: attend.id,
                    is_attend: false,
                    in_time: '9:00',
                    out_time: '17:00',
                },
                refetchQueries: [
                    {
                        query: staffsGql,
                    },
                ],
            })
            attend.is_attend = false
            attend.in_time = '9:00'
            attend.out_time = '17:00'
        },
        async setAttend(attend) {
            attend.is_attend = true
            this.updateAttend(attend)
        },
        async updateAttend(attend) {
            await this.$apollo.mutate({
                mutation: updateAttendGql,
                variables: {
                    id: attend.id,
                    is_attend: attend.is_attend,
                    in_time: attend.in_time,
                    out_time: attend.out_time
                },
                refetchQueries: [
                    {
                        query: staffsGql,
                    },
                ],
            })
        },
        dateStringJp(iso8601) {
            const d = new Date(iso8601)
            let dofw = d.getDay()
            const dofwStr = (dofw = ['日', '月', '火', '水', '木', '金', '土'][
                dofw
            ])
            return `${d.getMonth() + 1}/${d.getDate()}(${dofwStr})`
        },
    },
    computed: {},
    watch: {
        isThisWeek(isThisWeek) {
            if (!isThisWeek) {
                this.startDateIndex = 7
                this.endDateIndex = 13
            } else {
                this.startDateIndex = 0
                this.endDateIndex = 6
            }
        },
    },
}
</script>

<style lang="stylus" scoped>
.attend {
  .attend_list {
    border-right: 1px solid rgba(0, 0, 0, 0.12);

    &:last-of-type, &:nth-of-type(7) {
      border-right: none;
    }
  }
}
</style>
