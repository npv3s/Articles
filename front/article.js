let edit_button = document.getElementById('edit-button')
let delete_button = document.getElementById('delete-button')

edit_button.onclick = () => {
    let body = document.getElementById('body')
    let text = body.innerText

    let field = document.createElement('div')
    field.className = 'field'

    let control = document.createElement('div')
    control.className = 'control'

    let edit = document.createElement('textarea')
    edit.className = 'textarea'
    edit.value = text

    control.appendChild(edit)
    field.appendChild(control)
    body.replaceWith(field)
}

delete_button.onclick = () => {
    fetch('/article/delete/' + article_id + '/')
        .then((r) => {
            let modal = document.getElementById('modal')
            let modal_box = document.getElementById('modal-box')
            modal_box.innerHTML = ''
            let text = document.createElement('p')
            let buttons = document.createElement('div')
            buttons.className = 'field is-grouped is-grouped-right'
            let button = document.createElement('button')
            button.innerText = 'Ok'

            if (r.status === 200) {
                text.innerText = 'Статья удалена'
                button.className = 'button is-info is-light is-pulled-right'
                button.onclick = () => document.location.href = '/'
            } else {
                text.innerText = 'Не удалось удалить статью'
                button.className = 'button is-info is-light'
                button.onclick = () => modal.classList.remove('is-active')
            }

            modal_box.appendChild(text)
            buttons.appendChild(button)
            modal_box.appendChild(buttons)
            modal.classList.add('is-active')
        })
}