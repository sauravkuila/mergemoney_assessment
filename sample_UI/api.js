/* Shared API helper for authenticated fetches with automatic refresh retry.
   Supports '1fa' and '2fa' token sets stored in localStorage as:
     - onefa_token, onefa_refresh
     - twofa_token, twofa_refresh
*/

async function fetchWithAuth(path, opts = {}, authType = '1fa') {
  // normalize
  authType = (authType === '2fa') ? '2fa' : '1fa'

  const tokenKey = (authType === '1fa') ? 'onefa_token' : 'twofa_token'
  const refreshKey = (authType === '1fa') ? 'onefa_refresh' : 'twofa_refresh'

  let accessToken = localStorage.getItem(tokenKey)
  const refreshToken = localStorage.getItem(refreshKey)

  if (!accessToken) {
    return Promise.reject({unauthenticated: true})
  }

  // ensure headers
  opts.headers = opts.headers || {}
  if (!opts.headers['Content-Type'] && !(opts.body instanceof FormData)) {
    opts.headers['Content-Type'] = 'application/json'
  }
  // server expects token in x-access-token header (no 'Bearer ' prefix)
  opts.headers['x-access-token'] = accessToken

  // first attempt
  let res = await fetch(path, opts)
  if (res.status !== 401) return res

  // got 401; attempt refresh once
  if (!refreshToken) {
    // no refresh token available
    localStorage.removeItem(tokenKey)
    localStorage.removeItem(refreshKey)
    return res
  }

  try {
    const refreshPath = (authType === '1fa') ? '/v1/1fa/refresh' : '/v1/2fa/refresh'
    const refreshBody = JSON.stringify({ refresh_token: refreshToken })
    // refresh requires Authorization header with (possibly expired) access token
  const rres = await fetch(refreshPath, { method: 'POST', headers: { 'Content-Type': 'application/json', 'x-access-token': accessToken }, body: refreshBody })
    if (!rres.ok) {
      // refresh failed -> clear tokens and return original 401
      localStorage.removeItem(tokenKey)
      localStorage.removeItem(refreshKey)
      return res
    }
    const rdata = await rres.json()
    // store new access token depending on response shape
    if (authType === '1fa' && rdata && rdata.data && rdata.data.one_fa_access_token) {
      const newAccess = rdata.data.one_fa_access_token
      localStorage.setItem('onefa_token', newAccess)
      accessToken = newAccess
    } else if (authType === '2fa' && rdata && rdata.data && rdata.data.access_token) {
      const newAccess = rdata.data.access_token
      localStorage.setItem('twofa_token', newAccess)
      accessToken = newAccess
    } else {
      // unexpected response, clear and return original 401
      localStorage.removeItem(tokenKey)
      localStorage.removeItem(refreshKey)
      return res
    }

    // retry original request once with updated access token
  opts.headers['x-access-token'] = accessToken
    const retry = await fetch(path, opts)
    return retry
  } catch (err) {
    // network or unexpected error during refresh; clear tokens and return original 401 response
    localStorage.removeItem(tokenKey)
    localStorage.removeItem(refreshKey)
    return res
  }
}

// small helper to logout and redirect to login UI
function forceLogout() {
  localStorage.removeItem('onefa_token')
  localStorage.removeItem('onefa_refresh')
  localStorage.removeItem('twofa_token')
  localStorage.removeItem('twofa_refresh')
  window.location.href = '/ui/index.html'
}
