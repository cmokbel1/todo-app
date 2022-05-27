import { useEffect, useState } from 'react';
import { default as Login } from './login';
import { ListDetail } from './ListDetail';
import { ToDoLists } from './ToDoLists';
import { getLists, addList } from '../http/lists';


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
        body =
                <ListDetail selectedList={selectedList} setSelectedList={setSelectedList} />

    }
    return (
        <>
            <div className="col-12 col-md-3 ">
                <ToDoLists lists={lists} setLists={setLists} selectedList={selectedList} setSelectedList={setSelectedList} addList={addList} />
            </div>
            <div className="col-12 col-md-9">
                {body}
            </div>
        </>
    )
}

