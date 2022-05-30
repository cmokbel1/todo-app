// Get functions
async function getLists() {
    try {
        const res = await fetch('/api/todos');
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}

async function getList(id) {
    try {
        const res = await fetch(`/api/todos/${id}`);
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}

// Add functions
async function addItem(id, item) {
    const data = { name: item, listId: id }
    try {
        const res = await fetch(`/api/todos/${id}/`,
            {
                method: 'POST',
                body: JSON.stringify(data)
            })
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}

async function addList(item) {
    const data = { name: item, completed: false }
    try {
        const res = await fetch('/api/todos/',
            {
                method: 'POST',
                body: JSON.stringify(data)
            })
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch (err) {
        console.log(err);
        return err.error;
    }
}

// update functions!
async function setCompletion(id, completed, listId) {
    const data = { id: id, completed: completed }
    try {
        const res = await fetch(`/api/todos/${listId}/${id}`,
            {
                method: 'PATCH',
                body: JSON.stringify(data)
            })
            const jsonResponse = await res.json()
            return jsonResponse
    } catch (err) {
        console.log(err);
        return err.error;
    };
}

async function updateListName(id, name) {
    const data = { name: name }
    try {
        const res = await fetch(`/api/todos/${id}`,
        {
            method: 'PATCH',
            body: JSON.stringify(data)
        })
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch(err) {
        console.log(err);
        return err.error;
    };
}
export { getLists, getList, addItem, addList, setCompletion, updateListName };