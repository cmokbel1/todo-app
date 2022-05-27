import { default as Login } from './login';
import { ToDoLists } from './ToDoLists';



export function Main({ userState, setUserState }) {
    let body = <Login setUserState={setUserState} />
    if (userState) {
        body = <ToDoLists userState={userState} />
    }
    return body
}

