import { openBarkModal } from "./barkModal.js"

class BarkButton extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
<div id="bark_button">
    <img src="bark_button.png">
</div>
    `
  }
}

customElements.define("bark-button", BarkButton)

document.getElementById("bark_button").addEventListener("click", function() {
    openBarkModal({ postAsBark: true })
})
