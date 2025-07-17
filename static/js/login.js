document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await api.login({ username, password });
            
            if (response.error) {
                alert(response.error);
                return;
            }

            // 保存token和用户信息
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));

            // 跳转到首页
            window.location.href = '/';
        } catch (error) {
            alert('登录失败，请稍后重试');
            console.error('Login error:', error);
        }
    });
}); 