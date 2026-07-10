import apiClient from '../utils/apiClient.js';

export const authAPI = {
  register: (data) => apiClient.post('/auth/register', data),
  login: (data) => apiClient.post('/auth/login', data),
  refresh: (data) => apiClient.post('/auth/refresh', data),
  logout: (data) => apiClient.post('/auth/logout', data),
};