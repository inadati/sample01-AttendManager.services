const cookieparser = require('cookieparser')

export default {
    async nuxtClientInit({ commit, state }, context) {
        let loggedInAccount = null
        const parsed = cookieparser.parse(document.cookie)
        try {
            loggedInAccount = parsed.loggedInAccount
        } catch (err) {
            localStorage.removeItem('token')
        }
        commit('mutateLoggedInAccount', loggedInAccount)
    },
}
