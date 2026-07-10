import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext.jsx';
import { ProtectedRoute } from './components/ProtectedRoute.jsx';
import { Login } from './pages/Login.jsx';
import { Register } from './pages/Register.jsx';
import { Profile } from './pages/Profile.jsx';
import { Unauthorized } from './pages/Unauthorized.jsx';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/unauthorized" element={<Unauthorized />} />
          
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/admin/*"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <div className="p-8">Админ панель (скоро)</div>
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/manager/*"
            element={
              <ProtectedRoute allowedRoles={['manager']}>
                <div className="p-8">Панель менеджера (скоро)</div>
              </ProtectedRoute>
            }
          />
          
          <Route
            path="/employee/*"
            element={
              <ProtectedRoute allowedRoles={['employee']}>
                <div className="p-8">Панель сотрудника (скоро)</div>
              </ProtectedRoute>
            }
          />
          
          <Route path="/" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;