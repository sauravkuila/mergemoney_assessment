document.addEventListener('DOMContentLoaded', () => {
  const accountsBtn = document.getElementById('view-accounts')
  const logoutBtn = document.getElementById('logout')

  accountsBtn && accountsBtn.addEventListener('click', () => {
    window.location.href = '/ui/accounts.html'
  })

  logoutBtn && logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('onefa_token')
    localStorage.removeItem('onefa_refresh')
    window.location.href = '/ui/index.html'
  })
})
