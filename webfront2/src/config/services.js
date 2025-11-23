// Configuration des URLs des services (mail, chat, docs)
export const SERVICES_CONFIG = {
  getMailURL() {
    // En dev local
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      return 'http://localhost:8081'
    }
    // En production, utiliser la variable d'environnement ou le domaine par défaut
    return import.meta.env.VITE_APP_MAIL_URL || 'https://mail.office1789.com'
  },
  
  getChatURL() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      return 'http://localhost:8083'
    }
    return import.meta.env.VITE_APP_CHAT_URL || 'https://chat.office1789.com'
  },
  
  getDocsURL() {
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      return 'http://localhost:8082'
    }
    return import.meta.env.VITE_APP_DOCS_URL || 'https://docs.office1789.com'
  }
}
