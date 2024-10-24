// Handle Add Pack Size form submission.
document.getElementById('add-pack-form').addEventListener('submit', function (e) {
    e.preventDefault();
    const size = document.getElementById('size').value; // ID matches the HTML input for size

    fetch('/api/v1/packs', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ size: parseInt(size) })
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('add-pack-message').textContent = data.message || data.error; // ID matches the message div
        })
        .catch(error => {
            document.getElementById('add-pack-message').textContent = 'Error: ' + error;
        });
});

// Handle Delete Pack Size form submission.
document.getElementById('delete-pack-form').addEventListener('submit', function (e) {
    e.preventDefault();
    const size = document.getElementById('delete-size').value; // ID matches the HTML input for delete size

    fetch('/api/v1/packs', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ size: parseInt(size) })
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('delete-pack-message').textContent = data.message || data.error; // ID matches the message div
        })
        .catch(error => {
            document.getElementById('delete-pack-message').textContent = 'Error: ' + error;
        });
});

// Handle Calculate Packs form submission.
document.getElementById('calculate-pack-form').addEventListener('submit', function (e) {
    e.preventDefault();
    const items = document.getElementById('items').value;

    fetch(`/api/v1/calculate?quantity=${items}`, {
        method: 'GET',
    })
        .then(response => response.json())
        .then(data => {
            let result = 'Packs required: ';
            for (const [size, count] of Object.entries(data.packs)) {
                result += `${count} x ${size} pack(s), `;
            }
            document.getElementById('calculate-result').textContent = result.slice(0, -2);
        })
        .catch(error => {
            document.getElementById('calculate-result').textContent = 'Error: ' + error;
        });
});
