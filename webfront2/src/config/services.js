// URLs des services (mail, chat, docs, matrix, roundcube) strictement issues des variables d'environnement
export const SERVICES_CONFIG = {
  getMailURL() {
    return import.meta.env.VITE_ROUNDCUBE_URL
  },
  getChatURL() {
    return import.meta.env.VITE_CHAT_URL
  },
  getDocsURL() {
    return import.meta.env.VITE_DOCS_URL
  },
  getMatrixURL() {
    return import.meta.env.VITE_MATRIX_URL
  }
}
