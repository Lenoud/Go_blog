// API 对象
const api = {
    async getPost(id) {
        const response = await fetch(`/api/posts/${id}`);
        if (!response.ok) {
            throw new Error('获取文章失败');
        }
        return await response.json();
    },

    async updatePost(id, data) {
        const response = await fetch(`/api/posts/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        if (!response.ok) {
            throw new Error('更新文章失败');
        }
        return await response.json();
    }
};

// 检查用户是否已登录
function requireAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// 初始化 Markdown 编辑器
let easyMDE;

document.addEventListener('DOMContentLoaded', async () => {
    // 检查用户是否已登录
    if (!requireAuth()) return;

    // 获取文章 ID
    const postId = window.location.pathname.split('/').pop();
    if (!postId) {
        alert('文章不存在');
        window.location.href = '/';
        return;
    }

    // 初始化 Markdown 编辑器
    easyMDE = new EasyMDE({
        element: document.getElementById('content'),
        spellChecker: false,
        status: false,
        autofocus: true,
        toolbar: [
            'bold', 'italic', 'heading', '|',
            'quote', 'unordered-list', 'ordered-list', '|',
            'link', 'image', '|',
            'preview', 'side-by-side', 'fullscreen', '|',
            'guide'
        ],
        placeholder: '使用 Markdown 编写文章...',
        previewRender: (text) => {
            return DOMPurify.sanitize(marked.parse(text));
        }
    });

    // 加载文章数据
    await loadPost(postId);

    // 处理表单提交
    const form = document.getElementById('editPostForm');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const post = {
            title: document.getElementById('title').value,
            content: easyMDE.value(),
            tags: document.getElementById('tags').value.split(',').map(tag => tag.trim()).filter(tag => tag)
        };
        
        try {
            const response = await api.updatePost(postId, post);
            alert('保存成功！');
            window.location.href = '/';
        } catch (error) {
            console.error('Error:', error);
            alert('保存失败，请重试');
        }
    });
});

// 加载文章数据
async function loadPost(postId) {
    try {
        const post = await api.getPost(postId);
        
        document.getElementById('title').value = post.title;
        easyMDE.value(post.content);
        document.getElementById('tags').value = post.tags.join(', ');
    } catch (error) {
        console.error('Error:', error);
        alert('加载文章失败，请刷新页面重试');
    }
}

// 导入 Markdown 文件
window.importMarkdown = function() {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.md,.markdown';
    
    input.onchange = e => {
        const file = e.target.files[0];
        const reader = new FileReader();
        
        reader.onload = function(event) {
            const content = event.target.result;
            easyMDE.value(content);
        };
        
        reader.readAsText(file);
    };
    
    input.click();
};

// 导出 Markdown 文件
window.exportMarkdown = function() {
    const content = easyMDE.value();
    const title = document.getElementById('title').value || 'untitled';
    const blob = new Blob([content], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${title}.md`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}; 