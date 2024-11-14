import * as utils from "./utils.js"
import * as api from "./api.js"

export async function authorisedUsersToHomePage() {
    const response = await api.validateToken()

    if (response.ok) {
        window.location.href = "/home.html"
    }
}

export function toProfilePage(dogId) {
    window.location.href = `/profile.html?id=${dogId}`
}

export function toBark(barkId) {
    window.location.href = `/bark.html?id=${barkId}`
}
