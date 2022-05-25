import { useState } from 'react';

export const ToDoLists = ({lists}) => {
    return (
        <div>
            <ul className="list-group">
                {lists.map((list,index) => <li className="list-group-item" index={index}><button className="btn">{list.name}</button></li>)}
            </ul>
        </div>
    )
}
