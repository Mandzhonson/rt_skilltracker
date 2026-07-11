import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { DndContext, closestCorners, DragOverlay, useSensor, useSensors, PointerSensor } from '@dnd-kit/core';
import { employeeAPI } from '../../api/employee.js';
import { Navbar } from '../../components/Layout/Navbar.jsx';
import { TaskCard } from '../../components/employee/TaskCard.jsx';
import { TaskColumn } from '../../components/employee/TaskColumn.jsx';

export const EmployeeTasks = () => {
  const navigate = useNavigate();
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [updatingTaskId, setUpdatingTaskId] = useState(null);
  const [activeId, setActiveId] = useState(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await employeeAPI.listPlans();
      const plans = response.data || [];
      
      const allTasks = [];
      plans.forEach((plan) => {
        const planData = plan.plan || plan;
        const tasksData = planData.tasks || plan.tasks || [];
        
        tasksData.forEach((task) => {
          allTasks.push({
            ...task,
            plan_title: planData.title || 'Без названия',
            plan_id: planData.id || plan.id
          });
        });
      });
      
      setTasks(allTasks);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка загрузки задач');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (taskId, newStatus) => {
    if (updatingTaskId === taskId) return;
    
    setUpdatingTaskId(taskId);
    setError('');
    setMessage('');

    try {
      await employeeAPI.updateTaskStatus(taskId, newStatus);
      
      setTasks(prevTasks => 
        prevTasks.map(task => 
          task.id === taskId ? { ...task, status: newStatus } : task
        )
      );
      
      setMessage('Статус задачи обновлен');
      setTimeout(() => setMessage(''), 3000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка обновления статуса');
    } finally {
      setUpdatingTaskId(null);
    }
  };

  const handleDragEnd = async ({ active, over }) => {
    setActiveId(null);

    if (!over) return;

    const activeTask = tasks.find((t) => t.id === active.id);
    if (!activeTask) return;

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

  const getTasksByStatus = (status) => {
    return tasks.filter(t => t.status === status);
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

  const todoTasks = getTasksByStatus('todo');
  const inProgressTasks = getTasksByStatus('in_progress');
  const doneTasks = getTasksByStatus('done');

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

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Мои задачи</h1>
          <span className="text-sm text-gray-500">Всего: {tasks.length}</span>
        </div>

        {message && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-4">
            {message}
          </div>
        )}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

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
      </div>
    </div>
  );
};