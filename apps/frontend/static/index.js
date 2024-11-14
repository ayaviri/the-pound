import * as redirect from "./redirect.js"

document.addEventListener("DOMContentLoaded", async function(event) {
    await redirect.authorisedUsersToHomePage()
})

document.getElementById("register").addEventListener("click", function(event) {
    window.location.href = "/register.html"
})

document.getElementById("login").addEventListener("click", function(event) {
    window.location.href = "/login.html"
})
