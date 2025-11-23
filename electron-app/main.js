const { app, BrowserWindow, session, ipcMain } = require('electron')
const path = require('path')

let mainWindow = null
let mailWindow = null
let chatWindow = null

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1024,
    minHeight: 768,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true,
      webSecurity: true,
      partition: 'persist:office1789'  // Session persistante
    },
    icon: path.join(__dirname, 'www', 'logo.png'),
    title: 'Office1789',
    backgroundColor: '#ffffff',
    show: false  // Don't show until ready
  })

  // Configure session for better security
  const ses = session.fromPartition('persist:office1789')
  ses.webRequest.onHeadersReceived((details, callback) => {
    callback({
      responseHeaders: {
        ...details.responseHeaders,
        'Content-Security-Policy': [
          "default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob: https://*.office1789.com; " +
          "connect-src 'self' http://localhost:* https://office1789.com https://backend.office1789.com https://*.office1789.com wss://*.office1789.com; " +
          "img-src 'self' data: blob: http://localhost:* https://office1789.com https://backend.office1789.com https://*.office1789.com; " +
          "script-src 'self' 'unsafe-inline' 'unsafe-eval' https://*.office1789.com; " +
          "style-src 'self' 'unsafe-inline' https://*.office1789.com; " +
          "font-src 'self' data: https://*.office1789.com; " +
          "media-src 'self' data: blob: https://*.office1789.com;"
        ]
      }
    })
  })

  // Load the app
  mainWindow.loadFile(path.join(__dirname, 'www', 'index.html'))

  // Show window when ready
  mainWindow.once('ready-to-show', () => {
    mainWindow.show()
    mainWindow.focus()
  })

  // Handle external links
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    require('electron').shell.openExternal(url)
    return { action: 'deny' }
  })

  // Open DevTools in development
  if (process.env.NODE_ENV === 'development') {
    mainWindow.webContents.openDevTools()
  }

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

app.whenReady().then(() => {
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// Handle app protocol for better security
app.setAsDefaultProtocolClient('office1789')

// IPC handlers for opening Mail and Chat windows
ipcMain.on('open-mail-window', (event, url) => {
  if (mailWindow && !mailWindow.isDestroyed()) {
    mailWindow.focus()
    mailWindow.loadURL(url)
    return
  }

  mailWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    title: 'Office1789 - Mail',
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      webSecurity: true,
      partition: 'persist:office1789'
    },
    parent: mainWindow,
    modal: false
  })

  mailWindow.loadURL(url)
  
  mailWindow.on('closed', () => {
    mailWindow = null
  })
})

ipcMain.on('open-chat-window', (event, url) => {
  if (chatWindow && !chatWindow.isDestroyed()) {
    chatWindow.focus()
    chatWindow.loadURL(url)
    return
  }

  chatWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    title: 'Office1789 - Chat',
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      webSecurity: true,
      partition: 'persist:office1789'
    },
    parent: mainWindow,
    modal: false
  })

  chatWindow.loadURL(url)
  
  chatWindow.on('closed', () => {
    chatWindow = null
  })
})
