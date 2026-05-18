import { useState, useEffect } from 'react';

export default function Profile({ token, onLogout }) {
  const [profile, setProfile] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState('');
  const [message, setMessage] = useState('');

  const fetchProfile = async () => {
    const res = await fetch('/api/profile/', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await res.json();
    setProfile(data);
    setEditName(data.display_name || data.DisplayName || '');
  };

  useEffect(() => { fetchProfile(); }, []);

  const handleUpdate = async () => {
    const res = await fetch('/api/v1/profile/', {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ display_name: editName, avatar_url: '' })
    });
    if (res.ok) {
      setMessage('Profile updated!');
      setIsEditing(false);
      fetchProfile();
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Are you sure? This cannot be undone.")) return;
    const res = await fetch('/api/v1/profile/', {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.ok) onLogout();
  };

  if (!profile) return <div>Loading...</div>;

  return (
    <div style={{ padding: '20px', maxWidth: '600px', margin: '0 auto' }}>
      <h2>👤 My Profile</h2>
      {message && <p style={{ color: '#1DB954' }}>{message}</p>}
      
      <div style={{ background: '#282c34', color: 'white', padding: '25px', borderRadius: '12px' }}>
        {isEditing ? (
          <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
            <input 
              value={editName} 
              onChange={(e) => setEditName(e.target.value)} 
              style={{ padding: '8px' }}
            />
            <button onClick={handleUpdate} style={{ background: '#1DB954', border: 'none', padding: '8px 15px', color: 'white', cursor: 'pointer' }}>Save</button>
            <button onClick={() => setIsEditing(false)} style={{ background: 'grey', border: 'none', padding: '8px 15px', color: 'white', cursor: 'pointer' }}>Cancel</button>
          </div>
        ) : (
          <div style={{ marginBottom: '15px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold' }}>{profile.display_name || profile.DisplayName}</div>
            <button onClick={() => setIsEditing(true)} style={{ background: 'none', border: 'none', color: '#1DB954', cursor: 'pointer', padding: '0' }}>Edit Name</button>
          </div>
        )}
        
        <p><strong>Email:</strong> {profile.email || profile.Email}</p>
        <p><strong>Role:</strong> {profile.role || profile.Role}</p>
        <p><strong>ID:</strong> {profile.user_id || profile.UserId}</p>
        
        <button onClick={handleDelete} style={{ marginTop: '20px', background: '#ff4d4f', color: 'white', border: 'none', padding: '10px 15px', cursor: 'pointer', borderRadius: '4px' }}>
          Delete Account
        </button>
      </div>
    </div>
  );
}