import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';

export const TaskColumn = ({
  id,
  title,
  icon,
  tasks,
  children,
  bg,
  headerBg,
  countBg,
  textColor,
}) => {
  const { setNodeRef, isOver } = useDroppable({
    id,
  });

  return (
    <div className="card">
      <div className={`flex justify-between items-center mb-3 ${headerBg}`}>
        <h3 className={`font-semibold ${textColor || 'text-gray-700'}`}>{title}</h3>
        <span className={`text-xs ${countBg || 'bg-gray-200 text-gray-600'} px-2 py-1 rounded-full`}>
          {tasks.length}
        </span>
      </div>

      <SortableContext
        id={id}
        items={tasks.map((t) => t.id)}
        strategy={verticalListSortingStrategy}
      >
        <div
          ref={setNodeRef}
          className={`space-y-2 min-h-[100px] transition-colors ${bg || ''} ${
            isOver ? 'ring-2 ring-blue-400 ring-inset rounded-lg' : ''
          }`}
        >
          {children}
        </div>
      </SortableContext>
    </div>
  );
};