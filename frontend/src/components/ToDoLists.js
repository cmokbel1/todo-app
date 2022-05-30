import { getList, getLists, addList, updateListName } from '../http/lists';
import { useState, useEffect } from 'react';
import { ListDetail } from './ListDetail';

export const ToDoLists = ({ userState }) => {
    const [lists, setLists] = useState([]);
    const [selectedList, setSelectedList] = useState();

    const [newListName, setNewListName] = useState('')
    const [messageState, setMessageState] = useState('');
    const [errorMessageState, setErrorMessageState] = useState('');
    const [updatedListName, setUpdatedListName] = useState(selectedList)

    useEffect(() => {
        getLists().then(res => {
            setLists(res)
            setSelectedList(res[0])
        })
    }, [userState])


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
    const handleAddList = async (event) => {
        if (event.charCode === 13) {
            if (!newListName) {
                setErrorMessageState('List name cannot be empty');
                return;
            }
            const res = await addList(newListName);
            if (res.error) {
                setErrorMessageState(res.error);
            } else {
                setErrorMessageState('');
                setMessageState('List Successfully Added');
                setLists([...lists, res])
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
                    <input type="text" name="item" className="form-input" onChange={(e) => { setNewListName(e.target.value) }} onKeyPress={(e) => handleAddList(e)} placeholder="+ add list" value={newListName}></input>
                    <p className="text-center">{messageState}</p><p className="text-center" style={{ color: 'red' }}>{errorMessageState}</p>
                </div>
                <div className='col-12 col-md-9'>
                    <ListDetail {...selectedList} />
                </div>
            </div>
    }

    return (
        <div>
            {body}
        </div>
    )
}
