document.getElementById('save').onclick = () => {
    let ta = document.querySelector('textarea')
    let title = document.getElementById('title')
    fetch('/article/new/', {
        method: 'POST',
        body: JSON.stringify({
            "title": title.value,
            "body": ta.value
        })
    }).then((r) => {
        if (r.status === 200) {
            r.json().then((j) => {
                document.location = '/article/' + j.article_id + '/'
            })
        } else {
            overlay("Не удалось создать статью",
                (modal) => modal.classList.remove('is-active'))
        }
    })
}