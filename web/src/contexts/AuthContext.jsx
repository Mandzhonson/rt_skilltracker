import { createContext, useContext, useState, useEffect } from 'react';
import { authAPI } from '../api/auth.js';
import { userAPI } from '../api/user.js';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const loadUser = async () => {
      const token = localStorage.getItem('accessToken');
      if (token) {
        try {
          const response = await userAPI.getProfile();
          setUser(response.data);
          setIsAuthenticated(true);
        } catch (error) {
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
        }
      }
      setLoading(false);
    };
    loadUser();
  }, []);

  const login = async (email, password) => {
    try {
      const response = await authAPI.login({ email, password });
      const { access_token, refresh_token } = response.data;
      localStorage.setItem('accessToken', access_token);
      localStorage.setItem('refreshToken', refresh_token);
      
      const userResponse = await userAPI.getProfile();
      setUser(userResponse.data);
      setIsAuthenticated(true);
      return response;
    } catch (error) {
      // Пробрасываем ошибку дальше в Login.jsx
      throw error;
    }
  };

  const register = async (data) => {
    const response = await authAPI.register(data);
    return response;
  };

  const logout = async () => {
    try {
      const refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        await authAPI.logout({ refresh_token: refreshToken });
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      setUser(null);
      setIsAuthenticated(false);
    }
  };

  const value = {
    user,
    setUser,
    login,
    register,
    logout,
    isAuthenticated,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};