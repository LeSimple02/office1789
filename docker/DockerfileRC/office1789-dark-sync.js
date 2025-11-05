// Détection et synchronisation du mode dark Office1789
(function() {
  // Lire le paramètre dark depuis l'URL
  const urlParams = new URLSearchParams(window.location.search);
  const darkMode = urlParams.get('dark') === '1';
  
  // Appliquer la classe dark sur html
  if (darkMode) {
    document.documentElement.classList.add('dark');
  }
  
  // Écouter les messages de l'app parent pour changer le mode en temps réel
  window.addEventListener('message', function(event) {
    // Vérifier l'origine pour sécurité (optionnel)
    if (event.data && event.data.type === 'darkModeChange') {
      if (event.data.isDark) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  });
})();
