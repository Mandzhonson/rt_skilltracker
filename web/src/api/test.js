import apiClient from '../utils/apiClient.js';

export const testAPI = {
  getTest: (planId) => apiClient.get(`/employee/plans/${planId}/test`),
  submitTest: (planId, answers) => apiClient.post(`/employee/plans/${planId}/test`, { answers }),
};