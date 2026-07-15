import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { skillsAPI } from '../../api/skills.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const EmployeeSkills = () => {
  const navigate = useNavigate();
  const [skills, setSkills] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadSkills();
  }, []);

  const loadSkills = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await skillsAPI.getMySkills();
      setSkills(response.data || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки навыков');
    } finally {
      setLoading(false);
    }
  };

  const getCategoryColor = (category) => {
    const colors = {
      frontend: 'bg-blue-100 text-blue-700',
      backend: 'bg-green-100 text-green-700',
      devops: 'bg-purple-100 text-purple-700',
      database: 'bg-yellow-100 text-yellow-700',
      testing: 'bg-red-100 text-red-700',
      cloud: 'bg-indigo-100 text-indigo-700',
      mobile: 'bg-pink-100 text-pink-700',
      architecture: 'bg-orange-100 text-orange-700',
      ai: 'bg-cyan-100 text-cyan-700',
      security: 'bg-rose-100 text-rose-700',
      soft_skills: 'bg-emerald-100 text-emerald-700',
      other: 'bg-gray-100 text-gray-700',
    };
    return colors[category] || colors.other;
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

  const formatDate = (dateString) => {
    if (!dateString) return '';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return '';
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      });
    } catch {
      return '';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="flex items-center justify-center h-64">
          <div className="text-gray-500">Загрузка навыков...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Мои навыки</h1>
          <span className="text-sm text-gray-500">Всего: {skills.length}</span>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

        {skills.length === 0 ? (
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-12 text-center">
            <div className="text-4xl mb-4">🎯</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">У вас пока нет навыков</h3>
            <p className="text-gray-500">
              Навыки будут появляться здесь после завершения планов и генерации через ИИ
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {skills.map((skill) => (
              <div
                key={skill.id}
                className="bg-white rounded-xl shadow-sm border border-gray-100 p-5 hover:shadow-md transition-shadow"
              >
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900 text-lg">{skill.name}</h3>
                    {skill.description && (
                      <p className="text-sm text-gray-600 mt-1">{skill.description}</p>
                    )}
                    <div className="flex flex-wrap items-center gap-3 mt-3">
                      {skill.category && (
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getCategoryColor(skill.category)}`}>
                          {getCategoryLabel(skill.category)}
                        </span>
                      )}
                      <span className="text-xs text-gray-400">
                        Получен: {formatDate(skill.confirmed_at || skill.created_at)}
                      </span>
                    </div>
                    {skill.plan_id && (
                      <button
                        onClick={() => navigate(`/employee/plans/${skill.plan_id}`)}
                        className="text-xs text-blue-600 hover:text-blue-800 mt-2 inline-block"
                      >
                        Из плана
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};