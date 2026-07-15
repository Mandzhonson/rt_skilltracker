import apiClient from '../utils/apiClient.js';

export const adminAPI = {
  listUsers: (params) => apiClient.get('/admin/users', { params }),
  getUser: (id) => apiClient.get(`/admin/users/${id}`),
  getUserAvatar: (id) => apiClient.get(`/admin/users/${id}/avatar`, { responseType: 'blob' }),
  updateRole: (id, data) => apiClient.patch(`/admin/users/${id}/role`, data),
  updatePosition: (id, data) => apiClient.patch(`/admin/users/${id}/position`, data),
  assignManager: (id, data) => apiClient.patch(`/admin/users/${id}/manager`, data),
  removeManager: (id) => apiClient.delete(`/admin/users/${id}/manager`),
  listEmployeesByManager: (managerId) => apiClient.get(`/admin/managers/${managerId}/employees`),
  listManagers: () => apiClient.get('/admin/users', { params: { role: 'manager', limit: 100 } }),
};