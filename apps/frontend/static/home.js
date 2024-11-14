import * as api from "./api.js"
import * as barkButton from "./barkButton.js"
import * as barkModal from "./barkModal.js"
import { createBark } from "./barkCreation.js"

document.addEventListener("DOMContentLoaded", async function(event) {
    event.preventDefault()
    await loadTimeline()
})

export async function loadTimeline() {
    const timelineDiv = document.getElementById("timeline")
    timelineDiv.innerHTML = ""
    const count = 5
    const offset = 0

    const timelineResponse = await api.getTimeline(count, offset)

    if (!timelineResponse.ok) {
    }

    const barks = (await timelineResponse.json()).barks

    if (barks.length > 0) {
        barks.forEach((b) => {
            timelineDiv.appendChild(createBark(b))
        })
    } else {
    }
}
