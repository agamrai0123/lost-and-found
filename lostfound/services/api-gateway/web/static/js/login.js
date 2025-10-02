// Optional: Simple client-side form validation
document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    form.addEventListener('submit', (e) => {
        const username = form.querySelector('#username').value.trim();
        const password = form.querySelector('#password').value.trim();

        if (!username || !password) {
            e.preventDefault();
            alert('Please enter both username and password.');
        }
    });
});
