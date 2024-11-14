import * as api from "./api.js"
import * as utils from "./utils.js"
import * as redirect from "./redirect.js"
import * as barkButton from "./barkButton.js"
import * as barkModal from "./barkModal.js"

document.getElementById("search").addEventListener("submit", async function(event) {
    event.preventDefault()
    await utils.alertUponException(async function() {
        const username = document.getElementById("search_field").value
        const response = await api.getDogByUsername(username)

        if (!response.ok) {
            console.log("request failed ?")
        }

        const dog = (await response.json()).dog
        appendSearchResult(dog)
    })
})

function appendSearchResult(dog) {
    const container = document.createElement("div")
    container.textContent = `@${dog.username}`

    container.addEventListener("click", function() {
        redirect.toProfilePage(dog.id)
    })

    document.getElementById("search_result").appendChild(container)
}
