import * as api from "./api.js"
import * as utils from "./utils.js"
import * as barkButton from "./barkButton.js"
import * as barkModal from "./barkModal.js"
import { createBark } from "./barkCreation.js"

document.addEventListener("DOMContentLoaded", async function(event) {
    event.preventDefault()
    const dogId = getDogIdFromQueryString()
    await Promise.all([
        showFollowButton(dogId),
        loadDogProfile(dogId),
        loadDogBarks(dogId)
    ])
})

document.getElementById("follow").addEventListener("click", async function(event) {
    const dogId = getDogIdFromQueryString()
    let response

    await utils.alertUponException(async function() {
        response = await api.follow(dogId)
    })

    if (!response) return

    if (!response.ok) {
    }

    const followResult = await response.json()
    showAndUpdateFollowButtonText(followResult)
    await loadDogBarks(dogId)
})

async function handleFollowButtonClick(event) {
    const dogId = getDogIdFromQueryString()
    let response

    try {
        response = await api.follow(dogId)
    } catch (error) {
        alert(":(")
    }

    if (!response) return

    if (!response.ok) {
        alert("follow failed")
    }

    const followResult = await response.json()
    showAndUpdateFollowButtonText(followResult)
}

function getDogIdFromQueryString() {
    const params = new URLSearchParams(window.location.search)

    return params.get("id")
}

async function showFollowButton(dogId) {
    const doesFollowResponse = await api.doesFollow(dogId)
    const followResult = await doesFollowResponse.json()
    showAndUpdateFollowButtonText(followResult)
}

function showAndUpdateFollowButtonText(followResult) {
    const followButton = document.getElementById("follow")

    if (!followResult["follow_request_exists"]) {
        followButton.textContent = "follow"
    } else {
        if (followResult["is_approved"]) {
            followButton.textContent = "unfollow"
        } else {
            followButton.textContent = "request pending"
        }
    }

    followButton.classList.remove("hidden")
}

async function loadDogProfile(dogId) {
    let response
    
    await utils.alertUponException(async function() {
        response = await api.getDogById(dogId)
    })

    if (!response || !response.ok) return

    const data = await response.json()
    const profileHeader = document.getElementById("profile")
    profileHeader.textContent = `${data.dog.username}'s ${profileHeader.textContent}`
}

async function loadDogBarks(dogId) {
    const barksDiv = document.getElementById("barks")
    barksDiv.innerHTML = ""
    const count = 5
    const offset = 0

    await utils.alertUponException(async function() {
        const barksResponse = await api.getBarks(dogId, count, offset)

        if (barksResponse.status == 403) {
            barksDiv.textContent = "must follow dog to view barks"
        } else if (barksResponse.status == 200) {
            const barks = (await barksResponse.json()).barks

            if (barks.length > 0) {
                barks.forEach((b) => {
                    barksDiv.appendChild(createBark(b))
                })
            } else {
                barksDiv.textContent = "no barks from this dog yet !"
            }
        }
    })
}
