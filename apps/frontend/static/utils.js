export const LS_BEARER_TOKEN = "bearer_token"

export function getBearerToken() {
    return localStorage.getItem(LS_BEARER_TOKEN)
}

export async function alertUponException(action, ...args) {
    try {
        await action(...args)
    } catch (error) {
        alert(`error: ${error.message}`)
    }
}

// Given an Authorization header from an HTTP response
// (eg. "Bearer sample_token"), strips the leaded prefix
// to return the token itself (eg. "sample_token")
function _extractTokenFromAuthHeader(token) {
    const prefix = "Bearer "

    return token.substring(prefix.length, token.length)
}

// Takes a response from an authorised endpoint, checks for the Authorization header,
// and stores the bearer token in local storage if the token is present in the header.
// Disregards the status code of the response, as it is assumed that a present 
// token implies a 2xx status and the status code is to be checked by the caller
// of this function. Returns the unaltered response
export function refreshTokenIfNecessary(response) {
    const authHeader = response.headers.get("Authorization")

    if (authHeader) {
        const newToken = _extractTokenFromAuthHeader(authHeader)
        localStorage.setItem(LS_BEARER_TOKEN, newToken)
    }

    return response
}

export async function authorisedFetch(url, options) {
    const t = getBearerToken()
    options = options || {}
    options.headers = options.headers || {}
    options.headers["Authorization"] = `Bearer ${t}`
    
    return await fetch(url, options)
}
