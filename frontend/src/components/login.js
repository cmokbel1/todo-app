import { useState } from 'react';
import { loginUser } from '../http/user'

const Login = () => {
    const [user, setUser] = useState({
        username: '',
        password: ''
    })

    const handleInputChange = ({ target: { name, value } }) => { setUser({ ...user, [name]: value }) }

    const handleLogin = event => {
        event.preventDefault();
        loginUser(user.username, user.password);
    }

    return (
        <>
            <form className="card align-middle" style={{ width: '18em' }}>
                <div className="input col-sm-8 mb-3">
                    <h5 className="input-text" id="username">username</h5>
                    <input type="text" name="username" className="form-control" aria-label="username-input" aria-describedby="username-input" onChange={handleInputChange} />
                </div>
                <div className="input col-sm-8 mb-3">
                    <h5 className="input-text" id="password">password</h5>
                    <input type="password" name="password" className="form-control" aria-label="password-input" aria-describedby="password-input" onChange={handleInputChange} />
                </div>
                <button type="submit" className="btn btn-primary" onClick={handleLogin}>Login</button>
            </form>
        </>
    )
}

export default Login