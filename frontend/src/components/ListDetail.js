import { useState } from 'react';
import { addItem } from "../http/lists";

export const ListDetail = ({ selectedList, setSelectedList }) => {
    const [messageState, setMessageState] = useState('')
    const [newItemName, setNewItemName] = useState('')
// Takes the value from the input at the bottom of the list
// and adds it to the list of items within a specific list
// checks to see if the key pressed was the enter key 
// we display a message below the input if successful or error
    const handleAddItem = async (event) => {
        if (event.charCode === 13) {
            const res = await addItem(selectedList.id, newItemName);
            if (res.error) {
                setMessageState("Input required");
            } else {
                setMessageState('Post Successful');

            }
            setNewItemName('');
            setTimeout(() => {
                setMessageState('');
            }, 1000)

        }
    }
    let body = <h1>Nothing to see here</h1>
    if (selectedList) {
        body = <>
            <h1><u>{selectedList.name}</u></h1>
            <ul className="list-group">
                {selectedList.items.map((item, index) => {
                    return (
                        <li className="list-group-item" key={index}>
                            {item.name}
                            <label className="" htmlFor="checkbox" name="completed">Completed</label>
                            <input className="form-check-input" type="checkbox" />
                        </li>
                    )
                })}
            </ul>
            <input type="text" name="item" className="form-input" onChange={(e) => { setNewItemName(e.target.value) }} onKeyPress={(e) => handleAddItem(e)} placeholder="Add Item" value={newItemName}></input>
            <p className="text-center">{messageState}</p>
        </>
    }
    return body
}
