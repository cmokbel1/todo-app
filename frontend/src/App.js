import { useEffect, useState } from 'react';
import './App.css';

import { default as Login } from './components/Login';
import { default as ToDoColumn } from './components/ToDoColumn';

import { getUser } from './http/user';


function App() {

  const [userState, setUserState] = useState()
  const [returnError, setReturnError] = useState(null)

  useEffect(() => {
    const res = getUser();
    if (res.name) {
      setUserState(res.name)
    } else {
      setReturnError(res)
    }
  }, [])

  if (userState) {
    return (
      <div className="row">
        <ToDoColumn userState={userState} />
        <div className="col" style={{ border: '1px solid blue' }}>
          <h2>Welcome Home</h2>
        </div>
      </div>
    )
  } else if (returnError) {
    return (
      <div className="row">
        <ToDoColumn userState={userState} />
        <div className="col" style={{ border: '1px solid blue' }}>
          <Login userState={userState} setUserState={setUserState} />
        </div>
      </div>

    )
  } else {
    return (
      <div className="row">
        <ToDoColumn userState={userState} />
        <div className="col align-middle" style={{ border: '1px solid blue' }}>
          <h1 className="text-center">Something has gone wrong</h1>
        </div>
      </div>
    )
  }
}

export default App;
