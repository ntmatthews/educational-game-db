// Educational Game Database - Student Interface
class EduGameDB {
  constructor() {
    this.apiBase = '/api';
    this.currentUser = null;
    this.init();
  }

  async init() {
    this.registerServiceWorker();
    this.setupEventListeners();
    this.loadUserSession();
  }

  // Service Worker Registration for PWA
  async registerServiceWorker() {
    if ('serviceWorker' in navigator) {
      try {
        const registration = await navigator.serviceWorker.register('/static/sw.js');
        console.log('Service Worker registered:', registration);
      } catch (error) {
        console.error('Service Worker registration failed:', error);
      }
    }
  }

  // Event Listeners
  setupEventListeners() {
    // Login form
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
      loginForm.addEventListener('submit', (e) => this.handleLogin(e));
    }

    // Register form
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
      registerForm.addEventListener('submit', (e) => this.handleRegister(e));
    }

    // Logout button
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
      logoutBtn.addEventListener('click', () => this.handleLogout());
    }

    // Profile update
    const profileForm = document.getElementById('profileForm');
    if (profileForm) {
      profileForm.addEventListener('submit', (e) => this.handleProfileUpdate(e));
    }

    // Install PWA button
    this.setupPWAInstall();
  }

  // PWA Installation
  setupPWAInstall() {
    let deferredPrompt;
    const installBtn = document.getElementById('installBtn');

    window.addEventListener('beforeinstallprompt', (e) => {
      e.preventDefault();
      deferredPrompt = e;
      
      if (installBtn) {
        installBtn.style.display = 'block';
        installBtn.addEventListener('click', async () => {
          if (deferredPrompt) {
            deferredPrompt.prompt();
            const { outcome } = await deferredPrompt.userChoice;
            console.log(`User response to install prompt: ${outcome}`);
            deferredPrompt = null;
            installBtn.style.display = 'none';
          }
        });
      }
    });

    window.addEventListener('appinstalled', () => {
      console.log('PWA was installed');
      if (installBtn) {
        installBtn.style.display = 'none';
      }
    });
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

  // Authentication
  async handleLogin(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const credentials = {
      username: formData.get('username'),
      password: formData.get('password')
    };

    try {
      this.showLoading(true);
      const result = await this.apiCall('/login', 'POST', credentials);
      
      this.currentUser = result.account;
      this.saveUserSession();
      this.showMessage('Login successful!', 'success');
      
      // Redirect to dashboard
      setTimeout(() => {
        this.showDashboard();
      }, 1000);

    } catch (error) {
      this.showMessage(error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  async handleRegister(e) {
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

    // Validate password confirmation
    const passwordConfirm = formData.get('passwordConfirm');
    if (accountData.password !== passwordConfirm) {
      this.showMessage('Passwords do not match', 'error');
      return;
    }

    try {
      this.showLoading(true);
      const result = await this.apiCall('/accounts', 'POST', accountData);
      
      this.showMessage('Account created successfully! Please log in.', 'success');
      
      // Switch to login form
      setTimeout(() => {
        this.showLogin();
      }, 2000);

    } catch (error) {
      this.showMessage(error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  handleLogout() {
    this.currentUser = null;
    this.clearUserSession();
    this.showLogin();
    this.showMessage('Logged out successfully', 'success');
  }

  async handleProfileUpdate(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const updateData = {
      first_name: formData.get('firstName'),
      last_name: formData.get('lastName'),
      grade: parseInt(formData.get('grade')) || 0,
      school: formData.get('school'),
      game_level: this.currentUser.game_level,
      experience: this.currentUser.experience,
      is_active: this.currentUser.is_active
    };

    try {
      this.showLoading(true);
      const result = await this.apiCall(`/accounts/${this.currentUser.id}`, 'PUT', updateData);
      
      this.currentUser = result;
      this.saveUserSession();
      this.showMessage('Profile updated successfully!', 'success');
      this.updateProfileDisplay();

    } catch (error) {
      this.showMessage(error.message, 'error');
    } finally {
      this.showLoading(false);
    }
  }

  // Session Management
  saveUserSession() {
    if (this.currentUser) {
      localStorage.setItem('eduGameDB_user', JSON.stringify(this.currentUser));
    }
  }

  loadUserSession() {
    const saved = localStorage.getItem('eduGameDB_user');
    if (saved) {
      this.currentUser = JSON.parse(saved);
      this.showDashboard();
    } else {
      this.showLogin();
    }
  }

  clearUserSession() {
    localStorage.removeItem('eduGameDB_user');
  }

  // UI Management
  showLogin() {
    this.hideAllSections();
    const loginSection = document.getElementById('loginSection');
    if (loginSection) {
      loginSection.style.display = 'block';
    }
  }

  showRegister() {
    this.hideAllSections();
    const registerSection = document.getElementById('registerSection');
    if (registerSection) {
      registerSection.style.display = 'block';
    }
  }

  showDashboard() {
    this.hideAllSections();
    const dashboardSection = document.getElementById('dashboardSection');
    if (dashboardSection) {
      dashboardSection.style.display = 'block';
      this.updateDashboard();
    }
  }

  hideAllSections() {
    const sections = ['loginSection', 'registerSection', 'dashboardSection'];
    sections.forEach(id => {
      const section = document.getElementById(id);
      if (section) {
        section.style.display = 'none';
      }
    });
  }

  updateDashboard() {
    if (!this.currentUser) return;

    // Update user info display
    this.updateElement('userFullName', `${this.currentUser.first_name} ${this.currentUser.last_name}`);
    this.updateElement('userUsername', this.currentUser.username);
    this.updateElement('userEmail', this.currentUser.email);
    this.updateElement('userGrade', this.currentUser.grade);
    this.updateElement('userSchool', this.currentUser.school);
    this.updateElement('userLevel', this.currentUser.game_level);
    this.updateElement('userExperience', this.currentUser.experience);
    this.updateElement('userJoined', new Date(this.currentUser.created_at).toLocaleDateString());

    // Update profile form
    this.updateProfileForm();
  }

  updateProfileForm() {
    if (!this.currentUser) return;

    this.updateFormField('firstName', this.currentUser.first_name);
    this.updateFormField('lastName', this.currentUser.last_name);
    this.updateFormField('grade', this.currentUser.grade);
    this.updateFormField('school', this.currentUser.school);
  }

  updateProfileDisplay() {
    this.updateDashboard();
  }

  updateElement(id, value) {
    const element = document.getElementById(id);
    if (element) {
      element.textContent = value || 'N/A';
    }
  }

  updateFormField(name, value) {
    const field = document.querySelector(`[name="${name}"]`);
    if (field) {
      field.value = value || '';
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

  // Game Progress Methods
  async updateGameProgress(level, experience) {
    if (!this.currentUser) return;

    const updateData = {
      first_name: this.currentUser.first_name,
      last_name: this.currentUser.last_name,
      grade: this.currentUser.grade,
      school: this.currentUser.school,
      game_level: level,
      experience: experience,
      is_active: this.currentUser.is_active
    };

    try {
      const result = await this.apiCall(`/accounts/${this.currentUser.id}`, 'PUT', updateData);
      this.currentUser = result;
      this.saveUserSession();
      this.updateDashboard();
      return true;
    } catch (error) {
      console.error('Failed to update game progress:', error);
      return false;
    }
  }

  // Offline Support
  handleOffline() {
    this.showMessage('You are currently offline. Some features may be limited.', 'warning');
  }

  handleOnline() {
    this.showMessage('Connection restored!', 'success');
  }
}

// Global functions for navigation
function showLogin() {
  window.eduGameDB.showLogin();
}

function showRegister() {
  window.eduGameDB.showRegister();
}

function showDashboard() {
  window.eduGameDB.showDashboard();
}

// Network status monitoring
window.addEventListener('online', () => {
  if (window.eduGameDB) {
    window.eduGameDB.handleOnline();
  }
});

window.addEventListener('offline', () => {
  if (window.eduGameDB) {
    window.eduGameDB.handleOffline();
  }
});

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  window.eduGameDB = new EduGameDB();
});
