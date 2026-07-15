import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const ManagerPlans = () => {
  const navigate = useNavigate();
  const [plans, setPlans] = useState([]);
  const [filteredPlans, setFilteredPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [employees, setEmployees] = useState([]);
  const [loadingEmployees, setLoadingEmployees] = useState(false);
  const [selectedEmployeeId, setSelectedEmployeeId] = useState('');
  const [formData, setFormData] = useState({
    employee_id: '',
    title: '',
    description: '',
  });
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadPlans();
    loadEmployees();
  }, []);

  useEffect(() => {
    if (selectedEmployeeId) {
      setFilteredPlans(plans.filter(p => p.employee_id === selectedEmployeeId));
    } else {
      setFilteredPlans(plans);
    }
  }, [selectedEmployeeId, plans]);

  const loadPlans = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await managerAPI.listPlans();
      console.log('Plans response:', response.data);
      
      let plansData = [];
      if (response.data) {
        if (Array.isArray(response.data)) {
          plansData = response.data;
        } else if (response.data.plans && Array.isArray(response.data.plans)) {
          plansData = response.data.plans;
        } else {
          plansData = [];
        }
      }
      
      setPlans(plansData);
      setFilteredPlans(plansData);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки планов');
    } finally {
      setLoading(false);
    }
  };

  const loadEmployees = async () => {
    setLoadingEmployees(true);
    try {
      const response = await managerAPI.getEmployees();
      console.log('Employees response:', response.data);
      
      let employeesData = [];
      if (response.data) {
        if (Array.isArray(response.data)) {
          employeesData = response.data;
        } else if (response.data.employees && Array.isArray(response.data.employees)) {
          employeesData = response.data.employees;
        } else {
          employeesData = [];
        }
      }
      
      setEmployees(employeesData);
    } catch (err) {
      console.error('Error loading employees:', err);
    } finally {
      setLoadingEmployees(false);
    }
  };

  const handleShowCreateForm = () => {
    setShowCreateForm(true);
    if (employees.length === 0) {
      loadEmployees();
    }
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setCreating(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.createPlan({
        employee_id: formData.employee_id,
        title: formData.title,
        description: formData.description,
      });
      setMessage('План успешно создан');
      setFormData({ employee_id: '', title: '', description: '' });
      setShowCreateForm(false);
      loadPlans();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания плана');
    } finally {
      setCreating(false);
    }
  };

  const handleCreateAI = async (e) => {
    e.preventDefault();
    setCreating(true);
    setError('');
    setMessage('');

    try {
      await managerAPI.createAIPlan({
        employee_id: formData.employee_id,
        topic: formData.title,
        description: formData.description || '',
      });
      setMessage('План успешно создан с помощью ИИ');
      setFormData({ employee_id: '', title: '', description: '' });
      setShowCreateForm(false);
      loadPlans();
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания плана с помощью ИИ');
    } finally {
      setCreating(false);
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

  const getGenerationStatusLabel = (status) => {
    const statuses = {
      pending: '⏳ Ожидание',
      processing: '🔄 Генерация...',
      ready: '✅ Готово',
      failed: '❌ Ошибка',
    };
    return statuses[status] || status;
  };

  const getGenerationStatusColor = (status) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-700',
      processing: 'bg-blue-100 text-blue-700 animate-pulse',
      ready: 'bg-green-100 text-green-700',
      failed: 'bg-red-100 text-red-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-600';
  };

  const getEmployeeName = (employeeId) => {
    const employee = employees.find(e => e.id === employeeId);
    if (employee) {
      return `${employee.first_name} ${employee.last_name}`;
    }
    return employeeId;
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return '';
      return date.toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'short',
      });
    } catch {
      return '';
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

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Мои планы</h1>
          <button
            onClick={handleShowCreateForm}
            className="btn btn-primary"
          >
            Создать план
          </button>
        </div>

        {message && (
          <div className="success-message mb-4">{message}</div>
        )}
        {error && (
          <div className="error-message mb-4">{error}</div>
        )}

        {showCreateForm && (
          <div className="card mb-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Создание нового плана</h2>
            <form className="space-y-4">
              <div>
                <label className="label">Сотрудник</label>
                <select
                  required
                  value={formData.employee_id}
                  onChange={(e) => setFormData({ ...formData, employee_id: e.target.value })}
                  className="input"
                >
                  <option value="">Выберите сотрудника</option>
                  {loadingEmployees ? (
                    <option disabled>Загрузка...</option>
                  ) : employees.length === 0 ? (
                    <option disabled>Нет доступных сотрудников</option>
                  ) : (
                    employees.map((emp) => (
                      <option key={emp.id} value={emp.id}>
                        {emp.first_name} {emp.last_name}
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div>
                <label className="label">Название плана</label>
                <input
                  type="text"
                  required
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="input"
                  placeholder="Введите название плана"
                />
              </div>

              <div>
                <label className="label">Описание</label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="input"
                  rows="3"
                  placeholder="Введите описание плана"
                />
              </div>

              <div className="flex flex-wrap gap-3 pt-2">
                <button
                  type="button"
                  onClick={handleSubmit}
                  disabled={creating || !formData.employee_id || !formData.title}
                  className="btn btn-primary disabled:opacity-50"
                >
                  {creating ? 'Создание...' : 'Создать'}
                </button>
                <button
                  type="button"
                  onClick={handleCreateAI}
                  disabled={creating || !formData.employee_id || !formData.title}
                  className="btn btn-success disabled:opacity-50"
                >
                  {creating ? 'Создание...' : 'Создать с помощью ИИ'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowCreateForm(false);
                    setFormData({ employee_id: '', title: '', description: '' });
                    setError('');
                  }}
                  className="btn btn-secondary"
                >
                  Отмена
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="flex items-center gap-3 mb-4">
          <span className="text-sm text-gray-500">Фильтр:</span>
          <select
            value={selectedEmployeeId}
            onChange={(e) => setSelectedEmployeeId(e.target.value)}
            className="input max-w-xs"
          >
            <option value="">Все сотрудники</option>
            {employees.map((emp) => (
              <option key={emp.id} value={emp.id}>
                {emp.first_name} {emp.last_name}
              </option>
            ))}
          </select>
          {selectedEmployeeId && (
            <button
              onClick={() => setSelectedEmployeeId('')}
              className="text-sm text-red-600 hover:text-red-800"
            >
              Сбросить
            </button>
          )}
          <span className="text-sm text-gray-500 ml-auto">
            Всего: {filteredPlans.length}
          </span>
        </div>

        <div className="card">
          {filteredPlans.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {plans.length === 0 ? 'У вас пока нет планов. Создайте первый план!' : 'Планы не найдены'}
            </div>
          ) : (
            <div className="space-y-4">
              {filteredPlans.map((plan) => (
                <div
                  key={plan.id}
                  className="border rounded-lg p-4 hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => navigate(`/manager/plans/${plan.id}`)}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900">{plan.title}</h3>
                      {plan.description && (
                        <p className="text-sm text-gray-600 mt-1">{plan.description}</p>
                      )}
                      <div className="flex flex-wrap items-center gap-4 mt-2 text-sm">
                        <span className="text-gray-600">
                          Сотрудник: {getEmployeeName(plan.employee_id)}
                        </span>
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(plan.status)}`}>
                          {getStatusLabel(plan.status)}
                        </span>
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getGenerationStatusColor(plan.generation_status)}`}>
                          {getGenerationStatusLabel(plan.generation_status)}
                        </span>
                        <span className="text-gray-600">
                          Прогресс: {plan.progress || 0}%
                        </span>
                        <span className="text-gray-500 text-xs">
                          Создан: {formatDate(plan.created_at)}
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-2 ml-4 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                      <button
                        onClick={() => navigate(`/manager/plans/${plan.id}`)}
                        className="btn btn-primary text-sm"
                      >
                        Открыть
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};