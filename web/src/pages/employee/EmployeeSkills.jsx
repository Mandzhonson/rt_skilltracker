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
      const skillsData = response.data?.skills || response.data || [];
      setSkills(skillsData);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки навыков');
    } finally {
      setLoading(false);
    }
  };

  const getCategoryColor = (category) => {
    const colors = {
      frontend: 'bg-blue-50 text-blue-700',
      backend: 'bg-emerald-50 text-emerald-700',
      devops: 'bg-purple-50 text-purple-700',
      database: 'bg-amber-50 text-amber-700',
      testing: 'bg-rose-50 text-rose-700',
      cloud: 'bg-indigo-50 text-indigo-700',
      mobile: 'bg-pink-50 text-pink-700',
      architecture: 'bg-orange-50 text-orange-700',
      ai: 'bg-cyan-50 text-cyan-700',
      security: 'bg-red-50 text-red-700',
      soft_skills: 'bg-teal-50 text-teal-700',
      other: 'bg-gray-50 text-gray-700',
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
        month: 'short',
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
          <div className="text-gray-500">Загрузка...</div>
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
          <div className="error-message mb-4">{error}</div>
        )}

        {skills.length === 0 ? (
          <div className="card text-center py-12 border border-black rounded-lg bg-white">
            <div className="text-4xl mb-3">🎯</div>
            <h3 className="text-lg font-medium text-gray-900 mb-1">У вас пока нет навыков</h3>
            <p className="text-sm text-gray-500">
              Навыки будут появляться здесь после завершения планов
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
           {skills.map((skill) => (
              <div
                key={skill.id}
                className="bg-white rounded-lg border-card card-padding text-left"
              >
                <div className="item-spacing">
                  {/* Название навыка */}
                  <p className="font-medium text-gray-900 text-sm">
                    {skill.name}
                  </p>
                  
                  {/* Описание навыка */}
                  {skill.description && (
                    <p className="text-sm text-gray-500">
                      {skill.description}
                    </p>
                  )}
                  
                  {/* Блок с категорией и датой */}
                  <div className="tags-container">
                    {skill.category && (
                      <span className={`px-2 py-0.5 rounded text-xs font-medium ${getCategoryColor(skill.category)}`}>
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
        )}
      </div>
    </div>
  );
};