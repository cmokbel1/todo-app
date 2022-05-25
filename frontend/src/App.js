import { useEffect, useState } from 'react';
import './App.css';

import { default as Login } from './components/login';
import { default as ToDoColumn } from './components/ToDoColumn';

import { getUser, logout } from './http/user';


function App() {

  const [userState, setUserState] = useState()
  const [returnError, setReturnError] = useState(null)

  useEffect(() => {
    getUser().then(res => {
      if (res.name) {
        setUserState(res.name)
      } else {
        setReturnError(res)
      }
    });
  }, [])
  
    const handleLogout = () => {
      logout().then(res => {
        if (res.ok) {
          setUserState();
        } else {
          setReturnError(res.error);
        }
      })
    }

  if (userState) {
    return (
      <div className="row">
        <ToDoColumn />
        <div className="col" style={{ border: '1px solid blue' }}>
        <div  className="text-end"><span>Hello {userState} </span><button className="btn btn-primary" onClick={handleLogout}>Logout</button></div>
          <h2>Welcome Home</h2>
        </div>
      </div>
    )
  } else if (returnError) {
    return (
      <div className="row">
        <ToDoColumn  />
        <div className="col" style={{ border: '1px solid blue' }}>
          <Login userState={userState} setUserState={setUserState} />
        </div>
      </div>

    )
  } else {
    return (
      <div className="row">
        <ToDoColumn />
        <div className="col align-middle" style={{ border: '1px solid blue' }}>
          <h1 className="text-center">Something has gone wrong</h1>
        </div>
      </div>
    )
  }
}

export default App;
