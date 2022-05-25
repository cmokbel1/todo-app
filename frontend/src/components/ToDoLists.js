import { useState } from 'react';
import { getList } from '../http/lists'

export const ToDoLists = ({lists}) => {
    return (
        <div>
            <ul className="list-group">
                {lists.map((list,index) => <li className="list-group-item" key={index}><button className="btn" onClick={() => getList(list.id)}>{list.name}</button></li>)}
            </ul>
        </div>
    )
}
