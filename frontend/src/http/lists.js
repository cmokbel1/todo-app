async function getLists() {
    try {
        const res = await fetch('/api/todos');
        const jsonResponse = await res.json();
        console.log(jsonResponse);
        return jsonResponse;
    } catch (err) {
        console.log(err.error);
        return err.error;
    }
}

export { getLists };