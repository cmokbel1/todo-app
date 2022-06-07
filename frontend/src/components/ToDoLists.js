import { getList, getLists, addList, updateListName, deleteList } from '../http/lists';
import { useState, useEffect } from 'react';
import { ListDetail } from './ListDetail';

export const ToDoLists = ({ userState, setReturnError, setMessageState }) => {
    const [lists, setLists] = useState([]);
    const [selectedList, setSelectedList] = useState();
    const [errorMessage, setErrorMessage] = useState('')
    const [newListName, setNewListName] = useState('')
  
    useEffect(() => {
        getLists().then(res => {
            setLists(res)
            setSelectedList(res[0])
        })
    }, [userState])


    // when the list button is clicked the API will return a list item and
    // we want to set that list item to a state which will then be passed up
    // this will allow us to render the current list item onto the main page
    async function handleListClick(id) {
        const res = await getList(id);
        if (res.error) {
            return res.error
        }
        setSelectedList(res);

    }

    // need to abstract away this function...
    const handleAddList = async (event) => {
        if (event.charCode === 13) {
            if (!newListName) {
                setErrorMessage('List name cannot be empty.');
                return;
            }
            const res = await addList(newListName);
            if (res.error) {
                setReturnError(res.error);
            } else {
                setErrorMessage('');
                setReturnError('');
                setMessageState('List added successfully.');
                setLists([...lists, res])
                setSelectedList(res);
            }
            setNewListName('');
            setTimeout(() => {
                setMessageState('');
            }, 1000)
        }
    }
    // handler for the list name update
    const handleListNameUpdate = async (id, name) => {
        const res = await updateListName(id, name)
        if (res.error) {
            setReturnError(res.error)
        }
        else {
            setReturnError('');
            setMessageState('List updated successfully.')
            const newLists = lists.map(l => l.id === id ? res : l)
            setLists(newLists)
            setSelectedList(res)
            setTimeout(() => {
                setMessageState('');
            }, 1000)
        }
    }
    // handler for deleting list
    const handleDeleteList = async (listId) => {
        // TODO(cmokbel1): use custom modal instead of window confirm
        if (!window.confirm("Delete entire list?")) {
            return;
        }
        const res = await deleteList(listId);
        if (res === "") {
            const newLists = lists.filter(l => l.id !== listId ? l : null)
            setLists(newLists)
            // only reset the selectedList if we delete the selectedList
            if (listId === selectedList.id) {
                setSelectedList(newLists[0])
            }
        } else {
            setReturnError('An error occurred.')
        }
    }

    const trash = <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-trash" viewBox="0 0 16 16">
        <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z" />
        <path fillRule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z" />
    </svg>

    let body = <p>Nothing to see here</p>
    if (lists) {
        body =
            <div className="row">
                <div className='col-12 col-md-3'>
                    <ul className="list-group border rounded border-dark shadow mb-2">
                        {lists.map((list, index) =>
                            <li className="list-group-item d-flex justify-content-between" key={index}>
                                <button className="btn" onClick={() => handleListClick(list.id)}>
                                    {list.name}
                                </button>
                                <button className="btn btn-secondary" onClick={() => handleDeleteList(list.id)}>{trash}</button>
                            </li>
                        )}
                    </ul>
                    <input type="text" name="item" className="form-input w-100 mt-2"
                        onChange={(e) => { setNewListName(e.target.value) }} onKeyPress={(e) => handleAddList(e)}
                        placeholder="+ add list" value={newListName}></input>
                        <p style={{color: 'red'}}>{errorMessage}</p>
                </div>
                <div className='col-12 col-md-9'>
                    <ListDetail {...selectedList} handleUpdate={handleListNameUpdate} removeList={handleDeleteList} setReturnError={setReturnError} setMessageState={setMessageState} />
                </div>
            </div>
    }

    return (
        <div>
            {body}
        </div>
    )
}