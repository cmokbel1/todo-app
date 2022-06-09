import { useState } from 'react';
import { registerUser } from '../http/user';
import { default as Login } from './login';

const Registeration = () => {

    const [userRegInput, setUserRegInput] = useState('');
    const [passwordRegInput, setPasswordRegInput] = useState('');
    const [emailRegInput, setEmailRegInput] = useState('');
    const [loginState, setLoginState] = useState(false);

    const [errorState, setErrorState] = useState('');

    const handleRegister = async (event) => {
        event.preventDefault();
        const res = await registerUser(userRegInput, emailRegInput, passwordRegInput);
        if (!res.name) {
            setErrorState(res.error);
            return res.error;
        } else {
            setUserRegInput('');
            setPasswordRegInput('');
            setEmailRegInput('');
        }

    }
    if (loginState) {
        return <Login />
    } else {
        return (
            <div className="col-12 offset-md-4 col-md-4">
                <form className="card mx-auto px-4 py-2">
                    <div className="card-body">
                        <p className="card-title lead text-center mb-4">Register an account</p>
                        <div className="input mb-3">
                            <label className="input-text pb-2" htmlFor="username">username</label>
                            <input type="text" id="rUsername" name="username" className="form-control"
                                aria-label="username-input" aria-describedby="username-input"
                                defaultValue={setUserRegInput} onChange={(e) => setUserRegInput(e.target.value)} />
                        </div>
                        <div className="input mb-3">
                            <label className="input-text pb-2" htmlFor="email">email</label>
                            <input type="text" id="rEmail" name="email" className="form-control"
                                aria-label="password-input" aria-describedby="password-input"
                                defaultValue={setEmailRegInput} onChange={(e) => setEmailRegInput(e.target.value)} />
                        </div>
                        <div className="input mb-3">
                            <label className="input-text pb-2" htmlFor="password">password</label>
                            <input type="password" id="rPassword" name="password" className="form-control"
                                aria-label="password-input" aria-describedby="password-input"
                                defaultValue={setPasswordRegInput} onChange={(e) => setPasswordRegInput(e.target.value)} />
                        </div>
                        <p style={{ color: 'red' }}>{errorState}</p>
                        <button type="submit" className="btn btn-primary w-100 text-center mt-3"
                            onClick={handleRegister}
                            disabled={userRegInput === "" || passwordRegInput === ""}>
                            Register</button>
                            <a href="!#" onClick={() => setLoginState(true)}>Login</a>
                    </div>
                </form>
            </div>
        )
    }
}
export default Registeration