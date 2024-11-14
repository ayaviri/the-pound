import * as utils from "./utils.js"

// Hits the login endpoint with the given credentials, stores the 
// returned bearer token in local storage if successful, and redirects
// to the home page, throwing upon failure
export async function loginUser(username, password) {
    const registrationForm = { username, password }
    const response = await fetch(
        `http://localhost:8000/login`,
        {
            method: 'POST',
            body: JSON.stringify(registrationForm),
        }
    )

    if (!response.ok) {
        throw new Error("login failed")
    } else {
        const data = await response.json()
        localStorage.setItem(utils.LS_BEARER_TOKEN, data.token)
        window.location.href = "/home.html"
    }
}

export async function registerUser(username, password) {
    const registrationForm = { username, password }
    const response = await fetch(
        `http://localhost:8000/register`,
        {
            method: 'POST',
            body: JSON.stringify(registrationForm),
        }
    )

    if (!response.ok) {
        throw new Error("registration failed")
    }
}

export async function validateToken() {
    const response = await utils.authorisedFetch(`http://localhost:8000/validate`)

    return utils.refreshTokenIfNecessary(response)
}

export async function getTimeline(count, offset) {
    const params = new URLSearchParams({
        count, offset
    })
    const url = `http://localhost:8000/timeline?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getBarks(dogId, count, offset) {
    const params = new URLSearchParams({
        "dog_id": dogId,
        count, offset
    })
    const url = `http://localhost:8000/barks?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function toggleTreat(barkId) {
    const body = {
        "bark_id": barkId
    }
    const response = await utils.authorisedFetch(
        'http://localhost:8000/treat',
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function toggleRebark(barkId) {
    const body = {
        "bark_id": barkId
    }
    const response = await utils.authorisedFetch(
        'http://localhost:8000/rebark',
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getBark(barkId) {
    const params = new URLSearchParams({
        "id": barkId
    })
    const url = `http://localhost:8000/bark?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function postBark(content) {
    const body = { content }
    const response = await utils.authorisedFetch(
        `http://localhost:8000/bark`,
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function postPaw(content, parentBarkId) {
    const body = { content, "bark_id": parentBarkId }
    const response = await utils.authorisedFetch(
        `http://localhost:8000/paw`,
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getThread(barkId) {
    const params = new URLSearchParams({
        "id": barkId
    })
    const url = `http://localhost:8000/thread?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getPaws(barkId) {
    const params = new URLSearchParams({
        "id": barkId
    })
    const url = `http://localhost:8000/paws?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getDogByUsername(username) {
    const params = new URLSearchParams({ username })
    const url = `http://localhost:8000/dog?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getDogById(id) {
    const params = new URLSearchParams({ id })
    const url = `http://localhost:8000/dog?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function doesFollow(dogId) {
    const params = new URLSearchParams({ "id": dogId })
    const url = `http://localhost:8000/does_follow?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function follow(dogId) {
    const body = { "dog_id": dogId }
    const response = await utils.authorisedFetch(
        `http://localhost:8000/follow`,
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function getNotifications(count, offset) {
    const params = new URLSearchParams({
        count, offset
    })
    const url = `http://localhost:8000/notifications?${params.toString()}`
    const response = await utils.authorisedFetch(
        url,
        {
            method: 'GET',
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function approveFollowRequest({ fromDogId, notificationId }) {
    const body = {
        "dog_id": fromDogId,
        "notification_id": notificationId
    }
    const response = await utils.authorisedFetch(
        'http://localhost:8000/approve',
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function rejectFollowRequest({ fromDogId, notificationId }) {
    const body = {
        "dog_id": fromDogId,
        "notification_id": notificationId
    }
    const response = await utils.authorisedFetch(
        'http://localhost:8000/reject',
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}

export async function readNotification(notificationId) {
    const body = {
        "notification_id": notificationId
    }
    const response = await utils.authorisedFetch(
        'http://localhost:8000/notification_read',
        {
            method: 'POST',
            body: JSON.stringify(body)
        }
    )

    return utils.refreshTokenIfNecessary(response)
}
