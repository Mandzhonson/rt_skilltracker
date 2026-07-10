import { Link } from 'react-router-dom';

export const Unauthorized = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div className="text-center">
      <div className="text-6xl font-bold text-red-600 mb-4">403</div>
      <h1 className="text-3xl font-bold text-gray-900 mb-2">Доступ запрещен</h1>
      <p className="text-gray-600 mb-6">
        У вас нет прав для доступа к этой странице.
      </p>
      <Link to="/profile" className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700">
        Вернуться в профиль
      </Link>
    </div>
  </div>
);