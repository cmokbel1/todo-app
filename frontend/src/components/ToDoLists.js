import { getList, getLists, addList, updateListName, deleteList } from '../http/lists';
import { useState, useEffect } from 'react';
import { ListDetail } from './ListDetail';

export const ToDoLists = ({ userState }) => {
    const [lists, setLists] = useState([]);
    const [selectedList, setSelectedList] = useState();

    const [newListName, setNewListName] = useState('')
    const [messageState, setMessageState] = useState('');
    const [errorMessageState, setErrorMessageState] = useState('');
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
                setErrorMessageState('List name cannot be empty.');
                return;
            }
            const res = await addList(newListName);
            if (res.error) {
                setErrorMessageState(res.error);
            } else {
                setErrorMessageState('');
                setMessageState('List successfully added.');
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
            setErrorMessageState(res.error)
        }
        else {
            setErrorMessageState('');
            setMessageState('Successfully updated list name.')
            const newLists = lists.map(l => l.id === id ? res : l)
            setLists(newLists)
            setSelectedList(res)
        }
    }
    // handler for deleting list
    const handleDeleteList = async (listId) => {
        // TODO(cmokbel1): use custom modal instead of window confirm
        if (!window.confirm("Are you sure?")) {
            return;
        }
        const res = await deleteList(listId);
        if (res === "") {
            const newLists = lists.filter(l => l.id !== listId ? l : null)
            setLists(newLists)
            setSelectedList(newLists[0])   
        } else {
            setErrorMessageState('An error occurred.')
        }
    }

    let body = <p>Nothing to see here</p>
    if (lists) {
        body =
            <div className='row'>
                <div className='col-12 col-md-3'>
                    <ul className="list-group">
                        {lists.map((list, index) =>
                            <li className="list-group-item" key={index}>
                                <button className="btn" onClick={() => handleListClick(list.id)}>
                                    {list.name}
                                </button>
                            </li>
                        )}
                    </ul>
                    <input type="text" name="item" className="form-input w-75"
                        onChange={(e) => { setNewListName(e.target.value) }} onKeyPress={(e) => handleAddList(e)}
                        placeholder="+ add list" value={newListName}></input>
                    <p className="text-center">{messageState}</p><p className="text-center" style={{ color: 'red' }}>{errorMessageState}</p>
                </div>
                <div className='col-12 col-md-9'>
                    <ListDetail {...selectedList} handleUpdate={handleListNameUpdate} removeList={handleDeleteList} />
                </div>
            </div>
    }

    return (
        <div>
            {body}
        </div>
    )
}
