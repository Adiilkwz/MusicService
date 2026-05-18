import { useState } from 'react';

export default function Register({ onNavigate }) {
  const [formData, setFormData] = useState({ email: '', password: '', display_name: '' });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    
    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const data = await response.json();
      if (!response.ok) throw new Error(data.error || 'Registration failed');

      setSuccess('Account created! You can now log in.');
      setTimeout(() => onNavigate('login'), 2000);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '100px auto', fontFamily: 'sans-serif' }}>
      <h2>Create an Account</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {success && <p style={{ color: '#1DB954' }}>{success}</p>}
      
      <form onSubmit={handleRegister} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <input 
          type="text" placeholder="Display Name" required style={{ padding: '10px' }}
          onChange={(e) => setFormData({...formData, display_name: e.target.value})} 
        />
        <input 
          type="email" placeholder="Email" required style={{ padding: '10px' }}
          onChange={(e) => setFormData({...formData, email: e.target.value})} 
        />
        <input 
          type="password" placeholder="Password (min 6 chars)" required style={{ padding: '10px' }}
          onChange={(e) => setFormData({...formData, password: e.target.value})} 
        />
        <button type="submit" style={{ padding: '10px', background: '#1DB954', color: 'white', border: 'none', cursor: 'pointer' }}>
          Register
        </button>
      </form>
      <button onClick={() => onNavigate('login')} style={{ marginTop: '15px', background: 'none', border: 'none', color: 'blue', cursor: 'pointer' }}>
        Already have an account? Log in
      </button>
    </div>
  );
}