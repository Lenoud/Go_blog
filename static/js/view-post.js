// 获取文章 ID
const postId = window.location.pathname.split('/').pop();

// 加载文章数据
async function loadPost() {
    try {
        const response = await fetch(`/api/posts/${postId}`);
        if (response.ok) {
            const post = await response.json();
            
            // 设置标题
            document.getElementById('post-title').textContent = post.title;
            
            // 设置日期
            const date = new Date(post.created_at);
            document.getElementById('post-date').textContent = `发布于 ${date.toLocaleString()}`;
            
            // 设置标签
            const tagsContainer = document.getElementById('post-tags');
            post.tags.forEach(tag => {
                const tagSpan = document.createElement('span');
                tagSpan.className = 'badge bg-secondary me-1';
                tagSpan.textContent = tag;
                tagsContainer.appendChild(tagSpan);
            });
            
            // 渲染 Markdown 内容
            document.getElementById('post-content').innerHTML = marked.parse(post.content);
            
            // 设置编辑和删除按钮
            document.getElementById('editBtn').onclick = () => {
                window.location.href = `/edit/${postId}`;
            };
            
            document.getElementById('deleteBtn').onclick = async () => {
                if (confirm('确定要删除这篇文章吗？')) {
                    try {
                        const response = await fetch(`/api/posts/${postId}`, {
                            method: 'DELETE'
                        });
                        
                        if (response.ok) {
                            window.location.href = '/';
                        } else {
                            alert('删除失败，请重试');
                        }
                    } catch (error) {
                        console.error('Error:', error);
                        alert('删除失败，请重试');
                    }
                }
            };
        } else {
            alert('加载文章失败');
            window.location.href = '/';
        }
    } catch (error) {
        console.error('Error:', error);
        alert('加载文章失败');
        window.location.href = '/';
    }
}

// 页面加载时获取文章数据
loadPost(); 