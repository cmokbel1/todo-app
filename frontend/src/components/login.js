import { useState } from 'react';
import { loginUser } from '../http/user'

const Login = ({ userState, setUserState }) => {
    const [userInput, setUserInput] = useState('');
    const [passwordInput, setPasswordInput] = useState('');
    const [errorState, setErrorState] = useState('');

    const handleLogin = event => {
        event.preventDefault();
        loginUser(userInput, passwordInput).then(res => {
            if (res.name) {
                 setUserState(res.name);
            } else {
                setErrorState(res);
            }
            console.log(res);
        });
    }
    return (
        <>
            <form className="card justify-content-center" style={{ width: '18em' }}>
                <p className="card-title">You must login to continue</p>
                <div className="input col-sm-8 mb-3">
                    <label className="input-text" htmlFor="username">username</label>
                    <input type="text" id="username" name="username" className="form-control" aria-label="username-input" aria-describedby="username-input" onChange={(e) => setUserInput(e.target.value)} />
                </div>
                <div className="input col-sm-8 mb-3">
                    <label className="input-text" htmlFor="password">password</label>
                    <input type="password" id="password" name="password" className="form-control" aria-label="password-input" aria-describedby="password-input" onChange={(e) => setPasswordInput(e.target.value)} />
                </div>
                <p style={{color: 'red'}}>{errorState}</p>
                <button type="submit" className="btn btn-primary" onClick={handleLogin} disabled={userInput === "" || passwordInput === ""}>Login</button>
                
            </form>
        </>
    )
}

export default Login