async function loginUser(username, password) {
    const data = { name: username, password: password };
    try {
       const res = await fetch('http://localhost:8080/api/user/login', { method: 'POST', body: JSON.stringify(data) })
        const jsonResponse = await res.json();
        console.log(jsonResponse);
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err;
    }
}

function checkUser() {
    try {
        const res = fetch('http://localhost:8080/api/user');
        if (!res.ok) {
            return res.status
        } else {
            return 'Success!'
        }
    } catch (error) {
        console.log(error);
        return error;
    }
}
export { loginUser, checkUser };