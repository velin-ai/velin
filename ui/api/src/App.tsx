import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'

function App() {
  const path = window.location.pathname.replace(/\/+$/, '') || '/'

  if (path === '/login') {
    return <Login />
  }

  if (path === '/register') {
    return <Register />
  }

  if (path === '/dashboard') {
    return <Dashboard />
  }

  return <Home />
}

export default App
