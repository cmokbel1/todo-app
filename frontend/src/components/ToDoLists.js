import { getList } from '../http/lists';
import { useState } from 'react';

export const ToDoLists = ({ lists, selectedList, setSelectedList }) => {
    // when the list button is clicked the API will return a list item and
    // we want to set that list item to a state which will then be passed up
    // this will allow us to render the current list item onto the main page
    const handleListClick = async (id) => {
        const list = await getList(id);
        if (list.name) {
            console.log(list)
            setSelectedList(list)
        } else {
            return list.error
        }
    }

    let body = <p>Nothing to see here</p>
    if (lists) {
        body = <ul className="list-group">
            {lists.map((list, index) =>
                <li className="list-group-item" key={index}>
                    <button className="btn" onClick={() => handleListClick(list.id)}>
                        {list.name}
                    </button>
                </li>
            )}
        </ul>
    }

    return (
        <div>
            {body}
        </div>
    )
}
