document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm');
    const loginButton = document.getElementById('loginButton');
    const togglePassword = document.getElementById('togglePassword');
    const passwordInput = document.getElementById('password');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        loginButton.classList.add('loading');
        loginButton.textContent = '';

        // Simulate login process
        await new Promise(resolve => setTimeout(resolve, 2000));

        loginButton.classList.remove('loading');
        loginButton.textContent = 'Login';
        alert('Login successful!');
    });

    // Toggle password visibility
    togglePassword.addEventListener('click', () => {
        const type = passwordInput.type === 'password' ? 'text' : 'password';
        passwordInput.type = type;

        // Change icon based on visibility
        togglePassword.textContent = type === 'password' ? '👁️' : '🕶️';
    });

    // Add animation to input fields
    const inputs = document.querySelectorAll('input');
    inputs.forEach(input => {
        input.addEventListener('focus', () => {
            input.style.transform = 'scale(1.05)';
        });
        input.addEventListener('blur', () => {
            input.style.transform = 'scale(1)';
        });
    });
});
