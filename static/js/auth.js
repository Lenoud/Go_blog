// 检查用户是否已登录
function checkAuth() {
    const token = localStorage.getItem('token');
    const user = JSON.parse(localStorage.getItem('user'));
    
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const dashboardBtn = document.getElementById('dashboardBtn');
    const logoutBtn = document.getElementById('logoutBtn');
    const userName = document.querySelector('.user-name');

    if (token && user) {
        // 用户已登录
        loginBtn.classList.add('hidden');
        registerBtn.classList.add('hidden');
        dashboardBtn.classList.remove('hidden');
        logoutBtn.classList.remove('hidden');
        userName.textContent = user.username;
    } else {
        // 用户未登录
        loginBtn.classList.remove('hidden');
        registerBtn.classList.remove('hidden');
        dashboardBtn.classList.add('hidden');
        logoutBtn.classList.add('hidden');
        userName.textContent = '未登录';
    }
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// 检查是否需要登录
function requireAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// 页面加载时检查认证状态
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();

    // 绑定退出按钮事件
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            logout();
        });
    }
}); 