export const Item = ({ id, name, completed, setCompleted, index}) => {
    // this function passes in a data from the item to a helper function
    // helper function sets the completed state of an item to true or false
    let styledName = name
    if (completed) {
        styledName = <s>{name}</s>
    }

    return ( 
        <li className="list-group-item" key={index}>
            <span className="mr-3">{styledName}</span>
            <input className="form-check-input p-6"
                checked={completed}
                type="checkbox"
                onClick={() => setCompleted(id, !completed)} />
        </li>
    )
}