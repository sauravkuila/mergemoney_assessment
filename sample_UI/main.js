document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('login-form');
  const msg = document.getElementById('msg');

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    msg.textContent = '';
    const country = document.getElementById('country').value.trim();
    const mobile = document.getElementById('mobile').value.trim();
    if (!mobile) { msg.textContent = 'Enter mobile'; return }

    try {
      const url = `/v1/generateOTP?mobile=${encodeURIComponent(mobile)}&country_code=${encodeURIComponent(country)}`;
      const res = await fetch(url, { method: 'GET' });
      const data = await res.json();
      if (!res.ok) {
        msg.textContent = (data.description || 'Failed to generate OTP');
        return
      }
      // show the OTPRequestId and OTP for demo environments
      msg.innerHTML = `OTP generated. Request ID: <code>${data.data.otp_request_id}</code><br/>OTP (demo): <strong>${data.data.otp}</strong><br/><a href="verify.html">Go verify</a>`;
    } catch (err) {
      msg.textContent = 'Network error';
    }
  })
})
