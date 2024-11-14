import * as api from "./api.js"

class BarkModal extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
<div class="bark_modal hidden">
    <div id="bark_modal_content">
        <h2>take your thought to the pound</h2>
        <textarea id="bark_textarea" placeholder="bark here..."></textarea>
        <button id="bark_submission">bark!</button> <!-- TODO: Replace this with a paw image or something -->
    </div>
</div>
    `
  }
}

customElements.define("bark-modal", BarkModal)

document.getElementById("bark_submission").addEventListener("click", async function() {
    const content = document.getElementById("bark_textarea").value 

    if (content == "") {
        alert("mute barks are not allowed !")
        return
    }

    const barkModal = document.querySelector(".bark_modal")
    const postAsBark = barkModal.getAttribute("data-bark-type") == "bark"

    try {
        let response 

        if (postAsBark) {
            response = await api.postBark(content)
        } else {
            const parentBarkId = barkModal.getAttribute("data-parent-bark-id")
            response = await api.postPaw(content, parentBarkId)
        }

        if (response.ok) {
            closeModal()
            // TODO: Is there a way to know what specifically to reload given what page 
            // this is on ? In order to avoid the reloading the whole page ?
            location.reload()
            // await home.loadTimeline()
        } else {
            alert("could not post bark :(")
        }
    } catch (error) {
        alert("could not connect to server :(")
    }
})

document.querySelector(".bark_modal").addEventListener("click", function() {
    if (event.target === document.querySelector(".bark_modal")) {
        closeModal()
    }
}) 

function closeModal() {
    document.querySelector(".bark_modal").classList.add("hidden")
}

// Removes the hidden class attribute and sets the type of bark to 
// sent. If _postAsBark_ is false, then the bark will be posted as a paw to
// the bark with the given _parentBarkId_. Otherwise, the bark will be posted
// as normal and the given _parentBarkId_ will be ignored
export function openBarkModal({ postAsBark, parentBarkId }) {
    const barkModal = document.querySelector(".bark_modal")
    barkModal.classList.remove("hidden")

    if (postAsBark) {
        barkModal.setAttribute("data-bark-type", "bark")
        barkModal.removeAttribute("data-parent-bark-id")
    } else {
        barkModal.setAttribute("data-bark-type", "paw")
        barkModal.setAttribute("data-parent-bark-id", parentBarkId)
    }
}
