import * as api from "./api.js"
import * as utils from "./utils.js"
import * as redirect from "./redirect.js"
import * as barkButton from "./barkButton.js"
import * as barkModal from "./barkModal.js"

document.addEventListener("DOMContentLoaded", async function(event) {
    event.preventDefault()
    await loadNotifications()
})

export async function loadNotifications() {
    const notificationsDiv = document.getElementById("notifications")
    notificationsDiv.innerHTML = ""
    const count = 5
    const offset = 0

    const notificationsResponse = await api.getNotifications(count, offset)

    if (!notificationsResponse.ok) {
    }

    const notifications = (await notificationsResponse.json()).notifications

    if (notifications.length > 0) {
        notifications.forEach((n) => {
            notificationsDiv.appendChild(createNotification(n))
        })
    } else {
        notificationsDiv.textContent = "no notifications for this dog yet !"
    }
}

function createNotification(notification) {
    let container

    switch (notification.type) {
        case "treat":
            container = createTreatNotification(notification) 
            break;
        case "rebark":
            container = createRebarkNotification(notification) 
            break;
        case "paw":
            container = createPawNotification(notification) 
            break;
        case "follow":
            container = createFollowNotification(notification) 
            break;
    }

    container.setAttribute("class", "notification")
    return container
}

function createNotificationHelper(imageSrc, message, onClickHandler, actionContainer) {
    const container = document.createElement("div")
    
    const contentContainer = document.createElement("div")
    contentContainer.setAttribute("class", "notification_content")

    const image = document.createElement("img")
    image.setAttribute("class", "notification_icon")
    image.src = imageSrc

    const text = document.createElement("p")
    text.textContent = message
    text.setAttribute("class", "notification_content")

    contentContainer.appendChild(image)
    contentContainer.appendChild(text)

    container.appendChild(contentContainer)
    
    if (actionContainer) {
        container.appendChild(actionContainer)
    }

    container.addEventListener("click", onClickHandler)

    return container
}

function createTreatNotification(notification) {
    const message = `@${notification.payload.from_dog_username} liked your bark <3`
    const onClickHandler = async function() {
        await api.readNotification(notification.id)
        redirect.toBark(notification.payload.bark_id)
    }
    return createNotificationHelper("treat.png", message, onClickHandler)
}

function createRebarkNotification(notification) {
    const message = `@${notification.payload.from_dog_username} rebarked you ♲`
    const onClickHandler = async function() {
        await api.readNotification(notification.id)
        redirect.toBark(notification.payload.bark_id)
    }
    return createNotificationHelper("rebark.png", message, onClickHandler)
}

function createPawNotification(notification) {
    const message = `@${notification.payload.from_dog_username} pawed: ${notification.payload.bark}`
    const onClickHandler = async function() {
        await api.readNotification(notification.id)
        redirect.toBark(notification.payload.bark_id)
    }
    return createNotificationHelper("paw.jpg", message, onClickHandler)
}

function createFollowNotification(notification) {
    const message = (
        notification.payload.is_approved ? 
        `@${notification.payload.from_dog_username} has followed you` :
        `@${notification.payload.from_dog_username} has requested to follow you`
    )
    const onClickHandler = function() {
        if (
            event.target.classList.contains("notification_content") || 
            event.target.classList.contains("notification_icon") || 
            event.target.classList.contains("notification")
        ) {
            redirect.toProfilePage(notification.payload.from_dog_id) 
        }
    }
    const actionContainer = document.createElement("div")
    actionContainer.setAttribute("class", "notification_action")

    if (!notification.payload.is_approved) {
        actionContainer.appendChild(createFollowApproveButton(notification))
        actionContainer.appendChild(createFollowRejectButton(notification))
    }

    return createNotificationHelper(
        "follow.svg", 
        message, 
        onClickHandler,
        actionContainer
    )
}

function createFollowApproveButton(notification) {
    const approveImage = document.createElement("img")
    approveImage.src = "checkmark.png"
    approveImage.setAttribute("class", "notification_action_icon")
    approveImage.addEventListener("click", async function(event) {
        let response

        await utils.alertUponException(async function() {
            response = await api.approveFollowRequest({
                fromDogId: notification.payload.from_dog_id,
                notificationId: notification.id
            })
        })

        if (!response || !response.ok) {
            return
        }

        location.reload()
    })

    return approveImage
}

function createFollowRejectButton(notification) {
    const rejectImage = document.createElement("img")
    rejectImage.src = "x.png"
    rejectImage.setAttribute("class", "notification_action_icon")
    rejectImage.addEventListener("click", async function(event) {
        let response

        await utils.alertUponException(async function() {
            response = await api.rejectFollowRequest({
                fromDogId: notification.payload.from_dog_id,
                notificationId: notification.id
            })
        })

        if (!response || !response.ok) {
            return
        }

        location.reload()
    })

    return rejectImage
}
