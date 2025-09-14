document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('verify-form');
  const msg = document.getElementById('msg');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    msg.textContent = '';
    const country = document.getElementById('country').value.trim();
    const mobile = document.getElementById('mobile').value.trim();
    const otp = document.getElementById('otp').value.trim();
    const reqid = document.getElementById('reqid').value.trim();
    if (!mobile || !otp || !reqid) { msg.textContent = 'All fields required'; return }

    try {
      const res = await fetch('/v1/verifyOTP', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mobile, country_code: country, otp, otp_request_id: reqid })
      })
      const data = await res.json();
      if (!res.ok) {
        msg.textContent = (data.description || 'OTP verify failed');
        return
      }
      // store tokens in localStorage (basic client-side session)
      try {
        if (data && data.data) {
          if (data.data.one_fa_access_token) localStorage.setItem('onefa_token', data.data.one_fa_access_token)
          if (data.data.one_fa_refresh_token) localStorage.setItem('onefa_refresh', data.data.one_fa_refresh_token)
          // if 2fa tokens are present (some flows), store them too
          if (data.data.two_fa_access_token) localStorage.setItem('twofa_token', data.data.two_fa_access_token)
          if (data.data.two_fa_refresh_token) localStorage.setItem('twofa_refresh', data.data.two_fa_refresh_token)
        }
      } catch (e) {
        console.warn('localStorage not available', e);
      }
      // redirect to home screen
      window.location.href = '/ui/home.html';
    } catch (err) {
      msg.textContent = 'Network error';
    }
  })
})
