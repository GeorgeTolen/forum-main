import React, { useState } from 'react'
import SearchBar from '../components/SearchBar'

export default function BoardsPage() {
	const [boards] = useState([
		{ id: 1, title: 'React проекты' },
		{ id: 2, title: 'Solidity смарт-контракты' },
		{ id: 3, title: 'Кибербезопасность' },
	])
	const [filtered, setFiltered] = useState(boards)

	const handleSearch = query => {
		setFiltered(
			boards.filter(b => b.title.toLowerCase().includes(query.toLowerCase()))
		)
	}

	return (
		<div>
			<h1 className='text-xl font-bold mb-4'>Доски</h1>
			<SearchBar onSearch={handleSearch} />
			<ul>
				{filtered.map(b => (
					<li key={b.id} className='p-2 border-b'>
						{b.title}
					</li>
				))}
			</ul>
		</div>
	)
}
