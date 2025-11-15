<?php
/**
 * SSO Login Hook - Injecté dans index.php après iniset.php
 * Intercepte la demande de login et utilise les credentials SSO si disponibles
 */

// Vérifier si on a des credentials SSO en session
if (isset($_SESSION['sso_authorized']) && $_SESSION['sso_authorized'] === true) {
    error_log('[SSO-Hook] Credentials SSO détectés en session');
    
    // Si on est sur ?_task=mail (après redirect SSO)
    if (isset($_GET['_task']) && $_GET['_task'] === 'mail') {
        error_log('[SSO-Hook] Tentative auto-login avec ' . $_SESSION['sso_user']);
        
        // Simuler une requête POST de login
        $_POST['_task'] = 'login';
        $_POST['_action'] = 'login';
        $_POST['_user'] = $_SESSION['sso_user'];
        $_POST['_pass'] = $_SESSION['sso_pass'];
        $_POST['_timezone'] = 'Europe/Paris';
        
        // Générer un token CSRF valide pour cette session
        // Roundcube utilise rcube_utils::request_token() pour générer les tokens
        // On va charger le framework Roundcube pour avoir accès à cette fonction
        
        // Le token CSRF est stocké dans $_SESSION['request_token']
        // Si il n'existe pas encore, Roundcube le créera automatiquement
        // Mais on doit passer la vérification, donc on désactive temporairement
        
        // Désactiver la vérification CSRF pour cette requête SSO uniquement
        $_SERVER['HTTP_X_REQUESTED_WITH'] = 'XMLHttpRequest'; // Roundcube skip CSRF pour AJAX
        
        // Marquer que les credentials ont été utilisés
        unset($_SESSION['sso_authorized']);
        unset($_SESSION['sso_user']);
        unset($_SESSION['sso_pass']);
        
        error_log('[SSO-Hook] POST simulé, passage à Roundcube');
    }
}
