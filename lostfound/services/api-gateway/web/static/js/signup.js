document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    const password = document.querySelector('#password');
    const confirm = document.querySelector('#confirm');

    form.addEventListener('submit', (e) => {
        if (password.value !== confirm.value) {
            e.preventDefault();
            alert('Passwords do not match.');
            confirm.focus();
        }
    });
});
