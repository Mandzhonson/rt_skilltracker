import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext.jsx';
import { ProtectedRoute } from './components/ProtectedRoute.jsx';
import { Login } from './pages/Login.jsx';
import { Register } from './pages/Register.jsx';
import { Profile } from './pages/Profile.jsx';
import { Unauthorized } from './pages/Unauthorized.jsx';
import { AdminUsers } from './pages/admin/AdminUsers.jsx';
import { AdminUserDetail } from './pages/admin/AdminUserDetail.jsx';
import { AdminManagerEmployees } from './pages/admin/AdminManagerEmployees.jsx';
import { ManagerPlans } from './pages/manager/ManagerPlans.jsx';
import { ManagerPlanDetail } from './pages/manager/ManagerPlanDetail.jsx';
import { CreatePlan } from './pages/manager/CreatePlan.jsx';
import { CreateAIPlan } from './pages/manager/CreateAIPlan.jsx';
import { EmployeePlans } from './pages/employee/EmployeePlans.jsx';
import { EmployeePlanDetail } from './pages/employee/EmployeePlanDetail.jsx';
import { EmployeeTasks } from './pages/employee/EmployeeTasks.jsx';
import { EmployeeSkills } from './pages/employee/EmployeeSkills.jsx';
import { EmployeeTest } from './pages/employee/EmployeeTest.jsx';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/unauthorized" element={<Unauthorized />} />
          
          <Route path="/profile" element={
            <ProtectedRoute allowedRoles={['admin', 'manager', 'employee']}>
              <Profile />
            </ProtectedRoute>
          } />
          
          {/* Админ */}
          <Route path="/admin/users" element={
            <ProtectedRoute allowedRoles={['admin']}>
              <AdminUsers />
            </ProtectedRoute>
          } />
          
          <Route path="/admin/users/:id" element={
            <ProtectedRoute allowedRoles={['admin']}>
              <AdminUserDetail />
            </ProtectedRoute>
          } />
          
          <Route path="/admin/managers/:managerId/employees" element={
            <ProtectedRoute allowedRoles={['admin']}>
              <AdminManagerEmployees />
            </ProtectedRoute>
          } />
          
          {/* Менеджер */}
          <Route path="/manager/plans/create" element={
            <ProtectedRoute allowedRoles={['manager', 'admin']}>
              <CreatePlan />
            </ProtectedRoute>
          } />
          
          <Route path="/manager/plans/create-ai" element={
            <ProtectedRoute allowedRoles={['manager', 'admin']}>
              <CreateAIPlan />
            </ProtectedRoute>
          } />
          
          <Route path="/manager/plans/:planId" element={
            <ProtectedRoute allowedRoles={['manager', 'admin']}>
              <ManagerPlanDetail />
            </ProtectedRoute>
          } />
          
          <Route path="/manager/plans" element={
            <ProtectedRoute allowedRoles={['manager', 'admin']}>
              <ManagerPlans />
            </ProtectedRoute>
          } />
          
          {/* Сотрудник */}
          <Route path="/employee/plans/:planId" element={
            <ProtectedRoute allowedRoles={['employee', 'manager', 'admin']}>
              <EmployeePlanDetail />
            </ProtectedRoute>
          } />
          
          <Route path="/employee/plans/:planId/test" element={
            <ProtectedRoute allowedRoles={['employee']}>
              <EmployeeTest />
            </ProtectedRoute>
          } />
          
          <Route path="/employee/plans" element={
            <ProtectedRoute allowedRoles={['employee', 'manager', 'admin']}>
              <EmployeePlans />
            </ProtectedRoute>
          } />
          
          <Route path="/employee/tasks" element={
            <ProtectedRoute allowedRoles={['employee']}>
              <EmployeeTasks />
            </ProtectedRoute>
          } />
          
          <Route path="/employee/skills" element={
            <ProtectedRoute allowedRoles={['employee']}>
              <EmployeeSkills />
            </ProtectedRoute>
          } />
          
          <Route path="/" element={<Login />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;