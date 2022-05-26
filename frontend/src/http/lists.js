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
export { getLists,  getList };