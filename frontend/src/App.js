import { useEffect, useState } from 'react';
import './App.css';

import { Main } from './components/Main';
import { getUser } from './http/user';
import { Header, Footer } from './components/Nav';

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
  }, [userState])

  return (
    <div className="row">
      <div className="col">
        <Header userState={userState} setUserState={setUserState} setReturnError={setReturnError} />
        <Main userState={userState} setUserState={setUserState} setReturnError={setReturnError} />
        <Footer />
      </div>
      <div className="col">
      </div>
    </div>
  )
}

export default App;
