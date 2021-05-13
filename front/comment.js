document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('textarea').forEach((t) => t.value = '')
})

function open_reply(id) {
    document.getElementById('reply-' + id).classList.remove('is-hidden')
}

function close_reply(id) {
    document.getElementById('reply-' + id).classList.add('is-hidden')
}

function new_comment(id, root) {
    let textarea = document.getElementById('reply-area-' + id)
    fetch('/comment/new/', {
        method: 'POST',
        body: JSON.stringify({
            "article_id": article_id,
            "root": root,
            "body": textarea.value
        })
    }).then((r) => {
        if (r.status === 200) {
            overlay('Успешно добавлен комментарий',
                () => document.location.reload())
        } else {
            overlay('Не удалось добавить комментарий',
                () => document.querySelectorAll('.reply').forEach((v) => v.classList.add('is-hidden')))
        }
    })
}