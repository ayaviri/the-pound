import * as date from "./date.js"
import * as api from "./api.js"
import * as redirect from "./redirect.js"
import { openBarkModal } from "./barkModal.js"

// Creates an element to render the given bark
export function createBark(bark) {
    const container = document.createElement("div")
    container.setAttribute("class", "bark")
    container.appendChild(createBarkHeader(bark))
    container.appendChild(createBarkBody(bark))
    container.appendChild(createBarkFooter(bark))

    container.addEventListener("click", function(event) {
        if (
            event.target.classList.contains("bark_body") ||
            event.target.classList.contains("bark_footer") ||
            event.target.classList.contains("bark_header") ||
            event.target.classList.contains("bark")
        ) {
            redirect.toBark(bark.id)
        }
    })

    return container
}

function createBarkHeader(bark) {
    const header = document.createElement("div")

    const username = document.createElement("p")
    username.textContent = "@" + bark.dog_username
    username.addEventListener("click", function(event) {
        window.location.href = `/profile.html?id=${bark.dog_id}`
    })

    const creationDate = document.createElement("p")
    creationDate.textContent = date.humaniseTime(bark.created_at)

    header.appendChild(username)
    header.appendChild(creationDate)
    header.setAttribute("class", "bark_header")

    return header
}

function createBarkBody(bark) {
    const body = document.createElement("div")
    body.setAttribute("class", "bark_body")
    body.textContent = bark.bark

    return body
}

function createBarkFooter(bark) {
    const footer = document.createElement("div")
    footer.setAttribute("class", "bark_footer")
    footer.appendChild(createTreatButton(bark))
    footer.appendChild(createRebarkButton(bark))
    footer.appendChild(createPawButton(bark))
    const treatButton = document.createElement("button")
    const rebarkButton = document.createElement("button")

    return footer
}

function createTreatButton(bark) {
    const container = document.createElement("div")
    container.setAttribute("class", "bark_interaction")
    const button = document.createElement("button")
    const count = document.createElement("p")
    count.setAttribute("class", "bark_interaction_count")
    button.textContent = "treat"
    count.textContent = bark.treat_count

    // TODO: This needs to be a toggle, so the bark response somehow
    // needs to know if the requesting individual has given any returned
    // bark a treat or a rebark
    button.addEventListener("click", async function() {
        const treatResponse = await api.toggleTreat(bark.id)
        const barkResponse = await api.getBark(bark.id) 
        const updatedBark = (await barkResponse.json()).bark
        count.textContent = updatedBark.treat_count
    })

    container.appendChild(button)
    container.appendChild(count)

    return container
}

function createRebarkButton(bark) {
    const container = document.createElement("div")
    container.setAttribute("class", "bark_interaction")
    const button = document.createElement("button")
    const count = document.createElement("p")
    count.setAttribute("class", "bark_interaction_count")
    button.textContent = "rebark"
    count.textContent = bark.rebark_count

    button.addEventListener("click", async function() {
        const rebarkResponse = await api.toggleRebark(bark.id)
        const barkResponse = await api.getBark(bark.id) 
        const updatedBark = (await barkResponse.json()).bark
        count.textContent = updatedBark.rebark_count
    })

    container.appendChild(button)
    container.appendChild(count)

    return container
}

function createPawButton(bark) {
    const container = document.createElement("div")
    container.setAttribute("class", "bark_interaction")
    const button = document.createElement("button")
    const count = document.createElement("p")
    count.setAttribute("class", "bark_interaction_count")
    button.textContent = "paw"
    count.textContent = bark.paw_count

    button.addEventListener("click", async function() {
        openBarkModal({ postAsBark: false, parentBarkId: bark.id })
    })

    container.appendChild(button)
    container.appendChild(count)

    return container
}
