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

  // Состояния для должности
  const [positionInput, setPositionInput] = useState('');
  const [isEditingPosition, setIsEditingPosition] = useState(false);
  const [positionLoading, setPositionLoading] = useState(false);

  // Состояния для навыков и планов
  const [skills, setSkills] = useState([]);
  const [plans, setPlans] = useState([]);
  const [loadingProfile, setLoadingProfile] = useState(false);

  useEffect(() => {
    loadUser();
    loadManagers();
    loadUserProfile();
    
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
      setPositionInput(response.data.position || '');
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

  const loadUserProfile = async () => {
    setLoadingProfile(true);
    try {
      const response = await adminAPI.getUserProfile(id);
      const data = response.data;
      setSkills(data.skills || []);
      setPlans(data.plans || []);
    } catch (err) {
      console.error('Error loading user profile:', err);
    } finally {
      setLoadingProfile(false);
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

  const handleUpdatePosition = async (e) => {
    e.preventDefault();
    setPositionLoading(true);
    setError('');
    try {
      await adminAPI.updatePosition(id, { position: positionInput });
      setMessage('Должность успешно обновлена');
      setIsEditingPosition(false);
      setTimeout(() => setMessage(''), 3000);
      loadUser();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка обновления должности');
      setTimeout(() => setError(''), 3000);
    } finally {
      setPositionLoading(false);
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

  const getStatusLabel = (status) => {
    const statuses = {
      draft: 'Черновик',
      active: 'Активный',
      completed: 'Завершен',
      archived: 'Архивный',
    };
    return statuses[status] || status;
  };

  const getStatusColor = (status) => {
    const colors = {
      draft: 'bg-gray-100 text-gray-600',
      active: 'bg-blue-100 text-blue-700',
      completed: 'bg-green-100 text-green-700',
      archived: 'bg-yellow-100 text-yellow-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-600';
  };

  const getCategoryLabel = (category) => {
    const labels = {
      frontend: 'Фронтенд',
      backend: 'Бэкенд',
      devops: 'DevOps',
      database: 'Базы данных',
      testing: 'Тестирование',
      cloud: 'Облачные технологии',
      mobile: 'Мобильная разработка',
      architecture: 'Архитектура',
      ai: 'ИИ и ML',
      security: 'Безопасность',
      soft_skills: 'Soft Skills',
      other: 'Другое',
    };
    return labels[category] || category;
  };

  const getCategoryColor = (category) => {
    const colors = {
      frontend: 'bg-blue-100 text-blue-700',
      backend: 'bg-emerald-100 text-emerald-700',
      devops: 'bg-purple-100 text-purple-700',
      database: 'bg-amber-100 text-amber-700',
      testing: 'bg-rose-100 text-rose-700',
      cloud: 'bg-indigo-100 text-indigo-700',
      mobile: 'bg-pink-100 text-pink-700',
      architecture: 'bg-orange-100 text-orange-700',
      ai: 'bg-cyan-100 text-cyan-700',
      security: 'bg-red-100 text-red-700',
      soft_skills: 'bg-teal-100 text-teal-700',
      other: 'bg-gray-100 text-gray-700',
    };
    return colors[category] || colors.other;
  };

  const formatDate = (dateString) => {
    if (!dateString) return '—';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return '—';
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch {
      return '—';
    }
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

        <div className="card mb-6">
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
              <div className="flex items-center justify-between">
                <label className="block text-sm font-medium text-gray-500">Должность</label>
                <button
                  onClick={() => {
                    setIsEditingPosition(!isEditingPosition);
                    setPositionInput(user.position || '');
                  }}
                  className="text-blue-600 hover:text-blue-800 text-sm"
                >
                  {isEditingPosition ? 'Отмена' : 'Изменить'}
                </button>
              </div>
              
              {isEditingPosition ? (
                <form onSubmit={handleUpdatePosition} className="mt-2">
                  <div className="flex gap-2">
                    <input
                      type="text"
                      value={positionInput}
                      onChange={(e) => setPositionInput(e.target.value)}
                      placeholder="Введите должность"
                      className="input flex-1"
                      disabled={positionLoading}
                    />
                    <button
                      type="submit"
                      className="btn btn-primary"
                      disabled={positionLoading}
                    >
                      {positionLoading ? 'Сохранение...' : 'Сохранить'}
                    </button>
                  </div>
                </form>
              ) : (
                <p className="text-gray-900">{user.position || 'Не указана'}</p>
              )}
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
              <p className="text-gray-900">{formatDate(user.created_at)}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500">Дата обновления</label>
              <p className="text-gray-900">{formatDate(user.updated_at)}</p>
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

        {/* Навыки — переделаны на список в один столбец со скроллом */}
        <div className="mb-6">
          <div className="flex justify-between items-center mb-3">
            <h2 className="text-xl font-semibold text-gray-900">Навыки</h2>
            {skills.length > 0 && (
              <span className="text-sm text-gray-400">{skills.length} навыков</span>
            )}
          </div>
          
          {skills.length === 0 ? (
            <div className="card text-center py-4 border-card">
              <p className="text-gray-500">Нет навыков</p>
            </div>
          ) : (
            <div className="scroll-container">
              <div className="skills-column-container">
                {skills.map((skill) => (
                  <div
                    key={skill.id}
                    className="skill-list-item hover:bg-gray-50 transition-colors text-left"
                  >
                    <div className="flex justify-between items-center gap-4">
                      <div>
                        <p className="font-medium text-gray-900 text-sm">{skill.name}</p>
                        {skill.description && (
                          <p className="text-xs text-gray-500 mt-1">{skill.description}</p>
                        )}
                      </div>
                      
                      <div className="flex items-center gap-2 flex-wrap shrink-0">
                        {skill.category && (
                          <span className={`text-xs px-2 py-0.5 rounded font-medium ${getCategoryColor(skill.category)}`}>
                            {getCategoryLabel(skill.category)}
                          </span>
                        )}
                        <span className="text-xs text-gray-400">
                          {formatDate(skill.created_at)}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Планы */}
        {plans.length > 0 && (
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Планы</h2>
            <div className="space-y-3">
              {plans.map((plan) => (
                <div
                  key={plan.id}
                  className="border rounded-lg p-4 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="font-medium text-gray-900">{plan.title}</h3>
                      {plan.description && (
                        <p className="text-sm text-gray-500 mt-1">{plan.description}</p>
                      )}
                      <div className="flex flex-wrap items-center gap-3 mt-2 text-sm">
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(plan.status)}`}>
                          {getStatusLabel(plan.status)}
                        </span>
                        <span className="text-gray-500">Прогресс: {plan.progress || 0}%</span>
                        <span className="text-gray-400 text-xs">
                          {formatDate(plan.created_at)}
                        </span>
                      </div>
                    </div>
                    <button
                      onClick={() => navigate(`/admin/plans/${plan.id}`)}
                      className="btn btn-primary text-sm"
                    >
                      Открыть
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};