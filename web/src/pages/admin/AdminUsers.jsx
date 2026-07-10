import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminAPI } from '../../api/admin.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const AdminUsers = () => {
  const navigate = useNavigate();
  const [users, setUsers] = useState([]);
  const [avatars, setAvatars] = useState({});
  const [loading, setLoading] = useState(true);
  const [loadingAvatars, setLoadingAvatars] = useState(false);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [filters, setFilters] = useState({
    page: 1,
    limit: 20,
    role: '',
    search: '',
  });

  useEffect(() => {
    loadUsers();
  }, [filters.page, filters.role, filters.search]);

  const loadUsers = async () => {
    setLoading(true);
    setError('');
    setAvatars({});
    try {
      const params = {};
      if (filters.page) params.page = filters.page;
      if (filters.limit) params.limit = filters.limit;
      if (filters.role) params.role = filters.role;
      if (filters.search) params.search = filters.search;

      const response = await adminAPI.listUsers(params);
      const usersData = response.data.users || [];
      setUsers(usersData);
      
      // Загружаем аватарки для всех пользователей
      if (usersData.length > 0) {
        loadAllAvatars(usersData);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки пользователей');
    } finally {
      setLoading(false);
    }
  };

  const loadAllAvatars = async (usersList) => {
    setLoadingAvatars(true);
    const avatarPromises = usersList.map(async (user) => {
      try {
        const response = await adminAPI.getUserAvatar(user.id);
        const url = URL.createObjectURL(response.data);
        return { id: user.id, url };
      } catch (err) {
        if (err.response?.status !== 404) {
          console.error(`Error loading avatar for user ${user.id}:`, err);
        }
        return { id: user.id, url: null };
      }
    });

    try {
      const results = await Promise.all(avatarPromises);
      const avatarMap = {};
      results.forEach(({ id, url }) => {
        if (url) {
          avatarMap[id] = url;
        }
      });
      setAvatars(avatarMap);
    } catch (err) {
      console.error('Error loading avatars:', err);
    } finally {
      setLoadingAvatars(false);
    }
  };

  const handleRoleChange = async (userId, newRole) => {
    try {
      await adminAPI.updateRole(userId, { role: newRole });
      setMessage('Роль успешно изменена');
      setTimeout(() => setMessage(''), 3000);
      loadUsers();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка изменения роли');
      setTimeout(() => setError(''), 3000);
    }
  };

  const handleRemoveManager = async (userId) => {
    if (!confirm('Удалить назначенного менеджера?')) return;
    try {
      await adminAPI.removeManager(userId);
      setMessage('Менеджер удален');
      setTimeout(() => setMessage(''), 3000);
      loadUsers();
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

  const getRoleColor = (role) => {
    const colors = {
      admin: 'bg-purple-100 text-purple-700',
      manager: 'bg-blue-100 text-blue-700',
      employee: 'bg-green-100 text-green-700',
    };
    return colors[role] || 'bg-gray-100 text-gray-700';
  };

  const getInitials = (firstName, lastName) => {
    return `${firstName?.[0] || ''}${lastName?.[0] || ''}`;
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Управление пользователями</h1>
          <button
            onClick={() => loadUsers()}
            className="btn btn-secondary"
          >
            Обновить
          </button>
        </div>

        {message && (
          <div className="success-message mb-4">
            {message}
          </div>
        )}
        
        {error && (
          <div className="error-message mb-4">
            {error}
          </div>
        )}

        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="label">Поиск</label>
              <input
                type="text"
                placeholder="Email или имя..."
                value={filters.search}
                onChange={(e) => setFilters({ ...filters, search: e.target.value, page: 1 })}
                className="input"
              />
            </div>
            <div>
              <label className="label">Роль</label>
              <select
                value={filters.role}
                onChange={(e) => setFilters({ ...filters, role: e.target.value, page: 1 })}
                className="input"
              >
                <option value="">Все роли</option>
                <option value="admin">Администратор</option>
                <option value="manager">Менеджер</option>
                <option value="employee">Сотрудник</option>
              </select>
            </div>
            <div className="flex items-end">
              <button
                onClick={() => setFilters({ page: 1, limit: 20, role: '', search: '' })}
                className="btn btn-secondary"
              >
                Сбросить
              </button>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow overflow-hidden">
          {loading ? (
            <div className="text-center py-8">Загрузка...</div>
          ) : users.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Пользователи не найдены</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Пользователь
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Email
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Роль
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Менеджер
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Действия
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {users.map((user) => (
                    <tr key={user.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center space-x-3">
                          {loadingAvatars ? (
                            <div className="avatar-placeholder-sm animate-pulse">
                              ...
                            </div>
                          ) : avatars[user.id] ? (
                            <img
                              src={avatars[user.id]}
                              alt="Avatar"
                              className="avatar-sm"
                            />
                          ) : (
                            <div className="avatar-placeholder-sm">
                              {getInitials(user.first_name, user.last_name)}
                            </div>
                          )}
                          <div className="text-sm font-medium text-gray-900">
                            {user.first_name} {user.last_name}
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-500">{user.email}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <select
                          value={user.role}
                          onChange={(e) => handleRoleChange(user.id, e.target.value)}
                          className={`px-2 py-1 text-sm font-medium rounded-full border-0 ${getRoleColor(user.role)}`}
                        >
                          <option value="admin" className="bg-white text-gray-900">Администратор</option>
                          <option value="manager" className="bg-white text-gray-900">Менеджер</option>
                          <option value="employee" className="bg-white text-gray-900">Сотрудник</option>
                        </select>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-500">
                          {user.manager_id ? (
                            <div className="flex items-center space-x-2">
                              <span>Назначен</span>
                              <button
                                onClick={() => handleRemoveManager(user.id)}
                                className="text-red-600 hover:text-red-800 text-xs"
                              >
                                Удалить
                              </button>
                            </div>
                          ) : (
                            <span className="text-gray-400">Не назначен</span>
                          )}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <button
                          onClick={() => navigate(`/admin/users/${user.id}`)}
                          className="text-blue-600 hover:text-blue-800"
                        >
                          Подробнее
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};