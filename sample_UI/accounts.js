document.addEventListener('DOMContentLoaded', async () => {
  const list = document.getElementById('accounts-list') || document.getElementById('accountsList')
  if (!list) return

  try {
    const res = await fetchWithAuth('/v1/accounts', { method: 'GET' }, '1fa')

    if (res.status === 401) {
      // token/refresh failed, force logout
      list.textContent = 'Session expired. Redirecting to login...'
      setTimeout(() => forceLogout(), 1000)
      return
    }

    if (!res.ok) {
      const err = await res.json().catch(() => ({ description: 'error' }))
      list.textContent = err.description || 'Failed to load accounts'
      return
    }

    const data = await res.json()
    const accounts = data.data || []
    if (!accounts.length) {
      list.textContent = 'No accounts found.'
      return
    }

    list.innerHTML = ''
    accounts.forEach(acc => {
      const el = document.createElement('div')
      el.className = 'account animate'
      const accId = acc.account_number || acc.wallet_id || acc.upi_id || acc.id || `sid:${acc.sid}`
      el.innerHTML = `<div style="display:flex;justify-content:space-between"><div><strong>${accId}</strong><div class="muted">${acc.account_type || acc.bank_name || ''}</div></div><div>${acc.available_balance || acc.balance || ''}</div></div>`
      el.addEventListener('click', () => {
        // dispatch custom event to open transfer modal and include the full account object
        const ev = new CustomEvent('openTransferModal', { detail: { account: acc } })
        document.dispatchEvent(ev)
      })
      list.appendChild(el)
    })
  } catch (e) {
    list.textContent = 'Network error'
    console.error(e)
  }
})
