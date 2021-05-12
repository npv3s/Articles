let edit_button = document.getElementById('edit-button')
let delete_button = document.getElementById('delete-button')

edit_button.onclick = () => {
    let body = document.getElementById('body')
    let text = body.innerText
    let buttons = document.getElementById('buttons')

    let edit_textarea = document.createElement('textarea')
    edit_textarea.className = 'textarea'
    edit_textarea.value = text

    body.replaceWith(edit_textarea)

    let confirm_button = document.createElement('button')
    confirm_button.id = 'confirm-button'
    confirm_button.className = 'button is-info is-light'
    confirm_button.innerText = 'Изменить'

    let return_buttons = () => {
        delete_button.classList.remove('is-hidden')
        edit_button.classList.remove('is-hidden')
        confirm_button.remove()
        cancel_button.remove()
    }

    confirm_button.onclick = () => {
        text = edit_textarea.value
        fetch('/article/update/' + article_id + '/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'title': "Bye",
                'body': text
            })
        })
            .then((r) => {
                let modal = document.getElementById('modal')
                let modal_box = document.getElementById('modal-box')
                modal_box.innerHTML = ''
                let text = document.createElement('p')
                text.classList.add('mb-3')
                let button = document.createElement('button')
                button.innerText = 'Ok'

                if (r.status === 200) {
                    text.innerText = 'Статья успешно изменена'
                    button.className = 'button is-info is-light'
                    button.onclick = () => {
                        document.location.reload()
                    }
                } else {
                    text.innerText = 'Не удалось изменить статью'
                    button.className = 'button is-info is-light'
                    button.onclick = () => {
                        modal.classList.remove('is-active')
                    }
                }

                modal_box.appendChild(text)
                modal_box.appendChild(button)
                modal.classList.add('is-active')
            })
    }

    let cancel_button = document.createElement('button')
    cancel_button.id = 'cancel-button'
    cancel_button.className = 'button is-danger is-light'
    cancel_button.innerText = 'Отменить'
    cancel_button.onclick = () => {
        edit_textarea.replaceWith(body)
        return_buttons()
    }

    delete_button.classList.add('is-hidden')
    edit_button.classList.add('is-hidden')

    buttons.appendChild(confirm_button)
    buttons.appendChild(cancel_button)
}

delete_button.onclick = () => {
    fetch('/article/delete/' + article_id + '/')
        .then((r) => {
            let modal = document.getElementById('modal')
            let modal_box = document.getElementById('modal-box')
            modal_box.innerHTML = ''
            let text = document.createElement('p')
            text.classList.add('mb-3')
            let button = document.createElement('button')
            button.innerText = 'Ok'

            if (r.status === 200) {
                text.innerText = 'Статья удалена'
                button.className = 'button is-info is-light'
                button.onclick = () => document.location.href = '/'
            } else {
                text.innerText = 'Не удалось удалить статью'
                button.className = 'button is-info is-light'
                button.onclick = () => modal.classList.remove('is-active')
            }

            modal_box.appendChild(text)
            modal_box.appendChild(button)
            modal.classList.add('is-active')
        })
}