import * as authForm from "./authForm.js"
import * as api from "./api.js"
import * as utils from "./utils.js"

document.getElementById("auth").addEventListener("submit", async function(event) {
    event.preventDefault()
    const username = document.getElementById("username").value
    const password = document.getElementById("password").value

    utils.alertUponException(api.registerUser, username, password)
})
