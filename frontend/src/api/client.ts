import axios, { AxiosInstance } from 'axios';

// Configure base URL - update this to your backend URL
const API_BASE_URL = 'https://matiks-leaderboard-backend-gdj4.onrender.com/api';

const client: AxiosInstance = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add response interceptor for error handling
client.interceptors.response.use(
    (response) => response,
    (error) => {
        console.error('API Error:', error.response?.data || error.message);
        return Promise.reject(error);
    }
);

export default client;
