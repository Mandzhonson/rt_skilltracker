import { useEffect } from 'react';

export const SkillsGenerationModal = ({ isOpen, onClose, message, isComplete }) => {
  useEffect(() => {
    if (isComplete) {
      const timer = setTimeout(() => {
        onClose();
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [isComplete, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full mx-4 p-8 transform transition-all">
        <div className="flex flex-col items-center text-center">
          {isComplete ? (
            <div className="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center mb-4">
              <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
          ) : (
            <div className="w-16 h-16 rounded-full bg-blue-100 flex items-center justify-center mb-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          )}
          
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            {isComplete ? 'Навыки сформированы!' : 'Формирование навыков'}
          </h3>
          
          <p className="text-gray-600 text-sm leading-relaxed">
            {message}
          </p>
          
          {!isComplete && (
            <p className="text-xs text-gray-400 mt-3">
              Вы можете продолжить работу, генерация будет завершена в фоне
            </p>
          )}
          
          {isComplete && (
            <button
              onClick={onClose}
              className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Понятно
            </button>
          )}
        </div>
      </div>
    </div>
  );
};