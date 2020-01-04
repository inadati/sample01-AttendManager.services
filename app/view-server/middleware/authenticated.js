export default function ({ store, redirect }) {
    // ユーザが認証されていない場合
    if (!store.state.loggedInAccount) {
        return redirect('/login')
    }
}