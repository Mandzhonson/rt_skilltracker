import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const CreateAIPlan = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [employees, setEmployees] = useState([]);
  const [loadingEmployees, setLoadingEmployees] = useState(true);
  const [formData, setFormData] = useState({
    employee_id: '',
    topic: '',
    description: '',
  });

  useEffect(() => {
    loadEmployees();
  }, []);

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
        } else if (response.data.data && Array.isArray(response.data.data)) {
          employeesData = response.data.data;
        } else {
          for (const key in response.data) {
            if (Array.isArray(response.data[key])) {
              employeesData = response.data[key];
              break;
            }
          }
        }
      }
      
      setEmployees(employeesData);
    } catch (err) {
      console.error('Error loading employees:', err);
      setError('Ошибка загрузки списка сотрудников');
    } finally {
      setLoadingEmployees(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      const response = await managerAPI.createAIPlan({
        employee_id: formData.employee_id,
        topic: formData.topic,
        description: formData.description || '',
      });
      navigate(`/manager/plans/${response.data.id}`);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания AI плана');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="max-w-2xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/manager/plans')}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад к списку
        </button>

        <div className="card">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Создать план с помощью ИИ</h1>
          <p className="text-gray-600 mb-4">ИИ создаст персонализированный план развития на основе темы и описания</p>

          {error && (
            <div className="error-message mb-4">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="label">Сотрудник</label>
              <select
                required
                value={formData.employee_id}
                onChange={(e) => setFormData({ ...formData, employee_id: e.target.value })}
                className="input"
                disabled={loadingEmployees}
              >
                <option value="">Выберите сотрудника</option>
                {employees.map((emp) => (
                  <option key={emp.id} value={emp.id}>
                    {emp.first_name} {emp.last_name} ({emp.email})
                  </option>
                ))}
              </select>
              {loadingEmployees && (
                <p className="text-sm text-gray-500 mt-1">Загрузка сотрудников...</p>
              )}
              {!loadingEmployees && employees.length === 0 && (
                <p className="text-sm text-yellow-600 mt-1">
                  У вас пока нет сотрудников в подчинении
                </p>
              )}
            </div>

            <div>
              <label className="label">Тема</label>
              <input
                type="text"
                required
                placeholder="Например: Go, React, Управление проектами"
                value={formData.topic}
                onChange={(e) => setFormData({ ...formData, topic: e.target.value })}
                className="input"
              />
            </div>

            <div>
              <label className="label">Описание</label>
              <textarea
                placeholder="Опишите цели и ожидания от плана (необязательно)"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input"
                rows="4"
              />
            </div>

            <div className="flex space-x-2">
              <button
                type="submit"
                className="btn btn-primary"
                disabled={loading || employees.length === 0}
              >
                {loading ? 'Создание...' : 'Создать план с ИИ'}
              </button>
              <button
                type="button"
                onClick={() => navigate('/manager/plans')}
                className="btn btn-secondary"
              >
                Отмена
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};