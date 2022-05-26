import { useEffect, useState } from 'react';
import './App.css';

import { default as Login } from './components/Login';
import { default as ToDoColumn } from './components/ToDoColumn';
import { List } from './components/List'

import { getUser } from './http/user';


function App() {

  const [userState, setUserState] = useState()
  const [returnError, setReturnError] = useState(null)
  const [listState, setListState] = useState()

  useEffect(() => {
    getUser().then(res => {
      if (res.name) {
        setUserState(res.name)
      } else {
        setReturnError(res)
      }
    });
  }, [])

  if (userState) {
    return (
      <div className="row">
        <ToDoColumn userState={userState} setUserState={setUserState} listState={listState} setListState={setListState} setReturnError={setReturnError}/>
        <div className="col">
          {listState ?
            <List listState={listState} />
            : <h2>List Not Found</h2>}
        </div>
      </div>
    )
  } else if (returnError) {
    return (
      <div className="row">
        <ToDoColumn />
        <div className="col">
          <Login userState={userState} setUserState={setUserState} />
        </div>
      </div>

    )
  } else {
    return (
      <div className="row">
        <ToDoColumn />
        <div className="col align-middle">
          <h1 className="text-center">Something has gone wrong</h1>
        </div>
      </div>
    )
  }
}

export default App;
