import { useState } from 'react';

export const ToDoLists = ({userState}) => {
    return (
        <div>
            <ul className="list-group">
                <li className="list-group-item active"><a className="link-dark" href="!#">Placeholder</a></li>
                <li className="list-group-item"><a className="link-dark" href="!#">Placeholder</a></li>
                <li className="list-group-item"><a className="link-dark" href="!#">Placeholder</a></li>
                <li className="list-group-item"><a className="link-dark" href="!#">Placeholder</a></li>
                <li className="list-group-item"><a className="link-dark" href="!#">Placeholder</a></li>
                {/* user.lists.map(list => <li><button>list.name</button></li>) */}
            </ul>
        </div>
    )
}
