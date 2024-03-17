const searchInput = document.getElementById('searchInput');
const searchResults = document.getElementById('searchResults');

searchInput.addEventListener('input', function() {
    const query = this.value;
    if (query.trim() !== '') {
        fetch(`/suggest?query=${query}`)
            .then(response => response.json())
            .then(data => {
                displayResults(data);
            })
            .catch(error => console.error('Error fetching suggestions:', error));
    } else {
        searchResults.innerHTML = '';
    }
});

function displayResults(results) {
    searchResults.innerHTML = '';
    results.forEach(result => {
        const div = document.createElement('div');
        div.textContent = result;
        searchResults.appendChild(div);
    });
}