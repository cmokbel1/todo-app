import { useEffect } from 'react'
import './App.css';
import { default as Login } from './components/login'
import { checkUser } from './http/user'


function App() {
  useEffect(() => {
    checkUser();
  }, [])

    return (
      <Login />
    )
}

export default App;
