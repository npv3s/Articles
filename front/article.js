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
        fetch('/article/update/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'article_id': article_id,
                'title': "Bye",
                'body': text
            })
        })
            .then((r) => {
                if (r.status === 200) {
                    overlay('Статья успешно изменена',
                        () => document.location.reload()
                    )
                } else {
                    overlay('Не удалось изменить статью',
                        (modal) => modal.classList.remove('is-active')
                    )
                }
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
    fetch('/article/delete/', {
        method: 'POST',
        body: JSON.stringify({
            "article_id": article_id
        })
    })
        .then((r) => {
            if (r.status === 200) {
                overlay('Статья удалена',
                    () => document.location.href = '/'
                )
            } else {
                overlay('Не удалось удалить статью',
                    (modal) => modal.classList.remove('is-active')
                )
            }
        })
}