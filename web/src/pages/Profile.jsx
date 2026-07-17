import { useState, useRef, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext.jsx';
import { userAPI } from '../api/user.js';
import { Navbar } from '../components/Layout/Navbar.jsx';

export const Profile = () => {
  const { user, setUser } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const fileInputRef = useRef(null);
  
  const [profileData, setProfileData] = useState({
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    email: user?.email || '',
  });
  
  const [passwordData, setPasswordData] = useState({
    old_password: '',
    new_password: '',
    confirm_password: '',
  });
  
  const [message, setMessage] = useState({ text: '', type: '' });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (user) {
      setProfileData({
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        email: user.email || '',
      });
      loadAvatar();
    }
  }, [user]);

  const loadAvatar = async () => {
    try {
      const response = await userAPI.getAvatar();
      const url = URL.createObjectURL(response.data);
      setAvatarPreview(url);
    } catch (error) {
      if (error.response?.status !== 404) {
        console.error('Error loading avatar:', error);
      }
    }
  };

  const getRoleDisplay = (role) => {
    const roles = {
      admin: 'Администратор',
      manager: 'Менеджер',
      employee: 'Сотрудник',
    };
    return roles[role] || role;
  };

  const handleProfileUpdate = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage({ text: '', type: '' });
    
    try {
      const response = await userAPI.updateProfile({
        email: profileData.email,
        first_name: profileData.first_name,
        last_name: profileData.last_name,
      });
      setUser(response.data);
      setIsEditing(false);
      setMessage({ text: 'Профиль успешно обновлен', type: 'success' });
      setTimeout(() => setMessage({ text: '', type: '' }), 3000);
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.error || 'Ошибка обновления профиля', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordChange = async (e) => {
    e.preventDefault();
    
    if (passwordData.new_password !== passwordData.confirm_password) {
      setMessage({ text: 'Пароли не совпадают', type: 'error' });
      return;
    }
    
    if (passwordData.new_password.length < 6) {
      setMessage({ text: 'Пароль должен быть минимум 6 символов', type: 'error' });
      return;
    }
    
    setLoading(true);
    setMessage({ text: '', type: '' });
    
    try {
      await userAPI.updatePassword({
        old_password: passwordData.old_password,
        new_password: passwordData.new_password,
      });
      setMessage({ text: 'Пароль успешно изменен', type: 'success' });
      setIsChangingPassword(false);
      setPasswordData({
        old_password: '',
        new_password: '',
        confirm_password: '',
      });
      setTimeout(() => setMessage({ text: '', type: '' }), 3000);
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.error || 'Ошибка смены пароля', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    
    const maxSize = 5 * 1024 * 1024;
    if (file.size > maxSize) {
      setMessage({ text: 'Файл слишком большой (макс. 5MB)', type: 'error' });
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      return;
    }
    
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
    if (!allowedTypes.includes(file.type)) {
      setMessage({ text: 'Поддерживаются только JPEG, PNG, GIF', type: 'error' });
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      return;
    }
    
    setLoading(true);
    setMessage({ text: '', type: '' });
    
    const formData = new FormData();
    formData.append('avatar', file);
    
    try {
      await userAPI.setAvatar(formData);
      loadAvatar();
      setMessage({ text: 'Аватар успешно загружен', type: 'success' });
      setTimeout(() => setMessage({ text: '', type: '' }), 3000);
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.error || 'Ошибка загрузки аватара', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleAvatarDelete = async () => {
    if (!confirm('Удалить аватар?')) return;
    
    setLoading(true);
    setMessage({ text: '', type: '' });
    
    try {
      await userAPI.deleteAvatar();
      setAvatarPreview(null);
      setMessage({ text: 'Аватар удален', type: 'success' });
      setTimeout(() => setMessage({ text: '', type: '' }), 3000);
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.error || 'Ошибка удаления аватара', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f9fafb' }}>
      <Navbar />
      
      <div style={{ maxWidth: '48rem', margin: '0 auto', padding: '2rem 1rem' }}>
        <div className="card">
          <h1 style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#111827', marginBottom: '1.5rem' }}>
            Профиль пользователя
          </h1>
          
          {message.text && (
            <div style={{
              marginBottom: '1rem',
              padding: '1rem',
              borderRadius: '0.375rem',
              backgroundColor: message.type === 'success' ? '#f0fdf4' : '#fef2f2',
              border: message.type === 'success' ? '1px solid #86efac' : '1px solid #fca5a5',
              color: message.type === 'success' ? '#166534' : '#991b1b'
            }}>
              {message.text}
            </div>
          )}

          {/* Аватар */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', marginBottom: '1.5rem', paddingBottom: '1.5rem', borderBottom: '1px solid #e5e7eb' }}>
            <div>
              {avatarPreview ? (
                <img
                  src={avatarPreview}
                  alt="Avatar"
                  style={{ width: '96px', height: '96px', borderRadius: '50%', objectFit: 'cover', border: '2px solid #e5e7eb' }}
                />
              ) : (
                <div style={{ width: '96px', height: '96px', borderRadius: '50%', backgroundColor: '#e5e7eb', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '1.5rem', color: '#6b7280' }}>
                  {user?.first_name?.[0]}{user?.last_name?.[0]}
                </div>
              )}
            </div>
            <div>
              <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                <button
                  type="button"
                  onClick={() => fileInputRef.current?.click()}
                  className="btn btn-primary"
                  disabled={loading}
                >
                  Добавить аватар
                </button>
                {avatarPreview && (
                  <button
                    type="button"
                    onClick={handleAvatarDelete}
                    className="btn btn-danger"
                    disabled={loading}
                  >
                    Удалить аватар
                  </button>
                )}
              </div>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                onChange={handleAvatarUpload}
                style={{ display: 'none' }}
              />
              <p style={{ fontSize: '0.75rem', color: '#6b7280', marginTop: '0.5rem' }}>
                JPEG, PNG, GIF до 5 MB
              </p>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div>
              <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' }}>Email</label>
              <p style={{ fontSize: '1.125rem', fontWeight: '500', color: '#111827' }}>{user?.email}</p>
            </div>

            <div>
              <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' }}>Роль</label>
              <p style={{ fontSize: '1.125rem', fontWeight: '600', color: '#2563eb' }}>{getRoleDisplay(user?.role)}</p>
            </div>

            <div>
              <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' }}>Должность</label>
              <p style={{ fontSize: '1.125rem', color: '#111827' }}>{user?.position || 'Не указана'}</p>
            </div>

            {isEditing ? (
              <form onSubmit={handleProfileUpdate} style={{ display: 'flex', flexDirection: 'column', gap: '1rem', paddingTop: '1rem', borderTop: '1px solid #e5e7eb' }}>
                <div>
                  <label className="label">Email</label>
                  <input
                    type="email"
                    required
                    value={profileData.email}
                    onChange={(e) => setProfileData({ ...profileData, email: e.target.value })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">Имя</label>
                  <input
                    type="text"
                    required
                    value={profileData.first_name}
                    onChange={(e) => setProfileData({ ...profileData, first_name: e.target.value })}
                    className="input"
                  />
                </div>
                
                <div>
                  <label className="label">Фамилия</label>
                  <input
                    type="text"
                    required
                    value={profileData.last_name}
                    onChange={(e) => setProfileData({ ...profileData, last_name: e.target.value })}
                    className="input"
                  />
                </div>
                
                <div style={{ display: 'flex', gap: '0.75rem', paddingTop: '0.5rem' }}>
                  <button
                    type="submit"
                    disabled={loading}
                    className="btn btn-primary"
                  >
                    {loading ? 'Сохранение...' : 'Сохранить'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setIsEditing(false);
                      setProfileData({
                        first_name: user?.first_name || '',
                        last_name: user?.last_name || '',
                        email: user?.email || '',
                      });
                    }}
                    className="btn btn-secondary"
                  >
                    Отмена
                  </button>
                </div>
              </form>
            ) : (
              <>
                <div style={{ paddingTop: '1rem', borderTop: '1px solid #e5e7eb' }}>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                    <div>
                      <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' }}>Имя</label>
                      <p style={{ fontSize: '1.125rem', color: '#111827' }}>{user?.first_name}</p>
                    </div>
                    <div>
                      <label style={{ fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' }}>Фамилия</label>
                      <p style={{ fontSize: '1.125rem', color: '#111827' }}>{user?.last_name}</p>
                    </div>
                  </div>
                </div>
                
                <div style={{ display: 'flex', gap: '0.75rem', marginTop: '0.5rem', flexWrap: 'wrap' }}>
                  <button
                    onClick={() => setIsEditing(true)}
                    className="btn btn-primary"
                  >
                    Редактировать профиль
                  </button>
                  <button
                    onClick={() => setIsChangingPassword(!isChangingPassword)}
                    className="btn btn-secondary"
                  >
                    Сменить пароль
                  </button>
                </div>
              </>
            )}

            {isChangingPassword && (
              <form onSubmit={handlePasswordChange} style={{ display: 'flex', flexDirection: 'column', gap: '1rem', paddingTop: '1rem', borderTop: '1px solid #e5e7eb' }}>
                <h3 style={{ fontWeight: '600', color: '#111827' }}>Смена пароля</h3>
                <div>
                  <label className="label">Текущий пароль</label>
                  <input
                    type="password"
                    required
                    value={passwordData.old_password}
                    onChange={(e) => setPasswordData({ ...passwordData, old_password: e.target.value })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">Новый пароль</label>
                  <input
                    type="password"
                    required
                    minLength={6}
                    value={passwordData.new_password}
                    onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">Подтверждение пароля</label>
                  <input
                    type="password"
                    required
                    value={passwordData.confirm_password}
                    onChange={(e) => setPasswordData({ ...passwordData, confirm_password: e.target.value })}
                    className="input"
                  />
                </div>
                <div style={{ display: 'flex', gap: '0.75rem' }}>
                  <button
                    type="submit"
                    disabled={loading}
                    className="btn btn-primary"
                  >
                    {loading ? 'Смена...' : 'Сменить пароль'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setIsChangingPassword(false);
                      setPasswordData({
                        old_password: '',
                        new_password: '',
                        confirm_password: '',
                      });
                    }}
                    className="btn btn-secondary"
                  >
                    Отмена
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};