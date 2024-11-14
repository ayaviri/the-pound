// Converts an ISO 8601 formatted datetime string into
// one formatted as eg. 3m ago
export function humaniseTime(datetimeString) {
    const givenTime = new Date(datetimeString)
    const currentTime = new Date()
    const diffInSeconds = Math.floor((currentTime - givenTime) / 1000)

    if (diffInSeconds < 60) {
        return `${diffInSeconds}s ago`
    }

    const diffInMinutes = Math.floor(diffInSeconds / 60)

    if (diffInMinutes < 60) {
        return `${diffInMinutes}m ago`
    }

    const diffInHours = Math.floor(diffInMinutes / 60)

    if (diffInHours < 24) {
        return `${diffInHours}h ago`
    }

    const diffInDays = Math.floor(diffInHours / 24)

    if (diffInDays < 7) {
        return `${diffInDays}d ago`
    }

    const diffInWeeks = Math.floor(diffInDays / 7)

    return `${diffInWeeks}w ago`
}
