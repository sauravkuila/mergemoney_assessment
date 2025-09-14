document.addEventListener('DOMContentLoaded', () => {
  const modal = document.getElementById('transfer-modal')
  const form = document.getElementById('transfer-form')
  const destFields = document.getElementById('dest-fields')
  const modalSource = document.getElementById('modal-source-account')
  const cancelBtn = document.getElementById('cancel-transfer')
  const msg = document.getElementById('transfer-msg')
  let selectedAccount = null // will hold the account object selected from accounts list

  function showModal(account) {
    modal.style.display = 'flex'
    modal.style.alignItems = 'center'
    modal.style.justifyContent = 'center'
    // account may be an object (from accounts API) or a legacy id string
    selectedAccount = account
    const display = (account && account.account_number) || (account && account.wallet_id) || (account && account.upi_id) || account.sid || account || 'unknown'
    modalSource.textContent = display
    msg.textContent = ''
    // show 2FA screen first
    render2FAScreen()
  }

  function hideModal() {
    modal.style.display = 'none'
  }

  // render dest fields based on type
  function renderDestFields(type) {
    destFields.innerHTML = ''
    if (type === 'bank') {
      destFields.innerHTML = `
        <label>Bank Account Number</label>
        <input id="bank-account" />
        <label>IFSC / Routing</label>
        <input id="bank-routing" />
      `
    } else if (type === 'wallet') {
      destFields.innerHTML = `
        <label>Wallet ID / Phone</label>
        <input id="wallet-id" />
      `
    } else {
      destFields.innerHTML = `
        <label>Pickup Location</label>
        <input id="pickup-location" />
        <label>Recipient Name</label>
        <input id="pickup-name" />
      `
    }
  }

  // wire dest type change
  const destTypeSel = document.getElementById('dest-type')
  destTypeSel.addEventListener('change', (e) => renderDestFields(e.target.value))
  renderDestFields(destTypeSel.value)

  cancelBtn.addEventListener('click', hideModal)

  // listen for custom event to open modal with account id
  document.addEventListener('openTransferModal', (e) => {
    // event now contains the account object under detail.account (from updated accounts.js)
    const account = e.detail && (e.detail.account || e.detail.accountId) ? (e.detail.account || e.detail.accountId) : 'unknown'
    showModal(account)
  })

  // Render a 2FA screen that blocks interaction with the transfer form until authenticated
  function render2FAScreen() {
    const container = modal.querySelector('.modal-content')
    container.innerHTML = `
      <h2>2FA Required</h2>
      <p class="muted">Enter your MPIN to derive a 2FA session for transfers from <strong id="twofa-source">${modalSource.textContent}</strong></p>
      <p class="muted">Demo hint: default MPIN for the seeded user is <code>1222</code></p>
      <input id="twofa-input" placeholder="Enter MPIN" />
      <div style="display:flex;gap:8px;margin-top:12px">
        <button id="twofa-auth" class="btn">Authenticate</button>
        <button id="twofa-cancel" class="btn outline">Cancel</button>
      </div>
      <div id="twofa-msg" class="muted"></div>
    `
    // wire buttons
    const authBtn = document.getElementById('twofa-auth')
    const cancelBtn2 = document.getElementById('twofa-cancel')
    const twofaMsg = document.getElementById('twofa-msg')
  authBtn.addEventListener('click', async () => {
      const code = document.getElementById('twofa-input').value.trim()
      if (!code) { twofaMsg.textContent = 'Enter 2FA code'; return }
      // Call server verifyMPIN endpoint using existing 1FA token to get 2FA tokens
      twofaMsg.textContent = 'Authenticating...'
      try {
        // fetchWithAuth requires an access token; for verifyMPIN we use the 1FA token directly in header
        const onefa = localStorage.getItem('onefa_token')
        if (!onefa) { twofaMsg.textContent = '1FA session missing. Please login.'; setTimeout(()=>{ forceLogout() }, 800); return }
        const vres = await fetch('/v1/verifyMPIN', { method: 'POST', headers: { 'Content-Type': 'application/json', 'x-access-token': onefa }, body: JSON.stringify({ mpin: code }) })
        if (vres.status === 401) { twofaMsg.textContent = 'Invalid MPIN or session expired'; return }
        if (!vres.ok) { const err = await vres.json().catch(()=>({description:'verify failed'})); twofaMsg.textContent = err.description || 'Verify MPIN failed'; return }
        const vdata = await vres.json()
        if (vdata && vdata.data && vdata.data.AccessToken) {
          // server returns AccessToken/RefreshToken in DTO VerifyMPINData
          localStorage.setItem('twofa_token', vdata.data.AccessToken)
          localStorage.setItem('twofa_refresh', vdata.data.RefreshToken)
          twofaMsg.textContent = 'Authenticated. You may proceed.'
          // reveal the transfer form now that we have a twofa session
          setTimeout(() => renderTransferForm(), 300)
          return
        }
        // fallback: check lowercased keys (some DTOs use different json names)
        if (vdata && vdata.data && vdata.data.access_token) {
          localStorage.setItem('twofa_token', vdata.data.access_token)
          if (vdata.data.refresh_token) localStorage.setItem('twofa_refresh', vdata.data.refresh_token)
          twofaMsg.textContent = 'Authenticated. You may proceed.'
          setTimeout(() => renderTransferForm(), 300)
          return
        }
        twofaMsg.textContent = 'Unexpected response from server'
      } catch (err) { console.error(err); twofaMsg.textContent = 'Network error during auth' }
    })
    cancelBtn2.addEventListener('click', hideModal)
  }

  // Replace modal content with the transfer form (original markup)
  function renderTransferForm() {
    const container = modal.querySelector('.modal-content')
    container.innerHTML = `
      <h2>Transfer from <span id="modal-source-account">${modalSource.textContent}</span></h2>
      <form id="transfer-form">
        <label>Source Currency</label>
        <select id="source-currency">
          <option value="USD">USD</option>
          <option value="INR">INR</option>
          <option value="EUR">EUR</option>
        </select>

        <label>Amount</label>
        <input id="source-amount" type="number" step="0.01" />

        <label>Destination Type</label>
        <select id="dest-type">
          <option value="bank">Bank Account</option>
          <option value="wallet">Digital Wallet</option>
          <option value="cash">Cash Pickup</option>
        </select>

        <div id="dest-fields"></div>

        <label>Destination Currency</label>
        <select id="dest-currency">
          <option value="USD">USD</option>
          <option value="INR">INR</option>
          <option value="EUR">EUR</option>
        </select>

        <div style="display:flex;gap:8px;margin-top:12px">
          <button type="submit" class="btn">Send</button>
          <button type="button" id="cancel-transfer" class="btn outline">Cancel</button>
        </div>
        <div id="transfer-msg" class="muted"></div>
      </form>
    `
    // re-bind elements and handlers
    const newForm = document.getElementById('transfer-form')
    const newDestFields = document.getElementById('dest-fields')
    const newModalSource = document.getElementById('modal-source-account')
    newModalSource.textContent = modalSource.textContent
    const newCancel = document.getElementById('cancel-transfer')
    newCancel.addEventListener('click', hideModal)

    function newRenderDestFields(type) {
      newDestFields.innerHTML = ''
      if (type === 'bank') {
        newDestFields.innerHTML = `
          <label>Bank Account Number</label>
          <input id="bank-account" />
          <label>IFSC / Routing</label>
          <input id="bank-routing" />
        `
      } else if (type === 'wallet') {
        newDestFields.innerHTML = `
          <label>Wallet ID / Phone</label>
          <input id="wallet-id" />
        `
      } else {
        newDestFields.innerHTML = `
          <label>Pickup Location</label>
          <input id="pickup-location" />
          <label>Recipient Name</label>
          <input id="pickup-name" />
        `
      }
    }

    const destTypeSelect = document.getElementById('dest-type')
    destTypeSelect.addEventListener('change', (e) => newRenderDestFields(e.target.value))
    newRenderDestFields(destTypeSelect.value)

    // rewire submit handling similar to original
    newForm.addEventListener('submit', async (ev) => {
      ev.preventDefault()
      const sourceCurrency = document.getElementById('source-currency').value
      const amount = parseFloat(document.getElementById('source-amount').value)
      const destType = document.getElementById('dest-type').value
      const destCurrency = document.getElementById('dest-currency').value
  // Determine sid: prefer the selectedAccount.sid (if available), otherwise try parsing modalSource text
  let sid = 0
  if (selectedAccount && typeof selectedAccount === 'object' && selectedAccount.sid) sid = parseInt(selectedAccount.sid)
  else sid = parseInt(modalSource.textContent) || 0
      const transferMsg = document.getElementById('transfer-msg')
      transferMsg.textContent = ''
      if (!amount || amount <= 0) { transferMsg.textContent = 'Enter a valid amount'; return }
      let dest = {}
      if (destType === 'bank') {
        dest.account = document.getElementById('bank-account').value.trim()
        dest.routing = document.getElementById('bank-routing').value.trim()
        if (!dest.account) { transferMsg.textContent = 'Enter bank account'; return }
      } else if (destType === 'wallet') {
        dest.wallet = document.getElementById('wallet-id').value.trim()
        if (!dest.wallet) { transferMsg.textContent = 'Enter wallet id'; return }
      } else {
        dest.location = document.getElementById('pickup-location').value.trim()
        dest.name = document.getElementById('pickup-name').value.trim()
        if (!dest.location) { transferMsg.textContent = 'Enter pickup location'; return }
      }

      // build destination to match server DTO (account + swift_code for bank)
      const destination = { currency: destCurrency }
      if (destType === 'bank') {
        destination.account = dest.account
        if (dest.routing) destination.swift_code = dest.routing
      } else if (destType === 'wallet') {
        destination.wallet_id = dest.wallet
      } else {
        destination.recipient_detail = dest
      }

      const payload = {
        source: { sid: sid || 0, currency: sourceCurrency, amount: amount },
        destination: destination,
        idempotency_key: `idemp_${Date.now()}_${Math.floor(Math.random()*1000)}`
      }

      try {
        transferMsg.textContent = 'Submitting...'
        const res = await fetchWithAuth('/v1/transfer', { method: 'POST', body: JSON.stringify(payload), headers: { 'Content-Type': 'application/json' } }, '2fa')
        if (res.status === 401) { transferMsg.textContent = 'Session expired'; setTimeout(()=>{ hideModal(); forceLogout() }, 800); return }
        if (!res.ok) { const d = await res.json().catch(()=>({description:'failed'})); transferMsg.textContent = d.description || 'Transfer failed'; return }
        const d = await res.json()
        // proceed to confirmation flow (reuse earlier logic)
        if (d && d.data && d.data.transfer_id) {
          transferMsg.innerHTML = `Transfer pending confirmation. Transfer ID: <code>${d.data.transfer_id}</code>`
          // show confirm UI
          const twoFAInput = document.createElement('input')
          twoFAInput.placeholder = 'Enter 2FA code'
          twoFAInput.id = 'twofa-code'
          const confirmBtn = document.createElement('button')
          confirmBtn.textContent = 'Confirm Transfer'
          confirmBtn.className = 'btn'
          const cancelBtnEl = document.createElement('button')
          cancelBtnEl.textContent = 'Cancel'
          cancelBtnEl.className = 'btn'
          cancelBtnEl.style.marginLeft = '8px'
          transferMsg.appendChild(document.createElement('br'))
          transferMsg.appendChild(twoFAInput)
          transferMsg.appendChild(confirmBtn)
          transferMsg.appendChild(cancelBtnEl)

          confirmBtn.addEventListener('click', async () => {
            const code = twoFAInput.value.trim()
            if (!code) { alert('Enter 2FA code'); return }
            transferMsg.textContent = 'Confirming...'
            try {
              const cres = await fetchWithAuth('/v1/transfer/confirm', { method: 'POST', body: JSON.stringify({ transfer_id: d.data.transfer_id, action: 'confirm' }), headers: { 'Content-Type': 'application/json' } }, '2fa')
              if (cres.status === 401) { transferMsg.textContent = '2FA session expired. Please login again.'; setTimeout(()=>{ forceLogout() }, 800); return }
              const cd = await cres.json().catch(()=>({description:'confirm failed'}))
              if (!cres.ok) { transferMsg.textContent = cd.description || 'Confirm failed'; return }
              transferMsg.textContent = cd.description || 'Transfer successful'
              setTimeout(()=>{ hideModal() }, 900)
            } catch (err) { console.error(err); transferMsg.textContent = 'Network error during confirm' }
          })

          cancelBtnEl.addEventListener('click', async () => {
            transferMsg.textContent = 'Cancelling...'
            try {
              const cres = await fetchWithAuth('/v1/transfer/confirm', { method: 'POST', body: JSON.stringify({ transfer_id: d.data.transfer_id, action: 'cancel' }), headers: { 'Content-Type': 'application/json' } }, '2fa')
              const cd = await cres.json().catch(()=>({description:'cancel failed'}))
              transferMsg.textContent = cd.description || 'Transfer cancelled'
              setTimeout(()=>{ hideModal() }, 900)
            } catch (err) { console.error(err); transferMsg.textContent = 'Network error during cancel' }
          })
        }
      } catch (err) { console.error(err); transferMsg.textContent = 'Network error' }
    })
  }

  form.addEventListener('submit', async (ev) => {
    ev.preventDefault()
    msg.textContent = ''
    const sourceCurrency = document.getElementById('source-currency').value
    const amount = parseFloat(document.getElementById('source-amount').value)
    const destType = document.getElementById('dest-type').value
    const destCurrency = document.getElementById('dest-currency').value
    const accountId = modalSource.textContent

    if (!amount || amount <= 0) { msg.textContent = 'Enter a valid amount'; return }

    // collect dest details
    let dest = {}
    if (destType === 'bank') {
      dest.account = document.getElementById('bank-account').value.trim()
      dest.routing = document.getElementById('bank-routing').value.trim()
      if (!dest.account) { msg.textContent = 'Enter bank account'; return }
    } else if (destType === 'wallet') {
      dest.wallet = document.getElementById('wallet-id').value.trim()
      if (!dest.wallet) { msg.textContent = 'Enter wallet id'; return }
    } else {
      dest.location = document.getElementById('pickup-location').value.trim()
      dest.name = document.getElementById('pickup-name').value.trim()
      if (!dest.location) { msg.textContent = 'Enter pickup location'; return }
    }

    // build destination similar to modal flow
    const destination = { currency: destCurrency }
    if (destType === 'bank') {
      destination.account = dest.account
      if (dest.routing) destination.swift_code = dest.routing
    } else if (destType === 'wallet') {
      destination.wallet_id = dest.wallet
    } else {
      destination.recipient_detail = dest
    }

    const payload = {
      source: {
        sid: parseInt(accountId) || 0,
        currency: sourceCurrency,
        amount: amount
      },
      destination: destination,
      idempotency_key: `idemp_${Date.now()}_${Math.floor(Math.random()*1000)}`
    }

    // Simulate submit: use fetchWithAuth to call a hypothetical /v1/transfer endpoint (1fa + 2fa required flow omitted)
    try {
      msg.textContent = 'Submitting...'
  const res = await fetchWithAuth('/v1/transfer', { method: 'POST', body: JSON.stringify(payload), headers: { 'Content-Type': 'application/json' } }, '2fa')
      if (res.status === 401) {
        msg.textContent = 'Session expired, please login again.'
        setTimeout(()=>{ hideModal(); forceLogout() }, 800)
        return
      }
      if (!res.ok) {
        const d = await res.json().catch(()=>({description:'failed'}))
        msg.textContent = d.description || 'Transfer failed'
        return
      }
      const d = await res.json()
      // After creating transfer, if server accepted, prompt for 2FA confirmation
      if (d && d.data && d.data.transfer_id) {
        msg.innerHTML = `Transfer pending confirmation. Transfer ID: <code>${d.data.transfer_id}</code>`
        // show a simple 2FA input
        const twoFAInput = document.createElement('input')
        twoFAInput.placeholder = 'Enter 2FA code'
        twoFAInput.id = 'twofa-code'
        const confirmBtn = document.createElement('button')
        confirmBtn.textContent = 'Confirm Transfer'
        confirmBtn.className = 'btn'
        const cancelBtn = document.createElement('button')
        cancelBtn.textContent = 'Cancel'
        cancelBtn.className = 'btn'
        cancelBtn.style.marginLeft = '8px'
        msg.appendChild(document.createElement('br'))
        msg.appendChild(twoFAInput)
        msg.appendChild(confirmBtn)
        msg.appendChild(cancelBtn)

        confirmBtn.addEventListener('click', async () => {
          const code = twoFAInput.value.trim()
          if (!code) { alert('Enter 2FA code'); return }
          // call confirm endpoint under 2fa
          msg.textContent = 'Confirming...'
          try {
            const cres = await fetchWithAuth('/v1/transfer/confirm', { method: 'POST', body: JSON.stringify({ transfer_id: d.data.transfer_id, action: 'confirm' }), headers: { 'Content-Type': 'application/json' } }, '2fa')
            if (cres.status === 401) { msg.textContent = '2FA session expired. Please login again.'; setTimeout(()=>{ forceLogout() }, 800); return }
            const cd = await cres.json().catch(()=>({description:'confirm failed'}))
            if (!cres.ok) { msg.textContent = cd.description || 'Confirm failed'; return }
            msg.textContent = cd.description || 'Transfer successful'
            setTimeout(()=>{ hideModal() }, 900)
          } catch (err) {
            console.error(err)
            msg.textContent = 'Network error during confirm'
          }
        })

        cancelBtn.addEventListener('click', async () => {
          // cancel the pending transfer
          msg.textContent = 'Cancelling...'
          try {
            const cres = await fetchWithAuth('/v1/transfer/confirm', { method: 'POST', body: JSON.stringify({ transfer_id: d.data.transfer_id, action: 'cancel' }), headers: { 'Content-Type': 'application/json' } }, '2fa')
            const cd = await cres.json().catch(()=>({description:'cancel failed'}))
            msg.textContent = cd.description || 'Transfer cancelled'
            setTimeout(()=>{ hideModal() }, 900)
          } catch (err) {
            console.error(err)
            msg.textContent = 'Network error during cancel'
          }
        })
        return
      }
      msg.textContent = d.description || 'Transfer submitted successfully'
      setTimeout(()=>{ hideModal() }, 900)
    } catch (err) {
      console.error(err)
      msg.textContent = 'Network error'
    }
  })
})
