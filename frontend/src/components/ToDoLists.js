import { getList } from '../http/lists'

export const ToDoLists = ({ lists, listState, setListState }) => {
    // when the list button is clicked the API will return a list item and
    // we want to set that list item to a state which will then be passed up
    // this will allow us to render the current list item onto the main page
    const handleListClick = async(id) => {
        const list = await getList(id);
        if (list.name) {
            console.log(list)
            setListState(list)
            return listState
        } else {
            return list.error
        }
    }

    return (
        <div>
            <ul className="list-group">
                {lists.map((list, index) =>
                    <li className="list-group-item" key={index}>
                        <button className="btn" onClick={() => handleListClick(list.id)}>
                            {list.name}
                        </button>
                    </li>
                )}
            </ul>
        </div>
    )
}
