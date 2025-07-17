document.addEventListener('DOMContentLoaded', async () => {
    const postsGrid = document.getElementById('postsGrid');
    const pagination = document.getElementById('pagination');

    try {
        // 获取文章列表
        const posts = await api.getPosts();
        
        if (posts.error) {
            postsGrid.innerHTML = `<p class="error">${posts.error}</p>`;
            return;
        }

        // 渲染文章列表
        postsGrid.innerHTML = posts.map(post => `
            <article class="post-card">
                <div class="post-content">
                    <h2 class="post-title">
                        <a href="/post.html?slug=${post.slug}">${post.title}</a>
                    </h2>
                    <p class="post-summary">${post.summary || post.content.substring(0, 150) + '...'}</p>
                    <div class="post-meta">
                        <span>作者：${post.author.username}</span>
                        <span>分类：${post.category.name}</span>
                        <span>${new Date(post.created_at).toLocaleDateString()}</span>
                    </div>
                </div>
            </article>
        `).join('');

        // 如果有分页数据，渲染分页控件
        if (posts.pagination) {
            const { current_page, total_pages } = posts.pagination;
            let paginationHTML = '';

            // 上一页
            if (current_page > 1) {
                paginationHTML += `
                    <a href="?page=${current_page - 1}" class="page-item">
                        <i class="fas fa-chevron-left"></i>
                    </a>
                `;
            }

            // 页码
            for (let i = 1; i <= total_pages; i++) {
                paginationHTML += `
                    <a href="?page=${i}" class="page-item ${i === current_page ? 'active' : ''}">
                        ${i}
                    </a>
                `;
            }

            // 下一页
            if (current_page < total_pages) {
                paginationHTML += `
                    <a href="?page=${current_page + 1}" class="page-item">
                        <i class="fas fa-chevron-right"></i>
                    </a>
                `;
            }

            pagination.innerHTML = paginationHTML;
        }
    } catch (error) {
        postsGrid.innerHTML = '<p class="error">加载文章失败，请稍后重试</p>';
        console.error('Load posts error:', error);
    }
}); 