import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { DndContext, closestCorners, DragOverlay, useSensor, useSensors, PointerSensor } from '@dnd-kit/core';
import { employeeAPI } from '../../api/employee.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';
import { TaskCard } from '../../components/employee/TaskCard.jsx';
import { TaskColumn } from '../../components/employee/TaskColumn.jsx';
import { SkillsGenerationModal } from '../../components/employee/SkillsGenerationModal.jsx';

export const EmployeePlanDetail = () => {
  const { planId } = useParams();
  const navigate = useNavigate();
  const [plan, setPlan] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [updatingTaskId, setUpdatingTaskId] = useState(null);
  const [activeId, setActiveId] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [generationMessage, setGenerationMessage] = useState('');
  const [isGenerationComplete, setIsGenerationComplete] = useState(false);
  const [isArchived, setIsArchived] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  useEffect(() => {
    if (planId) {
      loadPlan();
    }
  }, [planId]);

  const loadPlan = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await employeeAPI.getPlan(planId);
      const planData = response.data;
      
      if (planData.status === 'archived') {
        setIsArchived(true);
        setPlan(planData);
        setTasks([]);
        setLoading(false);
        return;
      }
      
      setIsArchived(false);
      setPlan(planData);
      const tasksData = planData.tasks || planData.plan?.tasks || [];
      setTasks(tasksData);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки плана');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (taskId, newStatus) => {
    if (updatingTaskId === taskId || isArchived) return;
    
    setUpdatingTaskId(taskId);
    setError('');
    setMessage('');

    const previousTasks = [...tasks];
    const updatedTasks = tasks.map(task => 
      task.id === taskId ? { ...task, status: newStatus } : task
    );
    setTasks(updatedTasks);

    try {
      await employeeAPI.updateTaskStatus(taskId, newStatus);
      
      const allDone = updatedTasks.every(task => task.status === 'done');
      const hasTasks = updatedTasks.length > 0;
      
      if (allDone && hasTasks) {
        setIsGenerationComplete(false);
        setGenerationMessage('Ваши навыки формируются на основе завершенного плана. Это может занять несколько минут. Вы можете продолжить работу.');
        setIsModalOpen(true);
        
        let attempts = 0;
        const maxAttempts = 60;
        
        const checkStatus = async () => {
          try {
            const status = await employeeAPI.getSkillsStatus(planId);
            if (status.data.isGenerated && status.data.count > 0) {
              setGenerationMessage(`Успешно сформировано ${status.data.count} навыков!`);
              setIsGenerationComplete(true);
              return true;
            }
            return false;
          } catch (err) {
            return false;
          }
        };

        const interval = setInterval(async () => {
          attempts++;
          if (attempts >= maxAttempts) {
            clearInterval(interval);
            setGenerationMessage('Генерация навыков может занять больше времени. Вы можете проверить их позже в разделе "Мои навыки".');
            setIsGenerationComplete(true);
            return;
          }

          const done = await checkStatus();
          if (done) {
            clearInterval(interval);
          }
        }, 3000);

        await checkStatus();
      } else {
        setMessage('Статус задачи обновлен');
        setTimeout(() => setMessage(''), 3000);
      }
    } catch (err) {
      setTasks(previousTasks);
      setError(err.response?.data?.error || 'Ошибка обновления статуса');
      setTimeout(() => setError(''), 3000);
    } finally {
      setUpdatingTaskId(null);
    }
  };

  const handleDragEnd = async ({ active, over }) => {
    setActiveId(null);

    if (!over || isArchived) return;

    const activeTask = tasks.find((t) => t.id === active.id);
    if (!activeTask) return;

    if (activeTask.status === 'done') {
      setMessage('Нельзя перемещать выполненные задачи');
      setTimeout(() => setMessage(''), 3000);
      return;
    }

    if (activeTask.title === 'Пройти тестирование') {
      setMessage('Для завершения этой задачи необходимо пройти тестирование');
      setTimeout(() => setMessage(''), 3000);
      return;
    }

    let newStatus = null;

    if (['todo', 'in_progress', 'done'].includes(over.id)) {
      newStatus = over.id;
    } else {
      const overTask = tasks.find((t) => t.id === over.id);
      if (overTask) {
        newStatus = overTask.status;
      }
    }

    if (!newStatus || newStatus === activeTask.status) {
      return;
    }

    await handleUpdateStatus(activeTask.id, newStatus);
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

  const planData = plan?.plan || plan || {};

  const todoTasks = tasks.filter(task => task.status === 'todo');
  const inProgressTasks = tasks.filter(task => task.status === 'in_progress');
  const doneTasks = tasks.filter(task => task.status === 'done');

  const hasTestTask = tasks.some(task => 
    task.title === 'Пройти тестирование' && task.status !== 'done'
  );

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
        <div className="max-w-6xl mx-auto px-4 py-8">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate('/employee/plans')} className="btn btn-primary mt-4">
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
        <div className="max-w-6xl mx-auto px-4 py-8">
          <div className="text-center">План не найден</div>
          <button onClick={() => navigate('/employee/plans')} className="btn btn-primary mt-4">
            Вернуться к планам
          </button>
        </div>
      </div>
    );
  }

  if (isArchived) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <div className="max-w-6xl mx-auto px-4 py-8">
          <button
            onClick={() => navigate('/employee/plans')}
            className="text-blue-600 hover:text-blue-800 mb-4 inline-block"
          >
            ← Назад к планам
          </button>
          <div className="card">
            <div className="text-center py-12">
              <div className="text-6xl mb-4">📦</div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">План архивирован</h2>
              <p className="text-gray-500">Этот план был архивирован менеджером и недоступен для просмотра</p>
              <button
                onClick={() => navigate('/employee/plans')}
                className="btn btn-primary mt-4"
              >
                Вернуться к планам
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <div className="max-w-6xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/employee/plans')}
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
          <h1 className="text-2xl font-bold text-gray-900 mb-2">{planData.title || 'Без названия'}</h1>
          {planData.description && (
            <p className="text-gray-600 mb-4">{planData.description}</p>
          )}
          <div className="flex flex-wrap items-center gap-3 text-sm">
            <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(planData.status)}`}>
              Статус: {getStatusLabel(planData.status)}
            </span>
            <span className="text-gray-600">
              Прогресс: {planData.progress || 0}%
            </span>
            <span className="text-gray-500 text-xs">
              {formatDate(planData.created_at)}
            </span>
            {hasTestTask && (
              <button
                onClick={() => navigate(`/employee/plans/${planId}/test`)}
                className="!bg-gray-800 !text-white px-3 py-1 rounded text-xs font-medium hover:!bg-gray-900 shadow-sm"
              >
                Пройти тест
              </button>
            )}
          </div>
        </div>

        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={({ active }) => setActiveId(active.id)}
          onDragEnd={handleDragEnd}
        >
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <TaskColumn
              id="todo"
              title="К выполнению"
              icon="📋"
              tasks={todoTasks}
              textColor="text-gray-700"
              countBg="bg-gray-200 text-gray-600"
            >
              {todoTasks.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">Нет задач</p>
              ) : (
                todoTasks.map((task) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    onNavigate={navigate}
                    formatDate={formatDate}
                  />
                ))
              )}
            </TaskColumn>

            <TaskColumn
              id="in_progress"
              title="В работе"
              icon="🔄"
              tasks={inProgressTasks}
              textColor="text-blue-700"
              countBg="bg-blue-100 text-blue-700"
            >
              {inProgressTasks.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">Нет задач</p>
              ) : (
                inProgressTasks.map((task) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    onNavigate={navigate}
                    formatDate={formatDate}
                  />
                ))
              )}
            </TaskColumn>

            <TaskColumn
              id="done"
              title="Выполнены"
              icon="✅"
              tasks={doneTasks}
              textColor="text-green-700"
              countBg="bg-green-100 text-green-700"
            >
              {doneTasks.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">Нет задач</p>
              ) : (
                doneTasks.map((task) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    onNavigate={navigate}
                    formatDate={formatDate}
                  />
                ))
              )}
            </TaskColumn>
          </div>

          <DragOverlay>
            {activeId ? (
              <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 shadow-lg opacity-90">
                {tasks.find(t => t.id === activeId)?.title}
              </div>
            ) : null}
          </DragOverlay>
        </DndContext>

        <SkillsGenerationModal
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false);
            if (isGenerationComplete) {
              window.location.reload();
            }
          }}
          message={generationMessage}
          isComplete={isGenerationComplete}
        />
      </div>
    </div>
  );
};