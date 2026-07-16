import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { employeeAPI } from '../../api/employee.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const EmployeePlans = () => {
  const navigate = useNavigate();
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadPlans();
  }, []);

  const loadPlans = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await employeeAPI.listPlans();
      const plansData = response.data || [];
      // Фильтруем архивированные планы
      const activePlans = plansData.filter(plan => {
        const planData = plan.plan || plan;
        return planData.status !== 'archived';
      });
      setPlans(activePlans);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки планов');
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
      draft: 'bg-gray-100 text-gray-700',
      active: 'bg-blue-100 text-blue-700',
      completed: 'bg-green-100 text-green-700',
      archived: 'bg-yellow-100 text-yellow-700',
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

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-5xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Мои планы</h1>

        {error && (
          <div className="error-message mb-4">{error}</div>
        )}

        <div className="card">
          {plans.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              У вас пока нет планов
            </div>
          ) : (
            <div className="space-y-4">
              {plans.map((plan) => {
                const planData = plan.plan || plan;
                const planId = planData.id || plan.id;
                
                return (
                  <div
                    key={planId}
                    className="border rounded-lg p-4 hover:bg-gray-50 cursor-pointer transition-colors"
                    onClick={() => {
                      if (planId) {
                        navigate(`/employee/plans/${planId}`);
                      }
                    }}
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900">{planData.title || 'Без названия'}</h3>
                        {planData.description && (
                          <p className="text-sm text-gray-600 mt-1">{planData.description}</p>
                        )}
                        <div className="flex flex-wrap items-center gap-4 mt-2 text-sm">
                          <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(planData.status)}`}>
                            {getStatusLabel(planData.status)}
                          </span>
                          <span className="text-gray-600">
                            Прогресс: {planData.progress || 0}%
                          </span>
                          <span className="text-gray-500 text-xs">
                            Создан: {formatDate(planData.created_at)}
                          </span>
                        </div>
                      </div>
                      <div className="flex gap-2 ml-4 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                        <button
                          onClick={() => {
                            if (planId) {
                              navigate(`/employee/plans/${planId}`);
                            }
                          }}
                          className="btn btn-primary text-sm"
                        >
                          Открыть
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};