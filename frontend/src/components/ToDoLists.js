import { getList } from '../http/lists';
import { useState } from 'react';

export const ToDoLists = ({ lists, selectedList, setSelectedList, addList }) => {
    const [newListName, setNewListName] = useState('')
    const [messageState, setMessageState] = useState('');
    const [errorMessageState, setErrorMessageState] = useState('');
    // when the list button is clicked the API will return a list item and
    // we want to set that list item to a state which will then be passed up
    // this will allow us to render the current list item onto the main page
    const handleListClick = async (id) => {
        const list = await getList(id);
        if (!list.name) {
            return list.error
        } else {
            console.log(list)
            setSelectedList(list)
        }
    }
    
    // need to abstract away this function...
    const handleAddItem = async (event) => {
        if (!newListName) {
            setErrorMessageState('Input Required');
            return;
        }
        if (event.charCode === 13) {
            const res = await addList(newListName);
            if (res.error) {
                setErrorMessageState(res.error);
            } else {
                setErrorMessageState('');
                setMessageState('Post Successful');
            }
            setNewListName('');
            setTimeout(() => {
                setMessageState('');
            }, 1000)
        }
    }


    let body = <p>Nothing to see here</p>
    if (lists) {
        body =
            <>
                <ul className="list-group">
                    {lists.map((list, index) =>
                        <li className="list-group-item" key={index}>
                            <button className="btn" onClick={() => handleListClick(list.id)}>
                                {list.name}
                            </button>
                        </li>
                    )}
                </ul>
                <input type="text" name="item" className="form-input" onChange={(e) => { setNewListName(e.target.value) }} onKeyPress={(e) => handleAddItem(e)} placeholder="+ add list" value={newListName}></input>
            </>
    }

    return (
        <div>
            {body}
        </div>
    )
}
