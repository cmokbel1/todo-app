import { useState } from 'react';
import { loginUser } from '../http/user'

const Login = ({ setUserState }) => {
    const [userInput, setUserInput] = useState('');
    const [passwordInput, setPasswordInput] = useState('');
    const [errorState, setErrorState] = useState('');

    const handleLogin = event => {
        //trim username and set to lowercase
        const userTrim = userInput.trim();
        const passTrim = passwordInput.trim();
        event.preventDefault();
        loginUser(userTrim, passTrim).then(res => {
            if (res.name) {
                setUserState(res.name);
            } else {
                setErrorState(res);
            }
        });
    }
    return (
        <div className="col-6">
            <form className="card" style={{ width: '18em' }}>
                <div className="card-body">
                    <p className="card-title">You must login to continue</p>
                    <div className="input mb-3">
                        <label className="input-text" htmlFor="username">username</label>
                        <input type="text" id="username" name="username" className="form-control" aria-label="username-input" aria-describedby="username-input" onChange={(e) => setUserInput(e.target.value)} />
                    </div>
                    <div className="input mb-3">
                        <label className="input-text" htmlFor="password">password</label>
                        <input type="password" id="password" name="password" className="form-control" aria-label="password-input" aria-describedby="password-input" onChange={(e) => setPasswordInput(e.target.value)} />
                    </div>
                    <p style={{ color: 'red' }}>{errorState}</p>
                    <div className="text-center">
                        <button type="submit" className="btn btn-primary w-75 text-center" onClick={handleLogin} disabled={userInput === "" || passwordInput === ""}>Login</button>
                    </div>
                </div>
            </form>
        </div>
    )
}

export default Login