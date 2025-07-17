document.addEventListener('DOMContentLoaded', () => {
    // 检查用户是否已登录
    requireAuth();

    // 加载分类列表
    loadCategories();

    // 监听表单提交
    const form = document.getElementById('createPostForm');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            title: document.getElementById('title').value,
            slug: document.getElementById('slug').value,
            category_id: parseInt(document.getElementById('category').value),
            summary: document.getElementById('summary').value,
            content: document.getElementById('content').value,
            status: document.getElementById('status').value
        };

        try {
            const response = await api.createPost(formData);
            if (response.error) {
                alert(response.error);
                return;
            }
            alert('文章创建成功！');
            window.location.href = '/dashboard.html';
        } catch (error) {
            alert('创建文章失败，请稍后重试');
            console.error('Create post error:', error);
        }
    });

    // 自动生成别名
    const titleInput = document.getElementById('title');
    const slugInput = document.getElementById('slug');
    titleInput.addEventListener('input', () => {
        const slug = titleInput.value
            .toLowerCase()
            .replace(/[^a-z0-9\u4e00-\u9fa5]/g, '-')
            .replace(/-+/g, '-')
            .replace(/^-|-$/g, '');
        slugInput.value = slug;
    });
});

// 加载分类列表
async function loadCategories() {
    try {
        const categories = await api.getCategories();
        const categorySelect = document.getElementById('category');
        
        categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category.id;
            option.textContent = category.name;
            categorySelect.appendChild(option);
        });
    } catch (error) {
        console.error('Load categories error:', error);
        alert('加载分类失败，请刷新页面重试');
    }
} 