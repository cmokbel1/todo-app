import { useEffect, useState } from 'react';
import './App.css';

import { default as Login } from './components/login';
import { default as ToDoList } from './components/toDoList';

import { getUser } from './http/user';


function App() {

  const [userState, setUserState] = useState()

  useEffect(() => {
    getUser();
  }, [])

  if (getUser) {
    return (
      <div className="row">
        <ToDoList />
        <div className="col" style={{ border: '1px solid blue' }}>
          <h2>Welcome Home</h2>
        </div>
      </div>
    )
  } else {
    return (
      <div className="row">
        <ToDoList />
        <div className="col" style={{ border: '1px solid blue' }}>
          <Login userState={userState} setUserState={setUserState} />
        </div>
      </div>

    )
  }
}

export default App;
