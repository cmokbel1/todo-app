import { useState } from 'react';
import { loginUser } from '../http/user'

const Login = ({ setUserState, setReturnError }) => {
    const [userInput, setUserInput] = useState('');
    const [passwordInput, setPasswordInput] = useState('');
    const [loginError, setLoginError] = useState('');


    const handleLogin = event => {
        //trim username and set to lowercase
        const userTrim = userInput.trim();
        const passTrim = passwordInput.trim();
        event.preventDefault();
        loginUser(userTrim, passTrim).then(res => {
            if (res.name) {
                setUserState(res.name);
            } else {

                setLoginError(res);
            }
        });
    }
    return (
        <div className="col-12 offset-md-4 col-md-4">
            <form className="card mx-auto px-4 py-2">
                <div className="card-body">
                    <p className="card-title lead text-center mb-4">You must login to continue</p>
                    <div className="input mb-3">
                        <label className="input-text pb-2" htmlFor="username">Username</label>
                        <input type="text" id="username" name="username" className="form-control" aria-label="username-input" aria-describedby="username-input" onChange={(e) => setUserInput(e.target.value)} />
                    </div>
                    <div className="input mb-3">
                        <label className="input-text pb-2" htmlFor="password">Password</label>
                        <input type="password" id="password" name="password" className="form-control" aria-label="password-input" aria-describedby="password-input" onChange={(e) => setPasswordInput(e.target.value)} />
                    </div>
                    <p style={{ color: 'red' }}>{loginError}</p>
                    <button type="submit" className="btn btn-primary w-100 text-center mt-3" onClick={handleLogin} disabled={userInput === "" || passwordInput === ""}>Login</button>
                </div>
            </form>
        </div>
    )
}

export default Login