import { setCompletion } from '../http/lists';
import { useState } from 'react';

export const Item = ({ name, completed, index, id, listId }) => {
    const [completionMessage, setCompletionMessage] = useState('');
    const [completionState, setCompletionState] = useState(completed)
    // this function passes in a data from the item to a helper function
    // helper function sets the completed state of an item to true or false
    const handleCompletion = async (id, completed, listId) => {
        const res = await setCompletion(id, completed, listId);
        if (res.error) {
            console.log(res.error)
        }
        setCompletionState(!completionState)
        setCompletionMessage('Successfully updated.');
        setTimeout(() => {
            setCompletionMessage('');
        }, 1000)
    }


    return (
        <>
            <li className="list-group-item" key={index}>
                {
                    completionState ?
                        <>
                            <span className="mr-3"><s>{name}</s></span>
                            <input className="form-check-input p-6" defaultChecked type="checkbox" onChange={() => setCompletionState(false)} onClick={() => handleCompletion(id, completed, listId)} />
                        </>
                        :
                        <>
                            <span className="mr-3">{name}</span>
                            <input className="form-check-input p-6" type="checkbox" onChange={() => setCompletionState(true)} onClick={() => handleCompletion(id, completed, listId)} />
                        </>
                }
            </li>
            <p>{completionMessage}</p>
        </>
    )
}