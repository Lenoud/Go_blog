document.addEventListener('DOMContentLoaded', async () => {
    // 检查用户是否已登录
    requireAuth();

    try {
        // 获取用户文章列表
        const posts = await api.getPosts();
        
        if (posts.error) {
            document.getElementById('postsTable').innerHTML = `
                <tr>
                    <td colspan="5" class="text-center">${posts.error}</td>
                </tr>
            `;
            return;
        }

        // 渲染文章列表
        document.getElementById('postsTable').innerHTML = posts.map(post => `
            <tr>
                <td>
                    <a href="/post.html?slug=${post.slug}" class="post-link">${post.title}</a>
                </td>
                <td>${post.category.name}</td>
                <td>
                    <span class="status-badge ${post.status === 'published' ? 'status-published' : 'status-draft'}">
                        ${post.status === 'published' ? '已发布' : '草稿'}
                    </span>
                </td>
                <td>${new Date(post.created_at).toLocaleDateString()}</td>
                <td>
                    <div class="action-buttons">
                        <a href="/edit-post.html?slug=${post.slug}" class="btn btn-sm btn-primary">
                            <i class="fas fa-edit"></i> 编辑
                        </a>
                        <button class="btn btn-sm btn-danger" onclick="deletePost('${post.slug}')">
                            <i class="fas fa-trash"></i> 删除
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');

    } catch (error) {
        document.getElementById('postsTable').innerHTML = `
            <tr>
                <td colspan="5" class="text-center">加载失败，请稍后重试</td>
            </tr>
        `;
        console.error('Load posts error:', error);
    }
});

// 删除文章
async function deletePost(slug) {
    if (!confirm('确定要删除这篇文章吗？')) {
        return;
    }

    try {
        const response = await api.deletePost(slug);
        if (response.error) {
            alert(response.error);
            return;
        }
        // 重新加载页面
        window.location.reload();
    } catch (error) {
        alert('删除失败，请稍后重试');
        console.error('Delete post error:', error);
    }
} 