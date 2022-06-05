import { default as Login } from './login';
import { ToDoLists } from './ToDoLists';



export function Main({ userState, setUserState, setReturnError, setMessageState }) {
    let body = <Login setUserState={setUserState} setReturnError={setReturnError} />
    if (userState) {
        body = <ToDoLists userState={userState} setReturnError={setReturnError} setMessageState={setMessageState} />
    }
    return body
}

