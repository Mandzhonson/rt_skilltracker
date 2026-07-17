import apiClient from '../utils/apiClient.js';

export const userAPI = {
  getProfile: () => apiClient.get('/users/me'),
  updateProfile: (data) => apiClient.patch('/users/me', data),
  updatePassword: (data) => apiClient.patch('/users/me/password', data),
  setAvatar: (data) => apiClient.put('/users/me/avatar', data, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }),
  getAvatar: () => apiClient.get('/users/me/avatar', { responseType: 'blob' }),
  deleteAvatar: () => apiClient.delete('/users/me/avatar'),
};