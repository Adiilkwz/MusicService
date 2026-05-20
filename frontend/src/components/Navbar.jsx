export default function Navbar({ token, currentPage, setCurrentPage, onLogout }) {
  const navStyle = (pageName) => ({
    background: 'none',
    border: 'none',
    cursor: 'pointer',
    fontSize: '16px',
    fontWeight: currentPage === pageName ? 'bold' : 'normal',
    color: currentPage === pageName ? '#1DB954' : 'white',
  });

  return (
    <nav style={{ background: '#121212', padding: '15px 30px', display: 'flex', alignItems: 'center', boxShadow: '0 2px 10px rgba(0,0,0,0.5)' }}>
      
      <div 
        onClick={() => setCurrentPage('dashboard')}
        style={{ color: 'white', fontWeight: 'bold', fontSize: '22px', cursor: 'pointer', marginRight: '40px' }}
      >
        GoMusic
      </div>
      
      <div style={{ display: 'flex', gap: '20px' }}>
        <button onClick={() => setCurrentPage('dashboard')} style={navStyle('dashboard')}>Home / Catalog</button>
      </div>

      <div style={{ marginLeft: 'auto', display: 'flex', gap: '15px', alignItems: 'center' }}>
        
        {token ? (
          <>
            <button onClick={() => setCurrentPage('profile')} style={navStyle('profile')}>Profile</button>
            <button onClick={() => setCurrentPage('admin')} style={navStyle('admin')}>Admin</button>
            <button 
              onClick={onLogout} 
              style={{ padding: '8px 16px', background: '#ff4d4f', color: 'white', border: 'none', borderRadius: '20px', cursor: 'pointer', fontWeight: 'bold' }}
            >
              Logout
            </button>
          </>
        ) : (
          <>
            <button 
              onClick={() => setCurrentPage('register')} 
              style={{ background: 'none', border: 'none', color: '#aaa', cursor: 'pointer', fontWeight: 'bold' }}
            >
              Sign Up
            </button>
            <button 
              onClick={() => setCurrentPage('login')} 
              style={{ padding: '8px 24px', background: 'white', color: 'black', border: 'none', borderRadius: '20px', cursor: 'pointer', fontWeight: 'bold' }}
            >
              Log In
            </button>
          </>
        )}
      </div>
    </nav>
  );
}