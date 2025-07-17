document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const slug = urlParams.get('slug');
    
    if (!slug) {
        window.location.href = '/';
        return;
    }

    try {
        const post = await api.getPost(slug);
        
        if (post.error) {
            document.getElementById('postDetail').innerHTML = `<p class="error">${post.error}</p>`;
            return;
        }

        // 更新页面标题
        document.title = `${post.title} - 博客系统`;

        // 更新文章内容
        document.querySelector('.post-title').textContent = post.title;
        document.querySelector('.post-author').textContent = `作者：${post.author.username}`;
        document.querySelector('.post-category').textContent = `分类：${post.category.name}`;
        document.querySelector('.post-date').textContent = new Date(post.created_at).toLocaleDateString();
        document.querySelector('.post-content').innerHTML = post.content;

        // 检查是否是作者，显示编辑和删除按钮
        const user = JSON.parse(localStorage.getItem('user'));
        if (user && user.id === post.author.id) {
            const postActions = document.getElementById('postActions');
            postActions.style.display = 'flex';
            
            // 编辑按钮
            document.getElementById('editPostBtn').href = `/edit-post.html?slug=${post.slug}`;
            
            // 删除按钮
            document.getElementById('deletePostBtn').addEventListener('click', async () => {
                if (confirm('确定要删除这篇文章吗？')) {
                    try {
                        const response = await api.deletePost(post.slug);
                        if (response.error) {
                            alert(response.error);
                            return;
                        }
                        window.location.href = '/';
                    } catch (error) {
                        alert('删除失败，请稍后重试');
                        console.error('Delete post error:', error);
                    }
                }
            });
        }
    } catch (error) {
        document.getElementById('postDetail').innerHTML = '<p class="error">加载文章失败，请稍后重试</p>';
        console.error('Load post error:', error);
    }
}); 