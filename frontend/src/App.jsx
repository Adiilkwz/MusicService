import { useState } from 'react';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import Profile from './pages/auth/Profile';
import Admin from './pages/auth/Admin';
import Dashboard from './pages/stream/Dashboard';

function App() {
  const [token, setToken] = useState(null);
  const [currentPage, setCurrentPage] = useState('login');

  const handleLogout = () => {
    setToken(null);
    setCurrentPage('login');
  };

  if (!token) {
    if (currentPage === 'register') return <Register onNavigate={setCurrentPage} />;
    return <Login onLoginSuccess={(jwt) => { setToken(jwt); setCurrentPage('dashboard'); }} onNavigate={setCurrentPage} />;
  }

  return (
    <div>
      <nav style={{ background: '#121212', padding: '15px 30px', display: 'flex', gap: '20px', color: 'white' }}>
        <button onClick={() => setCurrentPage('dashboard')} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Dashboard</button>
        <button onClick={() => setCurrentPage('profile')} style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer' }}>Profile</button>
        <button onClick={() => setCurrentPage('admin')} style={{ background: 'none', border: 'none', color: '#f39c12', cursor: 'pointer' }}>Admin Panel</button>
        <button onClick={handleLogout} style={{ background: '#ff4d4f', border: 'none', color: 'white', padding: '5px 10px', marginLeft: 'auto', cursor: 'pointer' }}>Logout</button>
      </nav>

      <main>
        {currentPage === 'dashboard' && <Dashboard token={token} />}
        {currentPage === 'profile' && <Profile token={token} onLogout={handleLogout} />}
        {currentPage === 'admin' && <Admin token={token} />}
      </main>
    </div>
  );
}

export default App;