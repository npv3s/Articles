let login_button = document.getElementById('login-button')
let login_input = document.getElementById('login')
let password_input = document.getElementById('password')

login_button.onclick = () => {
    login_button.classList.add('is-loading')
    fetch('/login/', {
        method: 'POST',
        body: JSON.stringify({
            'login': login_input.value,
            'password': password_input.value
        })
    })
        .then((r) => {
            if (r.status === 200) {
                alert('Ok')
            } else {
                alert('Not ok')
            }
        })
        .finally(() => login_button.classList.remove('is-loading'))
}