const API_BASE_URL = 'http://localhost:8080/api';

const api = {
    // 用户认证
    async register(userData) {
        const response = await fetch(`${API_BASE_URL}/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData),
        });
        return response.json();
    },

    async login(credentials) {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(credentials),
        });
        return response.json();
    },

    // 文章
    async getPosts() {
        const response = await fetch(`${API_BASE_URL}/posts`);
        return response.json();
    },

    async getPost(slug) {
        const response = await fetch(`${API_BASE_URL}/posts/${slug}`);
        return response.json();
    },

    async createPost(postData) {
        const response = await fetch(`${API_BASE_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify(postData),
        });
        return response.json();
    },

    async updatePost(slug, postData) {
        const response = await fetch(`${API_BASE_URL}/posts/${slug}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify(postData),
        });
        return response.json();
    },

    async deletePost(slug) {
        const response = await fetch(`${API_BASE_URL}/posts/${slug}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });
        return response.json();
    },

    // 分类
    async getCategories() {
        const response = await fetch(`${API_BASE_URL}/categories`);
        return response.json();
    },

    async getCategory(id) {
        const response = await fetch(`${API_BASE_URL}/categories/${id}`);
        return response.json();
    },

    async createCategory(categoryData) {
        const response = await fetch(`${API_BASE_URL}/categories`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify(categoryData),
        });
        return response.json();
    },

    async updateCategory(id, categoryData) {
        const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify(categoryData),
        });
        return response.json();
    },

    async deleteCategory(id) {
        const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
        });
        return response.json();
    },
}; 