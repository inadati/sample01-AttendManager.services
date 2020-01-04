export default function ({ store, redirect }) {
    // ユーザが認証されてホームページにリダイレクトされた場合
    if (store.state.loggedInAccount) {
        return redirect('/')
    }
}