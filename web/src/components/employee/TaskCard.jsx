import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

export const TaskCard = ({ 
  task, 
  onNavigate, 
  formatDate
}) => {
  const isTestTask = task.title === 'Пройти тестирование';
  const isDone = task.status === 'done';
  const isDisabled = isDone || isTestTask;

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ 
    id: task.id,
    disabled: isDisabled,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    cursor: isDisabled ? 'default' : 'grab',
  };

  const getStatusLabel = (status) => {
    const statuses = {
      todo: 'К выполнению',
      in_progress: 'В работе',
      done: 'Выполнена',
    };
    return statuses[status] || status;
  };

  const getStatusColor = (status) => {
    const colors = {
      todo: 'bg-gray-100 text-gray-700',
      in_progress: 'bg-blue-100 text-blue-700',
      done: 'bg-green-100 text-green-700',
    };
    return colors[status] || 'bg-gray-100 text-gray-700';
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`bg-gray-50 border border-gray-200 rounded-lg p-3 hover:shadow-sm transition-shadow ${
        isDisabled ? 'opacity-75 cursor-default' : 'cursor-grab active:cursor-grabbing'
      }`}
    >
      <div className="flex flex-col gap-2">
        <div className="flex items-start justify-between gap-2">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <p className={`font-medium text-sm ${isDone ? 'text-gray-500 line-through' : 'text-gray-800'}`}>
                {task.title}
              </p>
              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(task.status)}`}>
                {getStatusLabel(task.status)}
              </span>
              {isTestTask && !isDone && (
                <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-700">
                  Тест
                </span>
              )}
            </div>
            {task.description && (
              <p className={`text-xs mt-1 ${isDone ? 'text-gray-400 line-through' : 'text-gray-500'}`}>
                {task.description}
              </p>
            )}
            <div className="flex items-center gap-2 mt-2 text-xs text-gray-400">
              <button
                onClick={() => onNavigate(`/employee/plans/${task.plan_id}`)}
                className="text-blue-600 hover:text-blue-800"
              >
                {task.plan_title}
              </button>
              <span>•</span>
              <span>{formatDate(task.created_at)}</span>
            </div>
          </div>
        </div>
        
        {isTestTask && !isDone && (
          <div className="mt-2 pt-2 border-t border-gray-200">
            <button
              onClick={() => onNavigate(`/employee/plans/${task.plan_id}/test`)}
              className="!bg-blue-600 !text-white px-4 py-2.5 rounded-lg text-sm font-medium hover:!bg-blue-700 shadow-sm w-full"
            >
              Пройти тестирование
            </button>
          </div>
        )}
      </div>
    </div>
  );
};