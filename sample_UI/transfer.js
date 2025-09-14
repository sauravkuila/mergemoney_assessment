document.addEventListener('DOMContentLoaded', () => {
  const modal = document.getElementById('transfer-modal')
  const form = document.getElementById('transfer-form')
  const destFields = document.getElementById('dest-fields')
  const modalSource = document.getElementById('modal-source-account')
  const cancelBtn = document.getElementById('cancel-transfer')
  const msg = document.getElementById('transfer-msg')

  function showModal(accountId) {
    modal.style.display = 'flex'
    modal.style.alignItems = 'center'
    modal.style.justifyContent = 'center'
    modalSource.textContent = accountId
    msg.textContent = ''
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
    const accountId = e.detail && e.detail.accountId ? e.detail.accountId : 'unknown'
    showModal(accountId)
  })

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

    // create payload matching server dto.TransferRequest
    const payload = {
      source: {
        sid: parseInt(accountId) || 0,
        currency: sourceCurrency,
        amount: amount
      },
      destination: {
        currency: destCurrency,
        recipient_detail: dest,
        account: dest.account || undefined,
        upi: dest.upi || dest.upi_id || undefined,
        wallet_id: dest.wallet || undefined,
        type: destType
      },
      idempotency_key: `idemp_${Date.now()}_${Math.floor(Math.random()*1000)}`
    }

    // Simulate submit: use fetchWithAuth to call a hypothetical /v1/transfer endpoint (1fa + 2fa required flow omitted)
    try {
      msg.textContent = 'Submitting...'
      const res = await fetchWithAuth('/v1/transfer', { method: 'POST', body: JSON.stringify(payload), headers: { 'Content-Type': 'application/json' } }, '1fa')
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
            const cres = await fetchWithAuth('/v1/transfer/confirm', { method: 'POST', body: JSON.stringify({ transfer_id: d.data.transfer_id, action: 'confirm', twofa_token: code }), headers: { 'Content-Type': 'application/json' } }, '2fa')
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
