
export const List = ({listState}) => {
    return (
        <div className="container">
            <h1><u>{listState.name}</u></h1>
            <ul className="list-group">
                {listState.items.map((item, index) => {
                    return (
                        <li className="list-group-item" key={index}>
                            {item.name}
                            <label className="" htmlFor="checkbox" name="completed">Completed</label>
                            <input className="form-check-input" type="checkbox"/>
                        </li>
                    )
                })}
            </ul>
        </div>
    )
}
