import { default as Login } from './login';
import { ToDoLists } from './ToDoLists';



export function Main({ userState, setUserState, setReturnError }) {
    let body = <Login setUserState={setUserState} setReturnError={setReturnError} />
    if (userState) {
        body = <ToDoLists userState={userState} />
    }
    return body
}

