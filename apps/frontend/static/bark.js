import * as api from "./api.js"
import * as barkButton from "./barkButton.js"
import * as barkModal from "./barkModal.js"
import { createBark } from "./barkCreation.js"

document.addEventListener("DOMContentLoaded", async function(event) {
    event.preventDefault()
    const barkId = getBarkIdFromQueryString()
    await loadThread(barkId)
})

function getBarkIdFromQueryString() {
    const params = new URLSearchParams(window.location.search)

    return params.get("id")
}

async function loadThread(barkId) {
    const threadDiv = document.getElementById("thread")

    try {
        const [threadResponse, pawsResponse] = await Promise.all(
            [api.getThread(barkId), api.getPaws(barkId)]
        )

        if (!threadResponse.ok || !pawsResponse.ok) {
            alert("could not get thread or paws :(")
        }
    
        const threadBarks = (await threadResponse.json()).barks || []
        const paws = (await pawsResponse.json()).paws || []
        const combined = threadBarks.concat(paws)

        if (combined.length > 0) {
            combined.forEach((b) => {
                threadDiv.appendChild(createBark(b))
            })
        }
    } catch (error) {
        console.log(error)
        alert("could not connect to server :(")
    }
}
