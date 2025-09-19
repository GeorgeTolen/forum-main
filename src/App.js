import React from 'react'
import { Routes, Route } from 'react-router-dom'
import BoardsPage from './pages/BoardsPage'
import PostsPage from './pages/PostsPage'
import HomePage from './pages/HomePage'

function App() {
	return (
		<Routes>
			<Route path='/' element={<HomePage />} />
			<Route path='/boards' element={<BoardsPage />} />
			<Route path='/posts' element={<PostsPage />} />
		</Routes>
	)
}

export default App
