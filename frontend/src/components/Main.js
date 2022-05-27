import { useEffect, useState } from 'react';
import { default as Login } from './Login';
import { ListDetail } from './ListDetail';
import { ToDoLists } from './ToDoLists';
import { getLists } from '../http/lists'


export function Main({ userState, setUserState }) {
    const [lists, setLists] = useState([]);
    const [selectedList, setSelectedList] = useState();

    useEffect(() => {
        if (userState) {
            getLists().then(res => {
                setLists(res)
                setSelectedList(res[0])
            })
        } else {
            setLists([])
            setSelectedList()
        }
    }, [userState])

    let body = <Login setUserState={setUserState} />
    if (userState) {
        body = <ListDetail selectedList={selectedList} setSelectedList={setSelectedList} />
    }
    return (
        <>
            <ToDoLists lists={lists} selectedList={selectedList} setSelectedList={setSelectedList} />
            {body}
        </>
    )
}

