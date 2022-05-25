import { useState } from 'react';
import { getList } from '../http/lists'

export const ToDoLists = ({lists, listState, setListState}) => {

    const handleListClick = async(id) => {
        const list = await getList(id);
        if (list.ok) {
            setListState(list)
        } else {
            return list.error
        }
    }
    
    return (
        <div>
            <ul className="list-group">
                {lists.map((list,index) => <li className="list-group-item" key={index}><button className="btn" onClick={() => handleListClick(list.id)}>{list.name}</button></li>)}
            </ul>
        </div>
    )
}
