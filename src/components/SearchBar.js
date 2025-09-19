import React, { useState } from 'react'

export default function SearchBar({ onSearch }) {
	const [query, setQuery] = useState('')

	const handleChange = e => {
		setQuery(e.target.value)
		onSearch(e.target.value)
	}

	return (
		<input
			type='text'
			placeholder='Поиск...'
			value={query}
			onChange={handleChange}
			className='border p-2 rounded w-full'
		/>
	)
}
