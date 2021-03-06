import { useEffect, useState } from 'react';
import './App.css';

import { Main } from './components/Main';
import { getUser } from './http/user';
import { Header, Footer } from './components/Nav';
import { FlashMessage } from './components/FlashMessage';

function App() {

  const [userState, setUserState] = useState()
  const [returnError, setReturnError] = useState(null)
  const [messageState, setMessageState] = useState(null)

  useEffect(() => {
    getUser().then(res => {
      if (res.name) {
        setUserState(res.name);
      }
    });
  }, [userState])

  return (
    <>
      <div className="app d-flex flex-column">
        <Header userState={userState} setUserState={setUserState} setReturnError={setReturnError} />
        <div className="container-fluid flex-fill">
          <div className="row">
            <FlashMessage messageState={messageState}
              setMessageState={setMessageState}
              returnError={returnError} />
            <Main userState={userState}
              setUserState={setUserState}
              setReturnError={setReturnError}
              setMessageState={setMessageState} />
          </div>
        </div>
        <Footer />
      </div>
    </>
  )
}

export default App;
