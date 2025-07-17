document.addEventListener('DOMContentLoaded', async () => {
    // 检查用户是否已登录
    requireAuth();

    const modal = document.getElementById('categoryModal');
    const modalTitle = document.getElementById('modalTitle');
    const categoryForm = document.getElementById('categoryForm');
    const createCategoryBtn = document.getElementById('createCategoryBtn');
    const closeModal = document.getElementById('closeModal');
    const cancelBtn = document.getElementById('cancelBtn');

    let editingCategoryId = null;

    // 加载分类列表
    async function loadCategories() {
        try {
            const categories = await api.getCategories();
            
            if (categories.error) {
                document.getElementById('categoriesGrid').innerHTML = `
                    <p class="error">${categories.error}</p>
                `;
                return;
            }

            document.getElementById('categoriesGrid').innerHTML = categories.map(category => `
                <div class="category-card">
                    <div class="category-content">
                        <h3 class="category-name">${category.name}</h3>
                        <p class="category-slug">${category.slug}</p>
                    </div>
                    <div class="category-actions">
                        <button class="btn btn-sm btn-primary" onclick="editCategory(${category.id}, '${category.name}', '${category.slug}')">
                            <i class="fas fa-edit"></i> 编辑
                        </button>
                        <button class="btn btn-sm btn-danger" onclick="deleteCategory(${category.id})">
                            <i class="fas fa-trash"></i> 删除
                        </button>
                    </div>
                </div>
            `).join('');

        } catch (error) {
            document.getElementById('categoriesGrid').innerHTML = `
                <p class="error">加载分类失败，请稍后重试</p>
            `;
            console.error('Load categories error:', error);
        }
    }

    // 显示模态框
    function showModal(title) {
        modalTitle.textContent = title;
        modal.style.display = 'block';
        categoryForm.reset();
        editingCategoryId = null;
    }

    // 隐藏模态框
    function hideModal() {
        modal.style.display = 'none';
        categoryForm.reset();
        editingCategoryId = null;
    }

    // 新建分类按钮点击事件
    createCategoryBtn.addEventListener('click', () => {
        showModal('新建分类');
    });

    // 关闭模态框按钮点击事件
    closeModal.addEventListener('click', hideModal);
    cancelBtn.addEventListener('click', hideModal);

    // 点击模态框外部关闭
    window.addEventListener('click', (e) => {
        if (e.target === modal) {
            hideModal();
        }
    });

    // 表单提交事件
    categoryForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            name: document.getElementById('categoryName').value,
            slug: document.getElementById('categorySlug').value
        };

        try {
            let response;
            if (editingCategoryId) {
                response = await api.updateCategory(editingCategoryId, formData);
            } else {
                response = await api.createCategory(formData);
            }

            if (response.error) {
                alert(response.error);
                return;
            }

            hideModal();
            loadCategories();
        } catch (error) {
            alert('保存失败，请稍后重试');
            console.error('Save category error:', error);
        }
    });

    // 编辑分类
    window.editCategory = (id, name, slug) => {
        editingCategoryId = id;
        document.getElementById('categoryName').value = name;
        document.getElementById('categorySlug').value = slug;
        showModal('编辑分类');
    };

    // 删除分类
    window.deleteCategory = async (id) => {
        if (!confirm('确定要删除这个分类吗？')) {
            return;
        }

        try {
            const response = await api.deleteCategory(id);
            if (response.error) {
                alert(response.error);
                return;
            }
            loadCategories();
        } catch (error) {
            alert('删除失败，请稍后重试');
            console.error('Delete category error:', error);
        }
    };

    // 初始加载分类列表
    loadCategories();
}); 