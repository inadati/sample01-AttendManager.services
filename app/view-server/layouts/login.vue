<template>
    <v-app id="inspire">
        <v-content>
            <v-container class="fill-height" fluid>
                <v-row align="center" justify="center">
                    <nuxt />
                    <v-col cols="12" sm="8" md="4">
                        <v-card class="elevation-12">
                            <v-toolbar color="primary" dark flat>
                                <v-toolbar-title>Login form</v-toolbar-title>
                            </v-toolbar>
                            <v-card-text>
                                <v-form
                                    ref="loginForm"
                                    v-model="loginForm.valid"
                                >
                                    <v-text-field
                                        v-model="loginForm.email.model"
                                        :rules="loginForm.email.rules"
                                        label="Login"
                                        name="login"
                                        prepend-icon="person"
                                        type="text"
                                    />

                                    <v-text-field
                                        v-model="loginForm.password.model"
                                        :rules="loginForm.password.rules"
                                        id="password"
                                        label="Password"
                                        name="password"
                                        prepend-icon="lock"
                                        type="password"
                                    />
                                </v-form>
                            </v-card-text>
                            <v-card-actions>
                                <v-spacer />
                                <v-btn color="primary" @click="login">
                                    Login
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </v-content>
    </v-app>
</template>

<script>
import getCertificateGql from '~/apollo/queries/getCertificate.gql'
import CryptoJS from 'crypto-js'
const Cookie = process.client ? require('js-cookie') : undefined

export default {
    data() {
        return {
            loginForm: {
                valid: true,
                email: {
                    model: '',
                    rules: [
                        v => !!v || 'メールアドレスが未入力です。',
                        v =>
                            /.+@.+\..+/.test(v) ||
                            '正しい形式のメールアドレスを入力してください。',
                    ],
                },
                password: {
                    model: '',
                    rules: [v => !!v || 'パスワードが未入力です。'],
                },
            },
        }
    },
    methods: {
        login() {
            if (!this.$refs.loginForm.validate()) {
                this.snackbar = true
                return
            }

            setTimeout(async () => {
                const res = await this.$apolloProvider.defaultClient.query({
                    query: getCertificateGql,
                    variables: {
                        email: this.loginForm.email.model,
                        password: this.loginForm.password.model,
                    },
                })
                console.dir(res)
                if (res.data.getCertificate.token) {
                    const certificate = res.data.getCertificate
                    this.$store.commit(
                        'mutateLoggedInAccount',
                        certificate.user_name
                    )
                    Cookie.set('loggedInAccount', certificate.user_name, {
                        expires: 3,
                    })


                    localStorage.setItem(
                        'token',
                        CryptoJS.AES.encrypt(
                            JSON.stringify(certificate.token),
                            process.env.SECRET_KEY
                        ).toString()
                    )
                    this.$router.push('/')
                }
            }, 1000)
        },
    },
    computed: {
        name() {
            return this.data
        },
    },
}
</script>
