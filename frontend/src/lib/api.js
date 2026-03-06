import axios from 'axios';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const api = axios.create({
  baseURL: `${BACKEND_URL}/api/v1/`,
});

// Helper function to get user ID from localStorage
export const getUserId = () => {
  const userStr = localStorage.getItem('sheleads_user');
  if (!userStr) return null;
  try {
    const user = JSON.parse(userStr);
    return user._id || user.id;
  } catch (e) {
    return null;
  }
};

// Helper function to add user_id to query params
export const addUserIdToParams = (params = {}) => {
  const userId = getUserId();
  if (userId) {
    return { ...params, user_id: userId };
  }
  return params;
};

// Helper function to add user_id to form data
export const addUserIdToFormData = (formData) => {
  const userId = getUserId();
  if (userId && formData instanceof FormData) {
    formData.append('user_id', userId);
  }
  return formData;
};

export default api;
