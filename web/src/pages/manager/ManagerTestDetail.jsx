import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { managerAPI } from '../../api/manager.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const ManagerTestDetail = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [test, setTest] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadTest();
  }, [planId]);

  const loadTest = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await managerAPI.getTest(planId);
      setTest(response.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки теста');
    } finally {
      setLoading(false);
    }
  };

  const getOptionLabel = (index) => {
    const labels = ['A', 'B', 'C', 'D'];
    return labels[index] || index;
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
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate(-1)} className="btn btn-primary mt-4">
            Вернуться назад
          </button>
        </div>
      </div>
    );
  }

  if (!test || !test.questions || test.questions.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="card text-center py-12">
            <div className="text-4xl mb-3">📝</div>
            <h3 className="text-lg font-medium text-gray-900 mb-1">Тест не найден</h3>
            <p className="text-sm text-gray-500">Для этого плана ещё не создан тест</p>
            <button onClick={() => navigate(-1)} className="btn btn-primary mt-4">
              Вернуться назад
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-3xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate(-1)}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад
        </button>

        <div className="card">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-2xl font-bold text-gray-900">Тест</h1>
            <span className="text-sm text-gray-500">
              {test.questions.length} вопросов
            </span>
          </div>

          <div className="space-y-6">
            {test.questions.map((question, index) => (
              <div key={question.id} className="border border-gray-200 rounded-lg p-4">
                <p className="font-medium text-gray-900 mb-3">
                  {index + 1}. {question.question}
                </p>
                <div className="space-y-2">
                  {question.options.map((option, optIndex) => {
                    const letter = getOptionLabel(optIndex);
                    const isCorrect = question.correct_answer === letter;
                    return (
                      <div
                        key={optIndex}
                        className={`flex items-center gap-3 p-2 rounded-lg ${
                          isCorrect ? 'bg-green-50 border border-green-200' : 'bg-gray-50'
                        }`}
                      >
                        <span className="text-gray-700">
                          {letter}. {option}
                        </span>
                        {isCorrect && (
                          <span className="text-xs text-green-600 font-medium ml-auto">✓ Правильный</span>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};