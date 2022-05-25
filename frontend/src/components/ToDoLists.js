import { useState } from 'react';

export const ToDoLists = ({lists}) => {
    return (
        <div>
            <ul className="list-group">
                {lists.map(list => <li><button className="btn btn-secondary">{list.name}</button></li>)}
            </ul>
        </div>
    )
}
