// create function deleteSpend to delete spend
function deleteSpend(id) {
    // /api/spends/delete/:id delete method
    fetch(`/api/spends/delete/${id}`, {
        method: 'DELETE'
    });

    // reload page
    location.reload();
}