// Educational Game Database - Admin Interface
class AdminPanel {
  constructor() {
    this.apiBase = '/api';
    this.accounts = [];
    this.stats = null;
    this.init();
  }

  async init() {
    this.setupEventListeners();
    await this.loadData();
    this.renderAll();
  }

  setupEventListeners() {
    // Create account form
    const createForm = document.getElementById('createAccountForm');
    if (createForm) {
      createForm.addEventListener('submit', (e) => this.handleCreateAccount(e));
    }

    // Search functionality
    const searchInput = document.getElementById('searchInput');
    if (searchInput) {
      searchInput.addEventListener('input', (e) => this.handleSearch(e.target.value));
    }

    // Filter functionality
    const filterSelect = document.getElementById('filterSelect');
    if (filterSelect) {
      filterSelect.addEventListener('change', (e) => this.handleFilter(e.target.value));
    }

    // Refresh button
    const refreshBtn = document.getElementById('refreshBtn');
    if (refreshBtn) {
      refreshBtn.addEventListener('click', () => this.refreshData());
    }

    // Modal close buttons
    document.addEventListener('click', (e) => {
      if (e.target.classList.contains('close-btn') || e.target.classList.contains('modal')) {
        this.closeModal();
      }
    });

    // Export button
    const exportBtn = document.getElementById('exportBtn');
    if (exportBtn) {
      exportBtn.addEventListener('click', () => this.exportData());
    }
  }

  // API Methods
  async apiCall(endpoint, method = 'GET', data = null) {
    const config = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
    };

    if (data) {
      config.body = JSON.stringify(data);
    }

    try {
      const response = await fetch(`${this.apiBase}${endpoint}`, config);
      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.error || 'An error occurred');
      }

      return result;
    } catch (error) {
      console.error('API call failed:', error);
      throw error;
    }
  }

  async loadData() {
    try {
      this.showLoading(true);
      const [accounts, stats] = await Promise.all([
        this.apiCall('/accounts'),
        this.apiCall('/stats')
      ]);
      
      this.accounts = accounts || [];
      this.stats = stats;
    } catch (error) {
      this.showMessage('Failed to load data: ' + error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  async refreshData() {
    await this.loadData();
    this.renderAll();
    this.showMessage('Data refreshed successfully', 'success');
  }

  renderAll() {
    this.renderStats();
    this.renderAccountsTable();
  }

  renderStats() {
    if (!this.stats) return;

    this.updateElement('totalAccounts', this.stats.total_accounts);
    this.updateElement('activeAccounts', this.stats.active_accounts);
    this.updateElement('averageLevel', this.stats.average_game_level.toFixed(1));
    this.updateElement('totalExperience', this.stats.total_experience.toLocaleString());
    
    // Calculate inactive accounts
    const inactiveAccounts = this.stats.total_accounts - this.stats.active_accounts;
    this.updateElement('inactiveAccounts', inactiveAccounts);
  }

  renderAccountsTable() {
    const tbody = document.getElementById('accountsTableBody');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (this.accounts.length === 0) {
      tbody.innerHTML = '<tr><td colspan="9" class="text-center">No accounts found</td></tr>';
      return;
    }

    this.accounts.forEach(account => {
      const row = this.createAccountRow(account);
      tbody.appendChild(row);
    });
  }

  createAccountRow(account) {
    const row = document.createElement('tr');
    
    const statusBadge = account.is_active 
      ? '<span class="badge badge-success">Active</span>'
      : '<span class="badge badge-danger">Inactive</span>';

    const createdDate = new Date(account.created_at).toLocaleDateString();
    
    row.innerHTML = `
      <td>${account.id}</td>
      <td>${account.username}</td>
      <td>${account.email}</td>
      <td>${account.first_name} ${account.last_name}</td>
      <td>${account.grade || 'N/A'}</td>
      <td>${account.school || 'N/A'}</td>
      <td>${account.game_level}</td>
      <td>${account.experience}</td>
      <td>${statusBadge}</td>
      <td>${createdDate}</td>
      <td>
        <button class="btn btn-sm btn-primary" onclick="adminPanel.editAccount(${account.id})">Edit</button>
        <button class="btn btn-sm btn-danger" onclick="adminPanel.deleteAccount(${account.id})">Delete</button>
      </td>
    `;
    
    return row;
  }

  async handleCreateAccount(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const accountData = {
      username: formData.get('username'),
      email: formData.get('email'),
      password: formData.get('password'),
      first_name: formData.get('firstName'),
      last_name: formData.get('lastName'),
      grade: parseInt(formData.get('grade')) || 0,
      school: formData.get('school')
    };

    try {
      this.showLoading(true);
      await this.apiCall('/accounts', 'POST', accountData);
      
      this.showMessage('Account created successfully!', 'success');
      this.closeModal();
      e.target.reset();
      await this.refreshData();

    } catch (error) {
      this.showMessage('Failed to create account: ' + error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  async editAccount(id) {
    const account = this.accounts.find(acc => acc.id === id);
    if (!account) return;

    // Populate edit form
    this.populateEditForm(account);
    this.showModal('editAccountModal');
  }

  populateEditForm(account) {
    const form = document.getElementById('editAccountForm');
    if (!form) return;

    form.querySelector('[name="accountId"]').value = account.id;
    form.querySelector('[name="firstName"]').value = account.first_name;
    form.querySelector('[name="lastName"]').value = account.last_name;
    form.querySelector('[name="grade"]').value = account.grade;
    form.querySelector('[name="school"]').value = account.school;
    form.querySelector('[name="gameLevel"]').value = account.game_level;
    form.querySelector('[name="experience"]').value = account.experience;
    form.querySelector('[name="isActive"]').checked = account.is_active;
  }

  async handleEditAccount(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const accountId = formData.get('accountId');
    const updateData = {
      first_name: formData.get('firstName'),
      last_name: formData.get('lastName'),
      grade: parseInt(formData.get('grade')) || 0,
      school: formData.get('school'),
      game_level: parseInt(formData.get('gameLevel')) || 1,
      experience: parseInt(formData.get('experience')) || 0,
      is_active: formData.get('isActive') === 'on'
    };

    try {
      this.showLoading(true);
      await this.apiCall(`/accounts/${accountId}`, 'PUT', updateData);
      
      this.showMessage('Account updated successfully!', 'success');
      this.closeModal();
      await this.refreshData();

    } catch (error) {
      this.showMessage('Failed to update account: ' + error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  async deleteAccount(id) {
    const account = this.accounts.find(acc => acc.id === id);
    if (!account) return;

    const confirmed = confirm(`Are you sure you want to delete account for ${account.username}?`);
    if (!confirmed) return;

    try {
      this.showLoading(true);
      await this.apiCall(`/accounts/${id}`, 'DELETE');
      
      this.showMessage('Account deleted successfully!', 'success');
      await this.refreshData();

    } catch (error) {
      this.showMessage('Failed to delete account: ' + error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  handleSearch(query) {
    const filteredAccounts = this.accounts.filter(account => {
      const searchableText = [
        account.username,
        account.email,
        account.first_name,
        account.last_name,
        account.school
      ].join(' ').toLowerCase();
      
      return searchableText.includes(query.toLowerCase());
    });

    this.renderFilteredAccounts(filteredAccounts);
  }

  handleFilter(filter) {
    let filteredAccounts = [...this.accounts];

    switch (filter) {
      case 'active':
        filteredAccounts = this.accounts.filter(acc => acc.is_active);
        break;
      case 'inactive':
        filteredAccounts = this.accounts.filter(acc => !acc.is_active);
        break;
      case 'recent':
        const oneWeekAgo = new Date();
        oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
        filteredAccounts = this.accounts.filter(acc => 
          new Date(acc.created_at) > oneWeekAgo
        );
        break;
      default:
        // 'all' or any other value shows all accounts
        break;
    }

    this.renderFilteredAccounts(filteredAccounts);
  }

  renderFilteredAccounts(accounts) {
    const tbody = document.getElementById('accountsTableBody');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (accounts.length === 0) {
      tbody.innerHTML = '<tr><td colspan="11" class="text-center">No accounts found</td></tr>';
      return;
    }

    accounts.forEach(account => {
      const row = this.createAccountRow(account);
      tbody.appendChild(row);
    });
  }

  exportData() {
    if (this.accounts.length === 0) {
      this.showMessage('No data to export', 'warning');
      return;
    }

    // Create CSV content
    const headers = ['ID', 'Username', 'Email', 'First Name', 'Last Name', 'Grade', 'School', 'Game Level', 'Experience', 'Status', 'Created At'];
    const csvContent = [
      headers.join(','),
      ...this.accounts.map(account => [
        account.id,
        account.username,
        account.email,
        account.first_name,
        account.last_name,
        account.grade,
        account.school,
        account.game_level,
        account.experience,
        account.is_active ? 'Active' : 'Inactive',
        new Date(account.created_at).toISOString()
      ].join(','))
    ].join('\n');

    // Download file
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `accounts_${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);

    this.showMessage('Data exported successfully!', 'success');
  }

  // Modal Management
  showModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.classList.add('active');
    }
  }

  closeModal() {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
      modal.classList.remove('active');
    });
  }

  // Utility Methods
  updateElement(id, value) {
    const element = document.getElementById(id);
    if (element) {
      element.textContent = value || 'N/A';
    }
  }

  showMessage(message, type = 'info') {
    // Remove existing alerts
    const existingAlerts = document.querySelectorAll('.alert');
    existingAlerts.forEach(alert => alert.remove());

    // Create new alert
    const alert = document.createElement('div');
    alert.className = `alert alert-${type}`;
    alert.textContent = message;

    // Insert at top of main content
    const main = document.querySelector('.main .container');
    if (main) {
      main.insertBefore(alert, main.firstChild);
    }

    // Auto remove after 5 seconds
    setTimeout(() => {
      alert.remove();
    }, 5000);
  }

  showLoading(show) {
    const loader = document.getElementById('loader');
    if (loader) {
      loader.style.display = show ? 'flex' : 'none';
    }
  }
}

// Global function for modal management
function showCreateModal() {
  window.adminPanel.showModal('createAccountModal');
}

// Initialize admin panel when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  window.adminPanel = new AdminPanel();
  
  // Setup edit form handler
  const editForm = document.getElementById('editAccountForm');
  if (editForm) {
    editForm.addEventListener('submit', (e) => {
      window.adminPanel.handleEditAccount(e);
    });
  }
});
