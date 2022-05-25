async function loginUser(username, password) {
    const data = { name: username, password: password };
    try {
        const res = await fetch('/api/user/login', { method: 'POST', body: JSON.stringify(data) })
        const jsonResponse = await res.json();
        console.log(jsonResponse);
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err;
    }
}

async function getUser() {
    try {
        const res = await fetch('/api/user');
        if (!res.ok) {
            return res.status;
        } else {
        const jsonResponse = await res.json();
        console.log(jsonResponse)
        return jsonResponse;
        }
    } catch (error) {
        console.log(error);
        return error;
    }
}
export { loginUser, getUser };