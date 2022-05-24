async function loginUser(username, password) {
    const data = { name: username, password: password };
    try {
        await fetch('http://localhost:8080/api/user/login', { method: 'POST', body: JSON.stringify(data) })
            .then(response => response.json())
            .then(data => console.log(data));
    } catch (error) {
        console.log(error);
    }
}

async function checkUser() {
    try {
    await fetch('http://localhost:8080/api/user', {method: 'GET'}).then(response => response.json())
    } catch(error) {
        console.log(error)
    }
}
export { loginUser, checkUser };