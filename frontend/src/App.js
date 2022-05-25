import { useEffect, useState } from 'react';
import './App.css';

import { default as Login } from './components/Login';
import { default as ToDoColumn } from './components/ToDoColumn';

import { getUser } from './http/user';


function App() {

  const [userState, setUserState] = useState()

  useEffect(() => {
    getUser();
  }, [])

  if (getUser === 200) {
    return (
      <div className="row">
        <ToDoColumn userState={userState}/>
        <div className="col" style={{ border: '1px solid blue' }}>
          <h2>Welcome Home</h2>
        </div>
      </div>
    )
  } else if (getUser === 401) {
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
        <div className="col" style={{ border: '1px solid blue' }}>
        <h1>Something has gone wrong</h1>
        </div>
      </div>
    )
  }
}

export default App;
