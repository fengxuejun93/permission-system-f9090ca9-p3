const API_BASE = '/api';

const state = {
    items: [],
    currentItem: null,
    editingId: null,
    meta: null,
    filters: {
        keyword: '',
        category: '',
        city: '',
        status: ''
    }
};

const $ = (selector) => document.querySelector(selector);
const $$ = (selector) => document.querySelectorAll(selector);

async function request(url, options = {}) {
    const response = await fetch(url, {
        headers: {
            'Content-Type': 'application/json',
            ...options.headers
        },
        ...options
    });
    const data = await response.json();
    if (data.code >= 400) {
        const error = new Error(data.message);
        error.code = data.code;
        throw error;
    }
    return data.data;
}

function showToast(message, type = 'info') {
    const toast = $('#toast');
    toast.textContent = message;
    toast.className = `toast show ${type}`;
    setTimeout(() => {
        toast.className = `toast ${type}`;
    }, 3000);
}

function getStatusLabel(status) {
    const map = {
        'active': '上架中',
        'offline': '已下架',
        'traded': '已置换',
        'pending': '待审核'
    };
    return map[status] || status;
}

function getStatusClass(status) {
    return `status-${status}`;
}

async function loadMeta() {
    state.meta = await request(`${API_BASE}/meta`);

    const categoryFilter = $('#categoryFilter');
    state.meta.categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat;
        option.textContent = cat;
        categoryFilter.appendChild(option);
    });

    const cityFilter = $('#cityFilter');
    state.meta.cities.forEach(city => {
        const option = document.createElement('option');
        option.value = city;
        option.textContent = city;
        cityFilter.appendChild(option);
    });

    const statusFilter = $('#statusFilter');
    state.meta.statuses.forEach(s => {
        const option = document.createElement('option');
        option.value = s.value;
        option.textContent = s.label;
        statusFilter.appendChild(option);
    });

    const formCategory = $('#formCategory');
    state.meta.categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat;
        option.textContent = cat;
        formCategory.appendChild(option);
    });

    const formCity = $('#formCity');
    state.meta.cities.forEach(city => {
        const option = document.createElement('option');
        option.value = city;
        option.textContent = city;
        formCity.appendChild(option);
    });

    const formCondition = $('#formCondition');
    state.meta.conditions.forEach(c => {
        const option = document.createElement('option');
        option.value = c;
        option.textContent = c;
        formCondition.appendChild(option);
    });

    const formDesiredCategory = $('#formDesiredCategory');
    state.meta.categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat;
        option.textContent = cat;
        formDesiredCategory.appendChild(option);
    });
}

async function loadStatistics() {
    const stats = await request(`${API_BASE}/statistics`);
    $('#statTotal').textContent = stats.totalItems;
    $('#statActive').textContent = stats.activeItems;
    $('#statOffline').textContent = stats.offlineItems;
    $('#statTraded').textContent = stats.tradedItems;
    $('#statViews').textContent = stats.totalViews.toLocaleString();
    $('#statFavorites').textContent = stats.totalFavorites.toLocaleString();
    $('#statIntents').textContent = stats.totalIntents.toLocaleString();
}

async function loadItems() {
    const params = new URLSearchParams();
    Object.entries(state.filters).forEach(([key, value]) => {
        if (value) params.append(key, value);
    });

    state.items = await request(`${API_BASE}/items?${params.toString()}`);
    renderItemList();
}

function renderItemList() {
    const listContainer = $('#itemList');

    if (state.items.length === 0) {
        listContainer.innerHTML = `
            <div style="padding: 40px; text-align: center; color: #909399;">
                <div style="font-size: 48px; margin-bottom: 12px;">📭</div>
                <div>暂无符合条件的货品</div>
            </div>
        `;
        return;
    }

    listContainer.innerHTML = state.items.map(item => `
        <div class="item-card ${state.currentItem?.id === item.id ? 'active' : ''}" data-id="${item.id}">
            <div class="item-card-header">
                <div class="item-card-title">${escapeHtml(item.title)}</div>
                <span class="status-badge ${getStatusClass(item.status)}">${getStatusLabel(item.status)}</span>
            </div>
            <div class="item-card-meta">
                <span class="meta-tag">${escapeHtml(item.category)}</span>
                <span class="meta-tag">${escapeHtml(item.city)}</span>
                <span class="meta-tag">${escapeHtml(item.condition)}</span>
            </div>
            <div class="item-card-stats">
                <span class="stat-item">👁️ ${item.viewCount}</span>
                <span class="stat-item">⭐ ${item.favoriteCount}</span>
                <span class="stat-item">💬 ${item.tradeIntentCount}</span>
            </div>
        </div>
    `).join('');

    $$('.item-card').forEach(card => {
        card.addEventListener('click', () => {
            const id = card.dataset.id;
            loadItemDetail(id);
        });
    });
}

async function loadItemDetail(id) {
    state.currentItem = await request(`${API_BASE}/items/${id}`);
    renderItemDetail();
    renderItemList();
    await loadStatistics();
}

function renderItemDetail() {
    const item = state.currentItem;
    if (!item) return;

    const isOffline = item.status === 'offline';
    const isTraded = item.status === 'traded';

    $('#detailContent').innerHTML = `
        <div class="detail-header">
            <div class="detail-title">${escapeHtml(item.title)}</div>
            <div class="detail-meta-row">
                <span class="detail-tag highlight">${escapeHtml(item.category)}</span>
                <span class="detail-tag">${escapeHtml(item.city)}</span>
                <span class="detail-tag">${escapeHtml(item.condition)}</span>
                <span class="detail-tag">期望置换: ${escapeHtml(item.desiredCategory)}</span>
                <span class="status-badge ${getStatusClass(item.status)}">${getStatusLabel(item.status)}</span>
            </div>
            <div class="detail-stats-row">
                <div class="detail-stat">
                    <div class="detail-stat-value">${item.viewCount}</div>
                    <div class="detail-stat-label">浏览</div>
                </div>
                <div class="detail-stat">
                    <div class="detail-stat-value">${item.favoriteCount}</div>
                    <div class="detail-stat-label">收藏</div>
                </div>
                <div class="detail-stat">
                    <div class="detail-stat-value">${item.tradeIntentCount}</div>
                    <div class="detail-stat-label">置换意向</div>
                </div>
            </div>
        </div>

        <div class="detail-section-block">
            <div class="detail-section-title">基本信息</div>
            <div class="detail-info-grid">
                <div class="detail-info-item">
                    <span class="detail-info-label">发布者</span>
                    <span class="detail-info-value">${escapeHtml(item.publisher)}</span>
                </div>
                <div class="detail-info-item">
                    <span class="detail-info-label">货品状态</span>
                    <span class="detail-info-value">${getStatusLabel(item.status)}</span>
                </div>
                <div class="detail-info-item">
                    <span class="detail-info-label">已收藏</span>
                    <span class="detail-info-value">${item.isFavorited ? '✅ 是' : '❌ 否'}</span>
                </div>
                <div class="detail-info-item">
                    <span class="detail-info-label">已沟通</span>
                    <span class="detail-info-value">${item.hasCommunicated ? '✅ 是' : '❌ 否'}</span>
                </div>
                <div class="detail-info-item">
                    <span class="detail-info-label">发布时间</span>
                    <span class="detail-info-value">${formatDate(item.createdAt)}</span>
                </div>
                <div class="detail-info-item">
                    <span class="detail-info-label">更新时间</span>
                    <span class="detail-info-value">${formatDate(item.updatedAt)}</span>
                </div>
            </div>
        </div>

        <div class="detail-section-block">
            <div class="detail-section-title">详细描述</div>
            <div class="detail-description">${escapeHtml(item.description) || '暂无详细描述'}</div>
        </div>

        <div class="detail-actions">
            <button class="action-btn action-btn-favorite ${item.isFavorited ? 'active' : ''}" id="btnFavorite">
                ${item.isFavorited ? '★ 已收藏' : '☆ 收藏'}
            </button>
            <button class="action-btn action-btn-intent" id="btnIntent" ${isOffline ? 'disabled' : ''}>
                💬 发起置换意向
            </button>
            <button class="action-btn action-btn-communicated ${item.hasCommunicated ? 'active' : ''}" id="btnCommunicated">
                ${item.hasCommunicated ? '✓ 已沟通' : '📞 标记已沟通'}
            </button>
            ${!isOffline && !isTraded ? `
                <button class="action-btn action-btn-offline" id="btnOffline">⬇️ 下架</button>
            ` : ''}
            ${isOffline ? `
                <button class="action-btn action-btn-relist" id="btnRelist">⬆️ 重新上架</button>
            ` : ''}
            <button class="action-btn action-btn-edit" id="btnEdit">✏️ 编辑</button>
        </div>
    `;

    $('#btnFavorite').addEventListener('click', () => handleToggleFavorite(item.id));
    $('#btnIntent').addEventListener('click', () => handleAddTradeIntent(item.id));
    $('#btnCommunicated').addEventListener('click', () => handleMarkCommunicated(item.id));
    const btnOffline = $('#btnOffline');
    if (btnOffline) btnOffline.addEventListener('click', () => handleOffline(item.id));
    const btnRelist = $('#btnRelist');
    if (btnRelist) btnRelist.addEventListener('click', () => handleRelist(item.id));
    $('#btnEdit').addEventListener('click', () => openEditModal(item));
}

async function handleToggleFavorite(id) {
    try {
        const updated = await request(`${API_BASE}/items/${id}/favorite`, { method: 'POST' });
        state.currentItem = updated;
        updateItemInList(updated);
        renderItemDetail();
        renderItemList();
        await loadStatistics();
        showToast(updated.isFavorited ? '已收藏' : '已取消收藏', 'success');
    } catch (error) {
        showToast(`收藏操作失败：${error.message}`, 'error');
    }
}

async function handleAddTradeIntent(id) {
    try {
        const updated = await request(`${API_BASE}/items/${id}/trade-intent`, { method: 'POST' });
        state.currentItem = updated;
        updateItemInList(updated);
        renderItemDetail();
        renderItemList();
        await loadStatistics();
        showToast('置换意向已发出，请耐心等待发布者回复', 'success');
    } catch (error) {
        if (error.message.includes('已下架')) {
            showToast('⚠️ 该货品已下架，无法发起置换意向', 'warning');
        } else {
            showToast(`发起失败：${error.message}`, 'error');
        }
    }
}

async function handleMarkCommunicated(id) {
    try {
        const updated = await request(`${API_BASE}/items/${id}/mark-communicated`, { method: 'POST' });
        state.currentItem = updated;
        updateItemInList(updated);
        renderItemDetail();
        renderItemList();
        showToast('已标记为已沟通', 'success');
    } catch (error) {
        showToast(`标记失败：${error.message}`, 'error');
    }
}

async function handleOffline(id) {
    if (!confirm('确定要下架该货品吗？下架后其他用户将无法发起置换意向。')) return;
    try {
        const updated = await request(`${API_BASE}/items/${id}/offline`, { method: 'POST' });
        state.currentItem = updated;
        updateItemInList(updated);
        renderItemDetail();
        renderItemList();
        await loadStatistics();
        showToast('货品已下架', 'success');
    } catch (error) {
        showToast(`下架失败：${error.message}`, 'error');
    }
}

async function handleRelist(id) {
    try {
        const updated = await request(`${API_BASE}/items/${id}/relist`, { method: 'POST' });
        state.currentItem = updated;
        updateItemInList(updated);
        renderItemDetail();
        renderItemList();
        await loadStatistics();
        showToast('货品已重新上架', 'success');
    } catch (error) {
        showToast(`重新上架失败：${error.message}`, 'error');
    }
}

function updateItemInList(updatedItem) {
    const index = state.items.findIndex(i => i.id === updatedItem.id);
    if (index !== -1) {
        state.items[index] = updatedItem;
    }
}

function openPublishModal() {
    state.editingId = null;
    $('#modalTitle').textContent = '发布货品';
    $('#itemForm').reset();
    $('#formPublisher').disabled = false;
    $('#formModal').classList.add('show');
}

function openEditModal(item) {
    state.editingId = item.id;
    $('#modalTitle').textContent = '编辑货品';
    $('#formTitle').value = item.title;
    $('#formCategory').value = item.category;
    $('#formCity').value = item.city;
    $('#formCondition').value = item.condition;
    $('#formPublisher').value = item.publisher;
    $('#formPublisher').disabled = true;
    $('#formDesiredCategory').value = item.desiredCategory;
    $('#formDescription').value = item.description || '';
    $('#formModal').classList.add('show');
}

function closeModal() {
    $('#formModal').classList.remove('show');
    state.editingId = null;
}

async function handleFormSubmit(e) {
    e.preventDefault();

    const data = {
        title: $('#formTitle').value.trim(),
        category: $('#formCategory').value,
        city: $('#formCity').value,
        condition: $('#formCondition').value,
        desiredCategory: $('#formDesiredCategory').value,
        description: $('#formDescription').value.trim()
    };

    try {
        if (state.editingId) {
            const updated = await request(`${API_BASE}/items/${state.editingId}`, {
                method: 'PUT',
                body: JSON.stringify(data)
            });
            state.currentItem = updated;
            updateItemInList(updated);
            showToast('货品更新成功', 'success');
        } else {
            data.publisher = $('#formPublisher').value.trim();
            const created = await request(`${API_BASE}/items`, {
                method: 'POST',
                body: JSON.stringify(data)
            });
            state.currentItem = created;
            showToast('货品发布成功', 'success');
        }

        closeModal();
        await loadItems();
        renderItemDetail();
        await loadStatistics();

        if (state.currentItem) {
            const card = document.querySelector(`.item-card[data-id="${state.currentItem.id}"]`);
            if (card) {
                card.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    } catch (error) {
        showToast(`提交失败：${error.message}`, 'error');
    }
}

function handleSearch() {
    state.filters.keyword = $('#searchInput').value.trim();
    loadItems();
}

function handleFilterChange() {
    state.filters.category = $('#categoryFilter').value;
    state.filters.city = $('#cityFilter').value;
    state.filters.status = $('#statusFilter').value;
    loadItems();
}

function handleResetFilters() {
    $('#searchInput').value = '';
    $('#categoryFilter').value = '';
    $('#cityFilter').value = '';
    $('#statusFilter').value = '';
    state.filters = {
        keyword: '',
        category: '',
        city: '',
        status: ''
    };
    loadItems();
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function initEventListeners() {
    $('#publishBtn').addEventListener('click', openPublishModal);
    $('#modalClose').addEventListener('click', closeModal);
    $('#cancelBtn').addEventListener('click', closeModal);
    $('#itemForm').addEventListener('submit', handleFormSubmit);
    $('#searchBtn').addEventListener('click', handleSearch);
    $('#searchInput').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') handleSearch();
    });
    $('#categoryFilter').addEventListener('change', handleFilterChange);
    $('#cityFilter').addEventListener('change', handleFilterChange);
    $('#statusFilter').addEventListener('change', handleFilterChange);
    $('#resetFilterBtn').addEventListener('click', handleResetFilters);

    $('#formModal').addEventListener('click', (e) => {
        if (e.target.id === 'formModal') {
            closeModal();
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && $('#formModal').classList.contains('show')) {
            closeModal();
        }
    });
}

async function init() {
    initEventListeners();
    await loadMeta();
    await loadStatistics();
    await loadItems();

    if (state.items.length > 0) {
        await loadItemDetail(state.items[0].id);
    }
}

document.addEventListener('DOMContentLoaded', init);
