import { useState, useEffect } from 'react';
import Navbar from './components/Navbar';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import Profile from './pages/auth/Profile';
import Admin from './pages/auth/Admin';
import Dashboard from './pages/stream/Dashboard';

function App() {
  const [token, setToken] = useState(() => localStorage.getItem('jwt_token') || null);
  
  const [currentPage, setCurrentPage] = useState('dashboard');

  useEffect(() => {
    if (token) {
      localStorage.setItem('jwt_token', token);
    } else {
      localStorage.removeItem('jwt_token');
    }
  }, [token]);

  const handleLoginSuccess = (jwt) => {
    setToken(jwt);
    setCurrentPage('dashboard');
  };

  const handleLogout = () => {
    setToken(null);
    setCurrentPage('dashboard');
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'login':
        return <Login onLoginSuccess={handleLoginSuccess} onNavigate={setCurrentPage} />;
      case 'register':
        return <Register onNavigate={setCurrentPage} />;
      case 'profile':
        return token ? <Profile token={token} onLogout={handleLogout} /> : <Login onLoginSuccess={handleLoginSuccess} onNavigate={setCurrentPage} />;
      case 'admin':
        return token ? <Admin token={token} /> : <Login onLoginSuccess={handleLoginSuccess} onNavigate={setCurrentPage} />;
      case 'dashboard':
      default:
        return <Dashboard token={token} />;
    }
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', background: '#121212', color: 'white' }}>
      <Navbar 
        token={token} 
        currentPage={currentPage} 
        setCurrentPage={setCurrentPage} 
        onLogout={handleLogout} 
      />

      <main style={{ flex: 1, overflowY: 'auto' }}>
        {renderPage()}
      </main>
    </div>
  );
}

export default App;