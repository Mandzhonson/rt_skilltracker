import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext.jsx';

export const Navbar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const getRoleLabel = (role) => {
    const roles = {
      admin: 'Администратор',
      manager: 'Менеджер',
      employee: 'Сотрудник',
    };
    return roles[role] || role;
  };

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <nav style={{ backgroundColor: 'white', boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1)' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 1rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', height: '64px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '2rem' }}>
            <Link to="/profile" style={{ fontSize: '1.25rem', fontWeight: 'bold', color: '#2563eb' }}>
              SkillTracker
            </Link>
            {user && (
              <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
                <Link to="/profile" style={{ color: '#374151' }}>
                  Профиль
                </Link>
                {user.role === 'admin' && (
                  <Link to="/admin/users" style={{ color: '#374151' }}>
                    Пользователи
                  </Link>
                )}
                {user.role === 'manager' && (
                  <Link to="/manager" style={{ color: '#374151' }}>
                    Управление
                  </Link>
                )}
                {user.role === 'employee' && (
                  <Link to="/employee" style={{ color: '#374151' }}>
                    Мои планы
                  </Link>
                )}
              </div>
            )}
          </div>

          {user && (
            <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
              <span style={{ fontSize: '0.875rem', color: '#4b5563' }}>
                {getRoleLabel(user.role)}
              </span>
              <span style={{ fontSize: '0.875rem', color: '#374151' }}>
                {user.first_name} {user.last_name}
              </span>
              <button
                onClick={handleLogout}
                style={{
                  padding: '0.5rem 1rem',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'white',
                  backgroundColor: '#dc2626',
                  borderRadius: '0.375rem',
                  border: 'none',
                  cursor: 'pointer',
                }}
                onMouseEnter={(e) => {
                  e.target.style.backgroundColor = '#b91c1c';
                }}
                onMouseLeave={(e) => {
                  e.target.style.backgroundColor = '#dc2626';
                }}
              >
                Выйти
              </button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
};