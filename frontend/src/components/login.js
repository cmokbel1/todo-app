import { useState } from 'react';
import { registerUser, loginUser } from '../http/user'

const Login = ({ setUserState, setMessage, setReturnError }) => {
    //login states
    const [userInput, setUserInput] = useState('');
    const [passwordInput, setPasswordInput] = useState('');
    //register states
    const [userRegInput, setUserRegInput] = useState('');
    const [passwordRegInput, setPasswordRegInput] = useState('');
    const [emailRegInput, setEmailRegInput] = useState('');
    //error and toggler states
    const [errorState, setErrorState] = useState('');
    const [showRegister, setShowRegister] = useState(false)


    //login function
    const handleLogin = event => {
        //trim username and set to lowercase
        const userTrim = userInput.trim();
        const passTrim = passwordInput.trim();
        event.preventDefault();
        loginUser(userTrim, passTrim).then(res => {
            if (res.name) {
                setReturnError("");
                setUserState(res.name);
            } else {
                setErrorState(res);
            }
        });
    }
    //register function
    const handleRegister = async (event) => {
        event.preventDefault();
        const res = await registerUser(userRegInput, emailRegInput, passwordRegInput);
        if (!res.name) {
            setErrorState(res.error);
            return res.error;
        } else {
            setReturnError("")
            setMessage('Successfully created account.')
            setUserRegInput('');
            setPasswordRegInput('');
            setEmailRegInput('');
            setTimeout(() => {
                setMessage('')
            }, 1000)
        }

    }
    const handleToggler = async (event) => {
        if (!showRegister) {
            setShowRegister(true);
        } else {
            setShowRegister(false);
        }

    }
    if (showRegister) {
        return (
            <div className="col-12 offset-md-4 col-md-4">
                <form className="card mx-auto px-4 py-2">
                    <div className="card-body">
                        <p className="card-title lead text-center mb-4">Register an account</p>
                        <div className="input mb-3">
                            <label className="input-text pb-2"
                                htmlFor="username">username</label>
                            <input
                                type="text"
                                id="rUsername"
                                name="username"
                                className="form-control"
                                aria-label="username-input"
                                aria-describedby="username-input"
                                defaultValue={userRegInput}
                                onChange={(e) => setUserRegInput(e.target.value)} />
                        </div>
                        <div className="input mb-3">
                            <label className="input-text pb-2"
                                htmlFor="email">email</label>
                            <input
                                type="text"
                                id="rEmail"
                                name="email"
                                className="form-control"
                                aria-label="email-input"
                                aria-describedby="email-input"
                                defaultValue={emailRegInput}
                                onChange={(e) => setEmailRegInput(e.target.value)} />
                        </div>
                        <div className="input mb-3">
                            <label className="input-text pb-2"
                                htmlFor="password">password</label>
                            <input
                                type="password"
                                id="rPassword"
                                name="password"
                                className="form-control"
                                aria-label="password-input"
                                aria-describedby="password-input"
                                defaultValue={passwordRegInput}
                                onChange={(e) => setPasswordRegInput(e.target.value)} />
                        </div>
                        <p style={{ color: 'red' }}>{errorState}</p>
                        <button
                            type="submit"
                            className="btn btn-primary w-100 text-center mt-3"
                            onClick={handleRegister}
                            disabled={userRegInput === "" || passwordRegInput === ""}>
                            Register</button>
                        <button
                            className="btn btn-secondary w-100 text-center mt-3"
                            onClick={(event) => handleToggler(event)}>Login</button>
                    </div>
                </form>
            </div>
        )
    } else {
        return (
            <div className="col-12 offset-md-4 col-md-4">
                <form className="card mx-auto px-4 py-2">
                    <div className="card-body">
                        <p className="card-title lead text-center mb-4">You must login to continue</p>
                        <div className="input mb-3">
                            <label
                                className="input-text pb-2"
                                htmlFor="username">Username</label>
                            <input
                                type="text"
                                id="username"
                                name="username"
                                className="form-control"
                                aria-label="username-input"
                                aria-describedby="username-input"
                                onChange={(e) => setUserInput(e.target.value)} />
                        </div>
                        <div className="input mb-3">
                            <label
                                className="input-text pb-2"
                                htmlFor="password">Password</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                className="form-control" aria-label="password-input"
                                aria-describedby="password-input"
                                onChange={(e) => setPasswordInput(e.target.value)} />
                        </div>
                        <p style={{ color: 'red' }}>{errorState}</p>
                        <button
                            type="submit"
                            className="btn btn-primary w-100 text-center mt-3"
                            onClick={handleLogin}
                            disabled={userInput === "" || passwordInput === ""}>Login</button>
                        <button
                            className="btn btn-secondary w-100 text-center mt-3"
                            onClick={(event) => handleToggler(event)}>Register</button>
                    </div>
                </form>
            </div>
        )
    }
}

export default Login