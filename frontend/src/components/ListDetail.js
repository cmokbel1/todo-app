import { useEffect, useState } from 'react';
import { addItem } from "../http/lists";
import { Item } from './Item';

export const ListDetail = ({ selectedList }) => {
    const [messageState, setMessageState] = useState('');
    const [errorMessageState, setErrorMessageState] = useState('');
    const [items, setItems] = useState([])
    const [newItemName, setNewItemName] = useState('');
    // takes an input value and adds it to the selectedList when enter is pressed
    const handleAddItem = async (event) => {
        if (event.charCode === 13) {
            if (!newItemName) {
                setErrorMessageState('Item name cannot be empty');
                return;
            }
            const res = await addItem(selectedList.id, newItemName);
            if (res.error) {
                setErrorMessageState(res.error);
            } else {
                setErrorMessageState('');
                setMessageState('Task Added Successfully');
                setItems([...items, res])
            }
            setNewItemName('');
            setTimeout(() => {
                setMessageState('');
            }, 1000)
        }
    }

    useEffect(() => {
        if (selectedList) {
            setItems(selectedList.items)
        }
    }, [selectedList])

    let body = <h1>Nothing to see here</h1>
    if (selectedList) {
        body = <>
            <h1><u>{selectedList.name}</u></h1>
            <ul className="list-group">
                {items.map((item, index) => {
                    return (
                        <Item name={item.name} completed={item.completed} key={index} id={item.id} listId={item.listId} />
                    )
                })}
            </ul>
            <input type="text" name="item" className="form-input" onChange={(e) => { setNewItemName(e.target.value) }} onKeyPress={(e) => handleAddItem(e)} placeholder="Add Item" value={newItemName}></input>
            <p className="text-center">{messageState}</p><p className="text-center" style={{ color: 'red' }}>{errorMessageState}</p>
        </>
    }
    return body
}
