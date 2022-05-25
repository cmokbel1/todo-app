// loginUser sends a requests to /api/user/login with the provided credentials.
// If an error occurs a string error message is returned otherwise a User object
// (with a name property) is returned.
async function loginUser(username, password) {
    const data = { name: username, password: password };
    try {
        const res = await fetch(
            '/api/user/login',
            {
                method: 'POST',
                body: JSON.stringify(data),
            },
        );
        const jsonResponse = await res.json();
        console.log(jsonResponse);
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}

// getUser returns either a User object (with a name property) or an error.
async function getUser() {
    try {
        const res = await fetch('/api/user');
        if (!res.ok) {
            if (res.status === 401) {
                return 'unauthorized'
            }
            return 'internal error';
        }
        const jsonResponse = await res.json();
        console.log(jsonResponse)
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}
export { loginUser, getUser };