import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

export const TaskCard = ({ 
  task, 
  onNavigate, 
  formatDate
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: task.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
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
      className="bg-gray-50 border border-gray-200 rounded-lg p-3 hover:shadow-sm transition-shadow cursor-grab active:cursor-grabbing"
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p className="font-medium text-gray-800 text-sm">{task.title}</p>
            <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(task.status)}`}>
              {getStatusLabel(task.status)}
            </span>
          </div>
          {task.description && (
            <p className="text-xs text-gray-500 mt-1">{task.description}</p>
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
    </div>
  );
};