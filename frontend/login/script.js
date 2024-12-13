document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("loginForm");
    const loginButton = document.getElementById("loginButton");
    const togglePassword = document.getElementById("togglePassword");
    const passwordInput = document.getElementById("password");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const curPass = passwordInput.value;
        const curLoginInput = document.getElementById("login").value;
        loginButton.classList.add("loading");
        loginButton.textContent = "";
        // Simulate login process

        // alert("Login successful!");
        let objData = {
            login: curLoginInput,
            password: curPass,
        };

        fetch("http://localhost:8080/auth", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(objData), // Ваши данные для отправки
        })
            .then((res) => res.json())
            .then((data) => {
                loginButton.classList.remove("loading");
                loginButton.textContent = "Login";
                if (data.error) {
                    alert("none");
                } else {
                    document.cookie = `token=${data.token}`;
                    window.location.href = "http://localhost/";
                }
            })
            .catch((error) => alert(error));
    });

    // Toggle password visibility
    togglePassword.addEventListener("click", () => {
        const type = passwordInput.type === "password" ? "text" : "password";
        passwordInput.type = type;

        // Change icon based on visibility
        togglePassword.textContent = type === "password" ? "👁️" : "🕶️";
    });

    // Add animation to input fields
    const inputs = document.querySelectorAll("input");
    inputs.forEach((input) => {
        input.addEventListener("focus", () => {
            input.style.transform = "scale(1.05)";
        });
        input.addEventListener("blur", () => {
            input.style.transform = "scale(1)";
        });
    });
});
