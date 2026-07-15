import apiClient from '../utils/apiClient.js';

export const skillsAPI = {
  getMySkills: () => apiClient.get('/employee/skills'),
  getEmployeeSkills: (employeeId) => apiClient.get(`/employee/${employeeId}/skills`),
};