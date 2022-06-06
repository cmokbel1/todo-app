import { useEffect, useState } from 'react';
import { addItem, setCompletion, deleteItem } from "../http/lists";
import { Item } from './Item';

export const ListDetail = ({ id, name, completed, items, handleUpdate, removeList, setReturnError, setMessageState }) => {
    const [itemsState, setItemsState] = useState(items ? items : []);
    const [newItemName, setNewItemName] = useState('');
    const [currentName, setCurrentName] = useState(name)
    const [errorMessage, setErrorMessage] = useState(null)


    // takes an input value and adds it to the selectedList when enter is pressed
    const handleAddItem = async (event) => {
        if (event.charCode === 13) {
            if (!newItemName) {
                setErrorMessage('Item name cannot be empty.');
                return;
            }
            const res = await addItem(id, newItemName);
            if (res.error) {
                setReturnError(res.error);
            } else {
                setReturnError(null);
                setErrorMessage(null)
                setMessageState('Item added successfully.');
                setItemsState([...itemsState, res])
            }
            setNewItemName('');
            setTimeout(() => {
                setMessageState('');
            }, 2000)
        }
    }

    const handleSetCompleted = async (itemId, completed) => {
        const res = await setCompletion(itemId, completed, id)
        if (res.error) {
            console.log(res.error)
            return
        }
        const newItems = itemsState.map(i => i.id === itemId ? res : i)
        setItemsState(newItems)
    }

    const handleListUpdate = (event) => {
        if (event.charCode === 13) {
            handleUpdate(id, currentName)
        }
    }

    const handleDeleteItem = async (itemId) => {
        const res = await deleteItem(id, itemId);
        if (res.error) {
            setReturnError(res.error);
            return;
        }
        const newItems = itemsState.filter(i => i.id !== itemId ? i : null)
        console.log(newItems)
        setItemsState(newItems)
    }


    useEffect(() => {
        if (id) {
            setItemsState(items)
            setCurrentName(name)
        }
    }, [items, id, name])

    let body = <h1>Nothing to see here</h1>
    if (id) {
        body = <div className="text-center border border-dark rounded shadow">
            <input className="fs-3 mb-4 mt-2 text-center w-50" rows="2" type="text" value={currentName} onChange={(e) => setCurrentName(e.target.value)} onKeyPress={(e) => handleListUpdate(e)}></input>
            <ul className="list-group mb-4">
                {itemsState.map((item, index) => {
                    return <Item id={item.id} name={item.name}
                        completed={item.completed}
                        setCompleted={handleSetCompleted}
                        key={index} deleteItem={handleDeleteItem} />
                })}
            </ul>
            <input type="text" name="item" className="form-input w-50"
                onChange={(e) => { setNewItemName(e.target.value) }} onKeyPress={(e) => handleAddItem(e)}
                placeholder="+ add item" value={newItemName}></input>
                <p style={{color:'red'}}>{errorMessage}</p>
            <button className="btn btn-danger mb-2" onClick={() => removeList(id)}>Delete</button>
        </div>
    }
    return body
}