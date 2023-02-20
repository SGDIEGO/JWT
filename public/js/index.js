let navBar = document.querySelector(".nav-bar");
let info = `<nav>
                <li>
                    <a href="/">Home</a>
                    <a href="/users">Users</a>
                    <a href="/register">Register</a>
                    <a href="/login">Login</a>
                    <a href="/logout">Logout</a>
                </li>
            </nav>`

// fetch('/www/shared/header.html')
// .then(file => file.text())
// .then(content => console.log(content))

document.addEventListener("DOMContentLoaded", () => {
    navBar.innerHTML = info
});
