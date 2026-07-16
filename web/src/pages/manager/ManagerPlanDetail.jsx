import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const ManagerPlanDetail = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [plan, setPlan] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [showTaskForm, setShowTaskForm] = useState(false);
  const [taskForm, setTaskForm] = useState({
    title: '',
    description: '',
  });
  const [loadingTask, setLoadingTask] = useState(false);
  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editTaskForm, setEditTaskForm] = useState({
    title: '',
    description: '',
  });
  const [editingPlan, setEditingPlan] = useState(false);
  const [editPlanForm, setEditPlanForm] = useState({
    title: '',
    description: '',
  });

  useEffect(() => {
    loadPlan();
  }, [planId]);

  const loadPlan = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await managerAPI.getPlan(planId);
      setPlan(response.data);
      setTasks(response.data.tasks || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки плана');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (e) => {
    e.preventDefault();
    if (planData.status === 'archived') {
      setError('Нельзя добавлять задачи в архивированный план');
      return;
    }
    
    setLoadingTask(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.createTask(planId, {
        title: taskForm.title,
        description: taskForm.description,
      });
      setMessage('Задача успешно создана');
      setTaskForm({ title: '', description: '' });
      setShowTaskForm(false);
      loadPlan();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания задачи');
    } finally {
      setLoadingTask(false);
    }
  };

  const handleUpdateTask = async (e) => {
    e.preventDefault();
    if (planData.status === 'archived') {
      setError('Нельзя редактировать задачи в архивированном плане');
      return;
    }
    
    setLoadingTask(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.updateTask(editingTaskId, {
        title: editTaskForm.title,
        description: editTaskForm.description,
      });
      setMessage('Задача обновлена');
      setEditingTaskId(null);
      setEditTaskForm({ title: '', description: '' });
      loadPlan();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка обновления задачи');
    } finally {
      setLoadingTask(false);
    }
  };

  const handleDeleteTask = async (taskId) => {
    if (planData.status === 'archived') {
      setError('Нельзя удалять задачи из архивированного плана');
      return;
    }
    
    if (!confirm('Удалить задачу?')) return;
    try {
      await managerAPI.deleteTask(taskId);
      setMessage('Задача удалена');
      loadPlan();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка удаления задачи');
    }
  };

  const handleDeletePlan = async () => {
    if (!confirm('Удалить план? Все задачи также будут удалены.')) return;
    setLoading(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.deletePlan(planId);
      setMessage('План успешно удален');
      setTimeout(() => {
        navigate('/manager/plans');
      }, 1500);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка удаления плана');
      setLoading(false);
    }
  };

  const handleUpdatePlan = async (e) => {
    e.preventDefault();
    if (planData.status === 'archived') {
      setError('Нельзя редактировать архивированный план');
      return;
    }
    
    setLoading(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.updatePlan(planId, {
        title: editPlanForm.title,
        description: editPlanForm.description || '',
      });
      setMessage('План успешно обновлен');
      setEditingPlan(false);
      loadPlan();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка обновления плана');
    } finally {
      setLoading(false);
    }
  };

  const handleArchivePlan = async () => {
    if (!confirm('Архивировать план? Он станет недоступен для сотрудника.')) return;
    setLoading(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.archivePlan(planId);
      setMessage('План успешно архивирован');
      setTimeout(() => {
        navigate('/manager/plans');
      }, 1500);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка архивации плана');
      setLoading(false);
    }
  };

  const startEditingPlan = () => {
    const planData = plan.plan || plan;
    setEditPlanForm({
      title: planData.title || '',
      description: planData.description || '',
    });
    setEditingPlan(true);
  };

  const cancelEditingPlan = () => {
    setEditingPlan(false);
    setEditPlanForm({ title: '', description: '' });
  };

  const startEditingTask = (task) => {
    setEditingTaskId(task.id);
    setEditTaskForm({
      title: task.title,
      description: task.description || '',
    });
  };

  const cancelEditingTask = () => {
    setEditingTaskId(null);
    setEditTaskForm({ title: '', description: '' });
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
      draft: 'bg-gray-100 text-gray-700',
      active: 'bg-blue-100 text-blue-700',
      completed: 'bg-green-100 text-green-700',
      archived: 'bg-yellow-100 text-yellow-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-700';
  };

  const getTaskStatusLabel = (status) => {
    const statuses = {
      pending: 'Ожидает',
      in_progress: 'В работе',
      done: 'Выполнена',
    };
    return statuses[status] || status;
  };

  const getTaskStatusColor = (status) => {
    const colors = {
      pending: 'bg-gray-100 text-gray-700',
      in_progress: 'bg-blue-100 text-blue-700',
      done: 'bg-green-100 text-green-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-700';
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'Не указана';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return 'Не указана';
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      });
    } catch {
      return 'Не указана';
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
          <button onClick={() => navigate('/manager/plans')} className="btn btn-primary mt-4">
            Вернуться к планам
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
          <button onClick={() => navigate('/manager/plans')} className="btn btn-primary mt-4">
            Вернуться к планам
          </button>
        </div>
      </div>
    );
  }

  const planData = plan.plan || plan;
  const isArchived = planData.status === 'archived';

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-4xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/manager/plans')}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад к планам
        </button>

        {message && (
          <div className="success-message mb-4">{message}</div>
        )}
        {error && (
          <div className="error-message mb-4">{error}</div>
        )}

        <div className="card mb-6">
          {editingPlan ? (
            <form onSubmit={handleUpdatePlan} className="space-y-4">
              <div>
                <label className="label">Название плана</label>
                <input
                  type="text"
                  required
                  value={editPlanForm.title}
                  onChange={(e) => setEditPlanForm({ ...editPlanForm, title: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="label">Описание</label>
                <textarea
                  value={editPlanForm.description}
                  onChange={(e) => setEditPlanForm({ ...editPlanForm, description: e.target.value })}
                  className="input"
                  rows="3"
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  disabled={loading || isArchived}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {loading ? 'Сохранение...' : 'Сохранить'}
                </button>
                <button
                  type="button"
                  onClick={cancelEditingPlan}
                  className="btn btn-secondary"
                >
                  Отмена
                </button>
              </div>
              {isArchived && (
                <p className="text-sm text-yellow-600">Архивированный план нельзя редактировать</p>
              )}
            </form>
          ) : (
            <>
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <h1 className="text-2xl font-bold text-gray-900 mb-2">{planData.title}</h1>
                  {planData.description && (
                    <p className="text-gray-600 mb-4">{planData.description}</p>
                  )}
                  <div className="flex flex-wrap items-center gap-4 text-sm">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(planData.status)}`}>
                      Статус: {getStatusLabel(planData.status)}
                    </span>
                    <span className="text-gray-600">
                      Прогресс: {planData.progress || 0}%
                    </span>
                    <span className="text-gray-500">
                      Создан: {formatDate(planData.created_at)}
                    </span>
                  </div>
                </div>
                <div className="flex flex-wrap gap-2 ml-4 flex-shrink-0">
                  {!isArchived && (
                    <>
                      <button
                        onClick={startEditingPlan}
                        className="px-3 py-1 text-sm text-blue-600 hover:text-blue-800 bg-blue-50 hover:bg-blue-100 rounded"
                      >
                        Редактировать
                      </button>
                      <button
                        onClick={handleArchivePlan}
                        className="px-3 py-1 text-sm text-yellow-600 hover:text-yellow-800 bg-yellow-50 hover:bg-yellow-100 rounded"
                      >
                        Архивировать
                      </button>
                    </>
                  )}
                  <button
                    onClick={handleDeletePlan}
                    className="px-3 py-1 text-sm text-red-600 hover:text-red-800 bg-red-50 hover:bg-red-100 rounded"
                  >
                    Удалить
                  </button>
                </div>
              </div>
            </>
          )}
        </div>

        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-gray-900">Задачи</h2>
          {!isArchived && (
            <button
              onClick={() => setShowTaskForm(!showTaskForm)}
              className="btn btn-primary"
            >
              Добавить задачу
            </button>
          )}
          {isArchived && (
            <span className="text-sm text-yellow-600">Архивированный план (только просмотр)</span>
          )}
        </div>

        {!isArchived && showTaskForm && (
          <div className="card mb-4">
            <h3 className="font-semibold text-gray-900 mb-3">Новая задача</h3>
            <form onSubmit={handleCreateTask} className="space-y-4">
              <div>
                <label className="label">Название</label>
                <input
                  type="text"
                  required
                  value={taskForm.title}
                  onChange={(e) => setTaskForm({ ...taskForm, title: e.target.value })}
                  className="input"
                  placeholder="Введите название задачи"
                />
              </div>
              <div>
                <label className="label">Описание</label>
                <textarea
                  value={taskForm.description}
                  onChange={(e) => setTaskForm({ ...taskForm, description: e.target.value })}
                  className="input"
                  rows="2"
                  placeholder="Введите описание задачи"
                />
              </div>
              <div className="flex gap-3">
                <button
                  type="submit"
                  disabled={loadingTask}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {loadingTask ? 'Создание...' : 'Создать'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowTaskForm(false);
                    setTaskForm({ title: '', description: '' });
                  }}
                  className="btn btn-secondary"
                >
                  Отмена
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Список задач */}
        <div className="card mb-6">
          {tasks.length > 0 ? (
            <div className="space-y-4">
              {tasks.map((task) => (
                <div
                  key={task.id}
                  className="border rounded-lg p-4 hover:bg-gray-50"
                >
                  {editingTaskId === task.id && !isArchived ? (
                    <form onSubmit={handleUpdateTask} className="space-y-3">
                      <div>
                        <label className="label">Название</label>
                        <input
                          type="text"
                          required
                          value={editTaskForm.title}
                          onChange={(e) => setEditTaskForm({ ...editTaskForm, title: e.target.value })}
                          className="input"
                        />
                      </div>
                      <div>
                        <label className="label">Описание</label>
                        <textarea
                          value={editTaskForm.description}
                          onChange={(e) => setEditTaskForm({ ...editTaskForm, description: e.target.value })}
                          className="input"
                          rows="2"
                        />
                      </div>
                      <div className="flex gap-3">
                        <button
                          type="submit"
                          disabled={loadingTask}
                          className="btn btn-primary disabled:opacity-50"
                        >
                          {loadingTask ? 'Сохранение...' : 'Сохранить'}
                        </button>
                        <button
                          type="button"
                          onClick={cancelEditingTask}
                          className="btn btn-secondary"
                        >
                          Отмена
                        </button>
                      </div>
                    </form>
                  ) : (
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h4 className="font-medium text-gray-900">{task.title}</h4>
                        {task.description && (
                          <p className="text-sm text-gray-600 mt-1">{task.description}</p>
                        )}
                        <div className="flex flex-wrap items-center gap-4 mt-2">
                          <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getTaskStatusColor(task.status)}`}>
                            {getTaskStatusLabel(task.status)}
                          </span>
                          <span className="text-xs text-gray-400">
                            Создана: {formatDate(task.created_at)}
                          </span>
                        </div>
                      </div>
                      {!isArchived && (
                        <div className="flex gap-2 ml-4 flex-shrink-0">
                          <button
                            onClick={() => startEditingTask(task)}
                            className="px-3 py-1 text-sm text-blue-600 hover:text-blue-800 bg-blue-50 hover:bg-blue-100 rounded"
                          >
                            Редактировать
                          </button>
                          <button
                            onClick={() => handleDeleteTask(task.id)}
                            className="px-3 py-1 text-sm text-red-600 hover:text-red-800 bg-red-50 hover:bg-red-100 rounded"
                          >
                            Удалить
                          </button>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">Задачи отсутствуют</p>
          )}
        </div>

        {/* Вынесенная кнопка просмотра теста менеджера в самом конце страницы */}
        <div className="flex justify-end">
          <button
            onClick={() => navigate(`/manager/plans/${planId}/test`)}
            className="btn !bg-purple-600 hover:!bg-purple-700 !text-white font-medium px-6 py-2.5 shadow-sm transition-all"
          >
            Просмотр теста плана
          </button>
        </div>
      </div>
    </div>
  );
};