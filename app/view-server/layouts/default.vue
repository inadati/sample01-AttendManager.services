<template>
    <v-app id="inspire">
        <v-navigation-drawer v-model="drawer" app>
            <v-list dense>
                <v-list-item to="/">
                    <v-list-item-action>
                        <v-icon>mdi-home</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>スタッフ</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item to="/attend">
                    <v-list-item-action>
                        <v-icon>mdi-contact-mail</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>出勤管理</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar app color="indigo" dark>
            <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
            <v-toolbar-title>Application</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="logout">ログアウト</v-btn>
        </v-app-bar>

        <v-content>
            <v-container class="text-center" fluid>
                <nuxt />
            </v-container>
        </v-content>
    </v-app>
</template>

<script>
const Cookie = process.client ? require('js-cookie') : undefined
export default {
    props: {
        source: String,
    },

    data: () => ({
        drawer: null,
    }),
    methods: {
        logout() {
            Cookie.remove('loggedInAccount')
            this.$store.commit('mutateLoggedInAccount', null)
            localStorage.removeItem('token')
            this.$router.push('/login')
        }
    },
}
</script>
