document.addEventListener('DOMContentLoaded', () => {
    // Navbar toggle for mobile
    const toggler = document.querySelector('.navbar-toggler');
    const collapse = document.querySelector('.collapse.navbar-collapse');
    if (toggler && collapse) {
        toggler.addEventListener('click', () => {
            collapse.classList.toggle('show');
        });
    }

    // Confirm password validation for forms
    const form = document.querySelector('form');
    const password = document.querySelector('#password');
    const confirm = document.querySelector('#confirm');

    if (form && password && confirm) {
        form.addEventListener('submit', (e) => {
            if (password.value !== confirm.value) {
                e.preventDefault();
                alert('Passwords do not match.');
                confirm.focus();
            }
        });
    }
});



