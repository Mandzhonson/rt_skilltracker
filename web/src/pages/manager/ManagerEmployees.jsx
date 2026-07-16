import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const ManagerEmployees = () => {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState([]);
  const [avatars, setAvatars] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadEmployees();
  }, []);

  const loadEmployees = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await managerAPI.getEmployees();
      let employeesData = [];
      if (response.data) {
        if (Array.isArray(response.data)) {
          employeesData = response.data;
        } else if (response.data.employees && Array.isArray(response.data.employees)) {
          employeesData = response.data.employees;
        }
      }
      setEmployees(employeesData);
      
      if (employeesData.length > 0) {
        loadAllAvatars(employeesData);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки сотрудников');
    } finally {
      setLoading(false);
    }
  };

  const loadAllAvatars = async (employeesList) => {
    const avatarPromises = employeesList.map(async (emp) => {
      try {
        const response = await managerAPI.getEmployeeAvatar(emp.id);
        const url = URL.createObjectURL(response.data);
        return { id: emp.id, url };
      } catch (err) {
        if (err.response?.status !== 404) {
          console.error(`Error loading avatar for employee ${emp.id}:`, err);
        }
        return { id: emp.id, url: null };
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
    }
  };

  const getInitials = (firstName, lastName) => {
    return `${firstName?.[0] || ''}${lastName?.[0] || ''}`;
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

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-6xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Мои сотрудники</h1>

        {error && (
          <div className="error-message mb-4">{error}</div>
        )}

        <div className="card">
          {employees.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              У вас пока нет сотрудников
            </div>
          ) : (
            <div className="space-y-3">
              {employees.map((employee) => (
                <div
                  key={employee.id}
                  className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => navigate(`/manager/employees/${employee.id}`)}
                >
                  <div className="flex items-center gap-4">
                    {avatars[employee.id] ? (
                      <img
                        src={avatars[employee.id]}
                        alt="Avatar"
                        className="avatar-sm" /* Используем стандартизированный класс .avatar-sm */
                      />
                    ) : (
                      <div className="avatar-placeholder-sm"> {/* Используем класс .avatar-placeholder-sm */}
                        {getInitials(employee.first_name, employee.last_name)}
                      </div>
                    )}
                    <div>
                      <p className="font-medium text-gray-900">
                        {employee.first_name} {employee.last_name}
                      </p>
                      <p className="text-sm text-gray-500">{employee.email}</p>
                    </div>
                  </div>
                  <div className="text-sm text-gray-400">
                    {employee.position || 'Сотрудник'}
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