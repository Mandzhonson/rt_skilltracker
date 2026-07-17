import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const ManagerEmployeeDetail = () => {
  const { employeeId } = useParams();
  const navigate = useNavigate();
  const [employee, setEmployee] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [skills, setSkills] = useState([]);
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadEmployeeProfile();
  }, [employeeId]);

  const loadEmployeeProfile = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await managerAPI.getEmployeeProfile(employeeId);
      const data = response.data;
      setEmployee(data.user);
      setSkills(data.skills || []);
      setPlans(data.plans || []);
      
      loadAvatar();
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки профиля сотрудника');
    } finally {
      setLoading(false);
    }
  };

  const loadAvatar = async () => {
    try {
      const response = await managerAPI.getEmployeeAvatar(employeeId);
      const url = URL.createObjectURL(response.data);
      setAvatarPreview(url);
    } catch (err) {
      if (err.response?.status !== 404) {
        console.error('Error loading avatar:', err);
      }
    }
  };

  const getInitials = (firstName, lastName) => {
    return `${firstName?.[0] || ''}${lastName?.[0] || ''}`;
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
        <div className="flex items-center justify-center h-64">
          <div className="text-gray-500">Загрузка...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate('/manager/employees')} className="btn btn-primary mt-4">
            Вернуться к сотрудникам
          </button>
        </div>
      </div>
    );
  }

  if (!employee) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="text-center">Сотрудник не найден</div>
          <button onClick={() => navigate('/manager/employees')} className="btn btn-primary mt-4">
            Вернуться к сотрудникам
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      {/* Здесь выставлен max-w-3xl (как у админа) вместо max-w-6xl, чтобы страница не была растянутой */}
      <div className="max-w-3xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/manager/employees')}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад к сотрудникам
        </button>

        {/* Карточка профиля с точной структурой и размерами из AdminUserDetail */}
        <div className="card mb-6">
          <div className="flex items-center space-x-4">
            <div>
              {avatarPreview ? (
                <img
                  src={avatarPreview}
                  alt="Avatar"
                  className="avatar"
                />
              ) : (
                <div className="avatar-placeholder">
                  {getInitials(employee.first_name, employee.last_name)}
                </div>
              )}
            </div>
            <div>
              <p className="text-xl font-semibold text-gray-900">
                {employee.first_name} {employee.last_name}
              </p>
              <p className="text-gray-500">{employee.email}</p>
              <p className="text-sm text-gray-400 mt-1">{employee.position || 'Сотрудник'}</p>
            </div>
          </div>
        </div>

        {/* Навыки в один вертикальный столбец с прокруткой */}
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
        <div>
          <h2 className="text-xl font-semibold text-gray-900 mb-3">Планы</h2>
          <div className="card">
            {plans.length === 0 ? (
              <p className="text-gray-500 text-center py-4">Нет планов</p>
            ) : (
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
                        onClick={() => navigate(`/manager/plans/${plan.id}`)}
                        className="btn btn-primary text-sm"
                      >
                        Открыть
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};