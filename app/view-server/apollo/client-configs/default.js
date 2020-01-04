import { createHttpLink } from 'apollo-link-http'
import { setContext } from 'apollo-link-context'
import { InMemoryCache } from 'apollo-cache-inmemory'
import CryptoJS from 'crypto-js'

export default () => {
    const httpLink = createHttpLink({
        uri: process.env.GRAPHQL_ENTRY_POINT,
        fetchOptions: {
            accept: '*/*',
            'content-type': 'application/json',
        },
    })

    const authLink = setContext((_, { headers }) => {
        let authorization = ''
        const encodeToken = localStorage.getItem('token')
        if (encodeToken) {
            const payload = CryptoJS.AES.decrypt(
                encodeToken,
                process.env.SECRET_KEY
            )
            const token = payload.toString(CryptoJS.enc.Utf8).replace(/\"/g, '')
            authorization = `Bearer ${token}`
        }
        return {
            headers: {
                ...headers,
                authorization,
            },
        }
    })

    return {
        link: authLink.concat(httpLink),
        cache: new InMemoryCache(),
    }
}
