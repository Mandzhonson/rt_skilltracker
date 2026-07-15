import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { testAPI } from '../../api/test.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';

export const EmployeeTest = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [test, setTest] = useState(null);
  const [answers, setAnswers] = useState({});
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [result, setResult] = useState(null);

  useEffect(() => {
    loadTest();
  }, [planId]);

  const loadTest = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await testAPI.getTest(planId);
      setTest(response.data);
      const initialAnswers = {};
      response.data.questions.forEach(q => {
        initialAnswers[q.id] = '';
      });
      setAnswers(initialAnswers);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки теста');
    } finally {
      setLoading(false);
    }
  };

  const handleAnswerChange = (questionId, value) => {
    setAnswers(prev => ({
      ...prev,
      [questionId]: value
    }));
  };

  const handleSubmit = async () => {
    const unanswered = Object.values(answers).some(a => !a);
    if (unanswered) {
      setError('Ответьте на все вопросы перед отправкой');
      return;
    }

    setSubmitting(true);
    setError('');

    try {
      const formattedAnswers = Object.entries(answers).map(([questionId, answer]) => ({
        question_id: questionId,
        answer: answer
      }));

      const response = await testAPI.submitTest(planId, formattedAnswers);
      setResult(response.data);
      
      if (response.data.passed && response.data.score === response.data.total) {
        setTimeout(() => {
          window.location.href = `/employee/plans/${planId}`;
        }, 2000);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка отправки теста');
    } finally {
      setSubmitting(false);
    }
  };

  const handleRetry = () => {
    setResult(null);
    setAnswers({});
    const initialAnswers = {};
    test?.questions.forEach(q => {
      initialAnswers[q.id] = '';
    });
    setAnswers(initialAnswers);
    setError('');
  };

  const getOptionLabel = (index) => {
    const labels = ['A', 'B', 'C', 'D'];
    return labels[index] || index;
  };

  const allAnswered = Object.values(answers).every(a => a !== '');
  const totalQuestions = test?.questions?.length || 0;
  const answeredCount = Object.values(answers).filter(a => a !== '').length;

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="flex items-center justify-center h-64">Загрузка теста...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate(`/employee/plans/${planId}`)} className="btn btn-primary mt-4">
            Вернуться к плану
          </button>
        </div>
      </div>
    );
  }

  if (result) {
    let score = result.score;
    let total = result.total;
    
    if (score > total && total > 0) {
      const correctAnswers = Math.round((score / 100) * total);
      score = correctAnswers;
    }
    
    const percent = Math.round((score / total) * 100);
    const isPassed = percent >= 70;
    
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-3xl mx-auto px-4 py-8">
          <div className="card">
            <div className="text-center">
              <div className={`text-6xl mb-4 ${isPassed ? 'text-green-500' : 'text-red-500'}`}>
                {isPassed ? '✅' : '❌'}
              </div>
              <h2 className={`text-2xl font-bold mb-4 ${isPassed ? 'text-green-600' : 'text-red-600'}`}>
                {isPassed ? 'Тест пройден!' : 'Тест не пройден'}
              </h2>
              
              <div className="space-y-2 text-gray-700">
                <p className="text-lg">
                  Результат: {score} из {total} ({percent}%)
                </p>
                <p className="text-sm text-gray-500">
                  Правильных ответов: {score} из {total}
                </p>
                {!isPassed && (
                  <p className="text-sm text-gray-500 mt-4">
                    Для прохождения теста необходимо набрать 70% правильных ответов
                  </p>
                )}
              </div>
              
              <div className="flex gap-3 justify-center mt-6">
                <button
                  onClick={() => navigate(`/employee/plans/${planId}`)}
                  className="btn btn-secondary"
                >
                  Вернуться к плану
                </button>
                {!isPassed && (
                  <button
                    onClick={handleRetry}
                    className="btn btn-primary"
                  >
                    Перепройти тест
                  </button>
                )}
              </div>
            </div>
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
          onClick={() => navigate(`/employee/plans/${planId}`)}
          className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
        >
          ← Назад к плану
        </button>

        <div className="card">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">Тестирование</h1>
          
          {error && (
            <div className="error-message mb-4">{error}</div>
          )}

          <div className="flex justify-between items-center mb-4">
            <p className="text-sm text-gray-500">
              Для прохождения теста необходимо набрать 70% правильных ответов
            </p>
            <span className="text-sm text-gray-400">
              {answeredCount} из {totalQuestions}
            </span>
          </div>

          <div className="space-y-6">
            {test?.questions.map((question, index) => (
              <div key={question.id} className="border border-gray-200 rounded-lg p-4">
                <p className="font-medium text-gray-900 mb-3">
                  {index + 1}. {question.question}
                </p>
                <div className="space-y-2">
                  {question.options.map((option, optIndex) => {
                    const letter = getOptionLabel(optIndex);
                    return (
                      <label key={optIndex} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 cursor-pointer">
                        <input
                          type="radio"
                          name={`question_${question.id}`}
                          value={letter}
                          checked={answers[question.id] === letter}
                          onChange={() => handleAnswerChange(question.id, letter)}
                          className="w-4 h-4 text-blue-600"
                        />
                        <span className="text-gray-700">
                          {letter}. {option}
                        </span>
                      </label>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>

          <button
            onClick={handleSubmit}
            disabled={!allAnswered || submitting}
            className="btn btn-primary w-full mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {submitting ? 'Отправка...' : 'Отправить ответы'}
          </button>
          
          {!allAnswered && (
            <p className="text-sm text-red-500 mt-2 text-center">
              Ответьте на все вопросы ({totalQuestions - answeredCount} осталось)
            </p>
          )}
        </div>
      </div>
    </div>
  );
};