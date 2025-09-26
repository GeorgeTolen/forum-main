// Простая система ответов на комментарии
document.addEventListener('DOMContentLoaded', function () {
	// Обработчик ссылок "Ответить"
	document.addEventListener('click', function (e) {
		if (e.target.classList.contains('reply-link')) {
			e.preventDefault()
			const commentId = e.target.getAttribute('data-comment-id')
			const replyForm = document.getElementById('reply-' + commentId)

			if (replyForm) {
				// Показать/скрыть форму ответа
				if (
					replyForm.style.display === 'none' ||
					replyForm.style.display === ''
				) {
					replyForm.style.display = 'block'
					// Фокус на textarea
					const textarea = replyForm.querySelector('textarea')
					if (textarea) {
						textarea.focus()
					}
				} else {
					replyForm.style.display = 'none'
				}
			}
		}
	})

	// Обработчик кнопок "Отмена"
	document.addEventListener('click', function (e) {
		if (e.target.classList.contains('cancel-reply')) {
			e.preventDefault()
			const commentId = e.target.getAttribute('data-comment-id')
			const replyForm = document.getElementById('reply-' + commentId)

			if (replyForm) {
				replyForm.style.display = 'none'
				// Очистить textarea и изображение
				const textarea = replyForm.querySelector('textarea')
				const fileInput = replyForm.querySelector('input[type="file"]')
				const preview = replyForm.querySelector('.image-preview')

				if (textarea) textarea.value = ''
				if (fileInput) fileInput.value = ''
				if (preview) preview.style.display = 'none'
			}
		}
	})

	// Обработчик загрузки изображений
	document.addEventListener('change', function (e) {
		if (e.target.classList.contains('image-upload-input')) {
			const file = e.target.files[0]
			const commentId = e.target.id.replace('reply-image-', '')
			const preview = document.getElementById('reply-preview-' + commentId)
			const previewImg = preview.querySelector('img')

			if (file && file.type.startsWith('image/')) {
				const reader = new FileReader()
				reader.onload = function (e) {
					previewImg.src = e.target.result
					preview.style.display = 'block'
				}
				reader.readAsDataURL(file)
			}
		}
	})

	// Обработчик удаления изображений
	document.addEventListener('click', function (e) {
		if (e.target.classList.contains('remove-image')) {
			e.preventDefault()
			const commentId = e.target.getAttribute('data-comment-id')
			const fileInput = document.getElementById('reply-image-' + commentId)
			const preview = document.getElementById('reply-preview-' + commentId)

			if (fileInput) fileInput.value = ''
			if (preview) preview.style.display = 'none'
		}
	})

	// Обработчик отправки форм ответов
	document.addEventListener('submit', function (e) {
		if (e.target.classList.contains('reply-form-content')) {
			e.preventDefault()

			const formData = new FormData(e.target)

			fetch('/api/comment', {
				method: 'POST',
				body: formData,
			})
				.then(response => {
					if (response.ok) {
						// Очистить форму
						e.target.querySelector('textarea').value = ''
						// Скрыть форму
						e.target.closest('.reply-form').style.display = 'none'
						// Перезагрузить страницу для показа нового комментария
						location.reload()
					} else {
						alert('Ошибка при отправке ответа')
					}
				})
				.catch(error => {
					console.error('Error:', error)
					alert('Ошибка при отправке ответа')
				})
		}
	})

	// Обработчик основной формы комментариев
	const mainCommentForm = document.querySelector(
		'form[action="/api/comment"]:not(.reply-form-content)'
	)
	if (mainCommentForm) {
		mainCommentForm.addEventListener('submit', function (e) {
			e.preventDefault()

			const formData = new FormData(this)

			fetch('/api/comment', {
				method: 'POST',
				body: formData,
				// Не устанавливаем Content-Type, браузер сам установит multipart/form-data
			})
				.then(response => {
					if (response.ok) {
						// Очистить форму
						this.querySelector('textarea').value = ''
						// Перезагрузить страницу
						location.reload()
					} else {
						alert('Ошибка при отправке комментария')
					}
				})
				.catch(error => {
					console.error('Error:', error)
					alert('Ошибка при отправке комментария')
				})
		})
	}

	// Обработчик лайков/дизлайков постов
	const postId = window.postId || 0
	if (postId > 0) {
		const likeLink = document.querySelector(
			'a[href="/post/' + postId + '/like"]'
		)
		const dislikeLink = document.querySelector(
			'a[href="/post/' + postId + '/dislike"]'
		)
		const likesHeader = document.querySelector('h3')

		function updatePostCounts(data) {
			if (likesHeader) {
				likesHeader.textContent =
					'Продвинуто: ' + data.likes + ' · Не нравится: ' + data.dislikes
			}
		}

		if (likeLink) {
			likeLink.addEventListener('click', function (e) {
				e.preventDefault()
				fetch(this.href, { headers: { Accept: 'application/json' } })
					.then(response => response.json())
					.then(data => updatePostCounts(data))
					.catch(() => {})
			})
		}

		if (dislikeLink) {
			dislikeLink.addEventListener('click', function (e) {
				e.preventDefault()
				fetch(this.href, { headers: { Accept: 'application/json' } })
					.then(response => response.json())
					.then(data => updatePostCounts(data))
					.catch(() => {})
			})
		}
	}

	// Обработчик лайков/дизлайков комментариев
	document.querySelectorAll('a[href^="/comment/"]').forEach(function (link) {
		if (/\/comment\/\d+\/(like|dislike)/.test(link.getAttribute('href'))) {
			link.addEventListener('click', function (e) {
				e.preventDefault()
				fetch(this.href, { headers: { Accept: 'application/json' } })
					.then(response => response.json())
					.then(data => {
						const container = this.closest('li')
						if (container) {
							const span = container.querySelector('span')
							if (span) {
								span.textContent =
									'Продвинуто: ' +
									data.likes +
									' · Не нравится: ' +
									data.dislikes
							}
						}
					})
					.catch(() => {})
			})
		}
	})
})
