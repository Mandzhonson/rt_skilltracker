import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { adminAPI } from '../../api/admin.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const AdminUserDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [selectedManagerId, setSelectedManagerId] = useState('');
  
  const [managers, setManagers] = useState([]);
  const [managerSearch, setManagerSearch] = useState('');
  const [showManagerDropdown, setShowManagerDropdown] = useState(false);
  const [loadingManagers, setLoadingManagers] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    loadUser();
    loadManagers();
    
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowManagerDropdown(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [id]);

  const loadUser = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await adminAPI.getUser(id);
      setUser(response.data);
      if (response.data.manager_id) {
        setSelectedManagerId(response.data.manager_id);
      }
      loadAvatar();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки пользователя');
    } finally {
      setLoading(false);
    }
  };

  const loadAvatar = async () => {
    try {
      const response = await adminAPI.getUserAvatar(id);
      const url = URL.createObjectURL(response.data);
      setAvatarPreview(url);
    } catch (err) {
      if (err.response?.status !== 404) {
        console.error('Error loading avatar:', err);
      }
    }
  };

  const loadManagers = async () => {
    setLoadingManagers(true);
    try {
      const response = await adminAPI.listManagers();
      const allManagers = response.data.users || [];
      setManagers(allManagers.filter(m => m.id !== id));
    } catch (err) {
      console.error('Error loading managers:', err);
    } finally {
      setLoadingManagers(false);
    }
  };

  const filteredManagers = managers.filter(manager => {
    if (!managerSearch) return true;
    const search = managerSearch.toLowerCase();
    return (
      manager.first_name?.toLowerCase().includes(search) ||
      manager.last_name?.toLowerCase().includes(search) ||
      manager.email?.toLowerCase().includes(search)
    );
  });

  const handleAssignManager = async (managerId) => {
    if (!managerId) {
      setError('Выберите менеджера');
      return;
    }
    try {
      await adminAPI.assignManager(id, { manager_id: managerId });
      setMessage('Менеджер успешно назначен');
      setSelectedManagerId(managerId);
      setShowManagerDropdown(false);
      setManagerSearch('');
      setTimeout(() => setMessage(''), 3000);
      loadUser();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка назначения менеджера');
      setTimeout(() => setError(''), 3000);
    }
  };

  const handleRemoveManager = async () => {
    if (!confirm('Удалить назначенного менеджера?')) return;
    try {
      await adminAPI.removeManager(id);
      setMessage('Менеджер удален');
      setSelectedManagerId('');
      setManagerSearch('');
      setTimeout(() => setMessage(''), 3000);
      loadUser();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка удаления менеджера');
      setTimeout(() => setError(''), 3000);
    }
  };

  const getRoleLabel = (role) => {
    const roles = {
      admin: 'Администратор',
      manager: 'Менеджер',
      employee: 'Сотрудник',
    };
    return roles[role] || role;
  };

  const getInitials = (firstName, lastName) => {
    return `${firstName?.[0] || ''}${lastName?.[0] || ''}`;
  };

  const getManagerName = (managerId) => {
    const manager = managers.find(m => m.id === managerId);
    if (manager) {
      return `${manager.first_name} ${manager.last_name}`;
    }
    return managerId;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="flex items-center justify-center h-64">Загрузка...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate('/admin/users')} className="btn btn-primary mt-4">
            Вернуться к списку
          </button>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="text-center">Пользователь не найден</div>
          <button onClick={() => navigate('/admin/users')} className="btn btn-primary mt-4">
            Вернуться к списку
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="max-w-3xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/admin/users')}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад к списку
        </button>

        <div className="card">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Детали пользователя</h1>

          {message && (
            <div className="success-message mb-4">{message}</div>
          )}
          {error && (
            <div className="error-message mb-4">{error}</div>
          )}

          <div className="flex items-center space-x-4 mb-6 pb-6 border-b">
            <div>
              {avatarPreview ? (
                <img
                  src={avatarPreview}
                  alt="Avatar"
                  className="avatar"
                />
              ) : (
                <div className="avatar-placeholder">
                  {getInitials(user.first_name, user.last_name)}
                </div>
              )}
            </div>
            <div>
              <p className="text-xl font-semibold text-gray-900">
                {user.first_name} {user.last_name}
              </p>
              <p className="text-gray-500">{user.email}</p>
            </div>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-500">ID</label>
              <p className="text-gray-900">{user.id}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Email</label>
              <p className="text-gray-900">{user.email}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Имя</label>
              <p className="text-gray-900">{user.first_name} {user.last_name}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Роль</label>
              <span className="px-3 py-1 inline-block rounded-full text-sm font-medium">
                {getRoleLabel(user.role)}
              </span>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Менеджер</label>
              {selectedManagerId ? (
                <div className="flex items-center space-x-2">
                  <span className="text-gray-900">{getManagerName(selectedManagerId)}</span>
                  <button
                    onClick={handleRemoveManager}
                    className="text-red-600 hover:text-red-800 text-sm"
                  >
                    Удалить
                  </button>
                </div>
              ) : (
                <p className="text-gray-400">Не назначен</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Дата создания</label>
              <p className="text-gray-900">{new Date(user.created_at).toLocaleString()}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Дата обновления</label>
              <p className="text-gray-900">{new Date(user.updated_at).toLocaleString()}</p>
            </div>

            <div className="pt-4 border-t">
              <h3 className="font-semibold text-gray-900 mb-3">Назначить менеджера</h3>
              <div ref={dropdownRef} className="relative">
                <input
                  type="text"
                  placeholder="Поиск менеджера по имени или email..."
                  value={managerSearch}
                  onChange={(e) => {
                    setManagerSearch(e.target.value);
                    setShowManagerDropdown(true);
                  }}
                  onFocus={() => setShowManagerDropdown(true)}
                  className="input"
                />
                
                {showManagerDropdown && (
                  <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {loadingManagers ? (
                      <div className="px-4 py-2 text-gray-500">Загрузка...</div>
                    ) : filteredManagers.length === 0 ? (
                      <div className="px-4 py-2 text-gray-500">Менеджеры не найдены</div>
                    ) : (
                      filteredManagers.map((manager) => (
                        <div
                          key={manager.id}
                          className="px-4 py-2 hover:bg-gray-50 cursor-pointer flex items-center space-x-3"
                          onClick={() => {
                            handleAssignManager(manager.id);
                          }}
                        >
                          <div className="avatar-placeholder-sm">
                            {getInitials(manager.first_name, manager.last_name)}
                          </div>
                          <div>
                            <div className="text-sm font-medium text-gray-900">
                              {manager.first_name} {manager.last_name}
                            </div>
                            <div className="text-xs text-gray-500">{manager.email}</div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};