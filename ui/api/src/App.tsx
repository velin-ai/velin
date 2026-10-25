import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'

function App() {
  const path = window.location.pathname.replace(/\/+$/, '') || '/'

  if (path === '/login') {
    return <Login />
  }

  if (path === '/register') {
    return <Register />
  }

  return <Home />
}

export default App
