document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');

    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        // 验证密码
        if (password !== confirmPassword) {
            alert('两次输入的密码不一致');
            return;
        }

        try {
            const response = await api.register({ username, email, password });
            
            if (response.error) {
                alert(response.error);
                return;
            }

            alert('注册成功，请登录');
            window.location.href = '/login.html';
        } catch (error) {
            alert('注册失败，请稍后重试');
            console.error('Register error:', error);
        }
    });
}); 