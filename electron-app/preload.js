const { contextBridge, ipcRenderer } = require('electron')

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
contextBridge.exposeInMainWorld('electron', {
  platform: process.platform,
  versions: {
    node: process.versions.node,
    chrome: process.versions.chrome,
    electron: process.versions.electron
  }
})

// Expose APIs for opening Mail and Chat windows
contextBridge.exposeInMainWorld('electronAPI', {
  openMailWindow: (url) => ipcRenderer.send('open-mail-window', url),
  openChatWindow: (url) => ipcRenderer.send('open-chat-window', url)
})
