import React, { useState } from 'react'
import SearchBar from '../components/SearchBar'

export default function PostsPage() {
	const [posts] = useState([
		{ id: 1, content: 'Изучаю React' },
		{ id: 2, content: 'Пишу смарт-контракт' },
		{ id: 3, content: 'Готовлюсь к экзамену' },
	])
	const [filtered, setFiltered] = useState(posts)

	const handleSearch = query => {
		setFiltered(
			posts.filter(p => p.content.toLowerCase().includes(query.toLowerCase()))
		)
	}

	return (
		<div>
			<h1 className='text-xl font-bold mb-4'>Посты</h1>
			<SearchBar onSearch={handleSearch} />
			<ul>
				{filtered.map(p => (
					<li key={p.id} className='p-2 border-b'>
						{p.content}
					</li>
				))}
			</ul>
		</div>
	)
}
