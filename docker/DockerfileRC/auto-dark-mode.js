// Office1789 - Auto Dark Mode pour Roundcube
(function() {
  // Vérifier si le paramètre dark=1 est dans l'URL
  const urlParams = new URLSearchParams(window.location.search);
  const isDarkMode = urlParams.get('dark') === '1';
  
  if (isDarkMode) {
    // Ajouter la classe dark-mode au HTML
    document.documentElement.classList.add('dark-mode');
    
    // Sauvegarder dans localStorage pour persistance
    localStorage.setItem('skin-mode', 'dark');
  }
  
  // Écouter les messages du parent (Vue.js) pour changer le mode
  window.addEventListener('message', function(event) {
    if (event.data && event.data.type === 'darkModeChange') {
      if (event.data.isDark) {
        document.documentElement.classList.add('dark-mode');
        localStorage.setItem('skin-mode', 'dark');
      } else {
        document.documentElement.classList.remove('dark-mode');
        localStorage.setItem('skin-mode', 'light');
      }
    }
  });
})();
