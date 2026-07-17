import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { adminAPI } from '../../api/admin.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const AdminManagerEmployees = () => {
  const { managerId } = useParams();
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadEmployees();
  }, [managerId]);

  const loadEmployees = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await adminAPI.listEmployeesByManager(managerId);
      setEmployees(response.data || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки сотрудников');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <div className="max-w-4xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">
          Сотрудники менеджера
        </h1>

        {error && (
          <div className="error-message mb-4">{error}</div>
        )}

        <div className="card">
          {loading ? (
            <div className="text-center py-8">Загрузка...</div>
          ) : employees.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Нет сотрудников</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Имя</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Роль</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {employees.map((emp) => (
                    <tr key={emp.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {emp.first_name} {emp.last_name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {emp.email}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-700">
                          Сотрудник
                        </span>
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