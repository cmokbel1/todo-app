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
    } catch(err) {
        console.log(err);
        return err.error;
    }
}

// Add functions
async function addItem(id, item) {
    const data = { name: item, listId: id }
    try {
        const res = await fetch(`/api/todos/${id}/`, 
        {method: 'POST',
         body: JSON.stringify(data)
        })
        const jsonResponse = await res.json();
        return jsonResponse;
    } catch(err) {
        console.log(err);
        return err.error;
    }
}
export { getLists,  getList, addItem };