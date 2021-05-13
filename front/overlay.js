function overlay(message, button_function) {
    let modal = document.getElementById('modal')
    let modal_box = document.getElementById('modal-box')
    let text = document.createElement('p')
    let button = document.createElement('button')

    modal_box.innerHTML = ''
    text.classList.add('mb-3')
    text.innerText = message
    button.innerText = 'Ok'
    button.className = 'button is-info is-light'
    button.onclick = () => button_function(modal)

    modal_box.appendChild(text)
    modal_box.appendChild(button)
    modal.classList.add('is-active')
}