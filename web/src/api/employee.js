import apiClient from '../utils/apiClient.js';

export const employeeAPI = {
  listPlans: () => apiClient.get('/employee/plans'),
  getPlan: (planId) => apiClient.get(`/employee/plans/${planId}`),
  updateTaskStatus: (taskId, status) => apiClient.patch(`/employee/tasks/${taskId}/status`, { status }),
};