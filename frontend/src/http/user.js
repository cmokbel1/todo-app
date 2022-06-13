//register a new user
async function registerUser(username, email, password) {
    const data = { name: username, email: email, password: password };
    try {
        const res = await fetch(
            '/api/users',
            {
                method: 'POST',
                body: JSON.stringify(data),
            })
        const jsonResponse = await res.json();
        if (!res.ok) {
            return jsonResponse;
        }
        return jsonResponse
    } catch (err) {
        return { "error": err.message }
    }
}
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
        if (!res.ok) {
            return jsonResponse;
        }
        return jsonResponse;
    } catch (err) {
        return { "error": err.message };
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
        return jsonResponse;
    } catch (err) {
        return { "error": err.message };
    }
}

// logout returns an object with an ok and error property. If ok is true then
// the logout was successful, otherwise, error contains the reason for the failure.
async function logout() {
    try {
        const res = await fetch('/api/user/logout', { method: 'DELETE' })
        if (res.status === 204) {
            console.log('Logout Success');
            return { ok: true };
        } else {
            return { error: 'Unauthorized use of button' };
        }
    } catch (err) {
        console.log(err);
        return { "error": err.message };
    }
}
export { loginUser, getUser, logout, registerUser };