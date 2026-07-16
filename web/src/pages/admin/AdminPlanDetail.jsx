import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { adminAPI } from '../../api/admin.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const AdminPlanDetail = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [plan, setPlan] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadPlan();
  }, [planId]);

  const loadPlan = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await adminAPI.getPlan(planId);
      const data = response.data;
      const planData = data.plan || data;
      setPlan(planData);
      setTasks(planData.tasks || data.tasks || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки плана');
    } finally {
      setLoading(false);
    }
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

  const getTaskStatusLabel = (status) => {
    const statuses = {
      todo: 'К выполнению',
      in_progress: 'В работе',
      done: 'Выполнена',
    };
    return statuses[status] || status;
  };

  const getTaskStatusColor = (status) => {
    const colors = {
      todo: 'bg-gray-100 text-gray-700',
      in_progress: 'bg-blue-100 text-blue-700',
      done: 'bg-green-100 text-green-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-700';
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
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate('/admin/users')} className="btn btn-primary mt-4">
            Вернуться к пользователям
          </button>
        </div>
      </div>
    );
  }

  if (!plan) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="text-center">План не найден</div>
          <button onClick={() => navigate('/admin/users')} className="btn btn-primary mt-4">
            Вернуться к пользователям
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-4xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate(-1)}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад
        </button>

        {/* Общая информация о плане */}
        <div className="card mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">{plan.title}</h1>
            {plan.description && (
              <p className="text-gray-600 mb-4">{plan.description}</p>
            )}
            <div className="flex flex-wrap items-center gap-4 text-sm">
              <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(plan.status)}`}>
                Статус: {getStatusLabel(plan.status)}
              </span>
              <span className="text-gray-600">
                Прогресс: {plan.progress || 0}%
              </span>
              <span className="text-gray-500">
                Создан: {formatDate(plan.created_at)}
              </span>
            </div>
          </div>
        </div>

        {/* Список задач */}
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Задачи</h2>
        <div className="card mb-6">
          {tasks.length === 0 ? (
            <p className="text-gray-500 text-center py-4">Задачи отсутствуют</p>
          ) : (
            <div className="space-y-3">
              {tasks.map((task) => (
                <div
                  key={task.id}
                  className="border rounded-lg p-4 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex justify-between items-center gap-4">
                    <div className="flex-1">
                      <h3 className="font-medium text-gray-900">{task.title}</h3>
                      {task.description && (
                        <p className="text-sm text-gray-500 mt-1">{task.description}</p>
                      )}
                      <div className="flex flex-wrap items-center gap-3 mt-2 text-sm">
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getTaskStatusColor(task.status)}`}>
                          {getTaskStatusLabel(task.status)}
                        </span>
                        <span className="text-gray-400 text-xs">
                          Создана: {formatDate(task.created_at)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Вынесенная кнопка просмотра теста в самом конце */}
        <div className="flex justify-end">
          <button
            onClick={() => navigate(`/admin/plans/${planId}/test`)}
            className="btn !bg-purple-600 hover:!bg-purple-700 !text-white font-medium px-6 py-2.5 shadow-sm transition-all"
          >
            Просмотр теста плана
          </button>
        </div>
      </div>
    </div>
  );
};