import { useEffect, useState } from 'react';
import { addItem, setCompletion } from "../http/lists";
import { Item } from './Item';

export const ListDetail = ({ id, name, completed, items, handleUpdate }) => {
    const [messageState, setMessageState] = useState('');
    const [errorMessageState, setErrorMessageState] = useState('');
    const [itemsState, setItemsState] = useState(items ? items : []);
    const [newItemName, setNewItemName] = useState('');
    const [currentName, setCurrentName] = useState(name)


    // takes an input value and adds it to the selectedList when enter is pressed
    const handleAddItem = async (event) => {
        if (event.charCode === 13) {
            if (!newItemName) {
                setErrorMessageState('Item name cannot be empty');
                return;
            }
            const res = await addItem(id, newItemName);
            if (res.error) {
                setErrorMessageState(res.error);
            } else {
                setErrorMessageState('');
                setMessageState('Task Added Successfully');
                setItemsState([...itemsState, res])
                console.log(items)
            }
            setNewItemName('');
            setTimeout(() => {
                setMessageState('');
            }, 1000)
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

    const handleListUpdate = (e) => {
        if (e.charCode === 13) {
            handleUpdate(id, currentName)
        }
    }

    useEffect(() => {
        if (id) {
            setItemsState(items)
            setCurrentName(name)
        }
    }, [items, id, name])

    let body = <h1>Nothing to see here</h1>
    if (id) {
        body = <>
            <input className="fs-3" rows="2" value={currentName} onChange={(e) => setCurrentName(e.target.value)} onKeyPress={(e) => handleListUpdate(e)}></input>
            <ul className="list-group">
                {itemsState.map((item, index) => {
                    return <Item id={item.id} name={item.name} completed={item.completed} setCompleted={handleSetCompleted} key={index} />
                })}
            </ul>
            <input type="text" name="item" className="form-input w-50" rows="2"
                onChange={(e) => { setNewItemName(e.target.value) }} onKeyPress={(e) => handleAddItem(e)}
                placeholder="Add Item" value={newItemName}></input>
            <p className="text-center">{messageState}</p><p className="text-center" style={{ color: 'red' }}>{errorMessageState}</p>
            <button className="btn btn-danger">Delete</button>
        </>
    }
    return body
}
