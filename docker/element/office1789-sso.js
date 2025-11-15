// Office1789 SSO pour Element/Matrix
// Ce script intercepte le token SSO dans l'URL et authentifie automatiquement

(function() {
    'use strict';
    
    console.log('[Office1789-SSO] Plugin chargé');
    
    // Secret partagé avec le backend Go
    const SSO_SECRET = 'Office1789-Matrix-SecretKey-ChangeInProduction';
    
    // Fonction pour décoder base64url
    function base64urlDecode(str) {
        // Remplacer les caractères base64url par base64 standard
        str = str.replace(/-/g, '+').replace(/_/g, '/');
        // Ajouter le padding si nécessaire
        while (str.length % 4) {
            str += '=';
        }
        return atob(str);
    }
    
    // Fonction pour vérifier la signature HMAC-SHA256
    async function verifyToken(token) {
        const parts = token.split('.');
        if (parts.length !== 2) {
            console.error('[Office1789-SSO] Token invalide - format incorrect');
            return null;
        }
        
        const [claimsB64, signatureB64] = parts;
        
        try {
            // Décoder les claims
            const claimsJSON = base64urlDecode(claimsB64);
            const claims = JSON.parse(claimsJSON);
            
            console.log('[Office1789-SSO] Claims décodés:', claims);
            
            // Vérifier l'expiration
            const now = Math.floor(Date.now() / 1000);
            if (claims.exp && claims.exp < now) {
                console.error('[Office1789-SSO] Token expiré');
                return null;
            }
            
            // Vérifier la signature avec Web Crypto API
            const encoder = new TextEncoder();
            const key = await crypto.subtle.importKey(
                'raw',
                encoder.encode(SSO_SECRET),
                { name: 'HMAC', hash: 'SHA-256' },
                false,
                ['sign', 'verify']
            );
            
            const signature = Uint8Array.from(atob(signatureB64.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
            const data = encoder.encode(claimsB64);
            
            const valid = await crypto.subtle.verify(
                'HMAC',
                key,
                signature,
                data
            );
            
            if (!valid) {
                console.error('[Office1789-SSO] Signature invalide');
                return null;
            }
            
            console.log('[Office1789-SSO] Token valide ✓');
            return claims;
            
        } catch (error) {
            console.error('[Office1789-SSO] Erreur validation token:', error);
            return null;
        }
    }
    
    // Fonction pour effectuer l'authentification Matrix
    async function authenticateMatrix(userID, password) {
        console.log('[Office1789-SSO] Tentative d\'authentification pour:', userID);
        
        try {
            // Utiliser l'API Matrix pour se connecter
            const response = await fetch('http://localhost:8008/_matrix/client/r0/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    type: 'm.login.password',
                    user: userID,
                    password: password,
                    initial_device_display_name: 'Office1789 SSO'
                })
            });
            
            if (response.ok) {
                const data = await response.json();
                console.log('[Office1789-SSO] Authentification réussie ✓');
                
                // Stocker les credentials dans localStorage pour Element
                localStorage.setItem('mx_access_token', data.access_token);
                localStorage.setItem('mx_user_id', data.user_id);
                localStorage.setItem('mx_device_id', data.device_id);
                localStorage.setItem('mx_hs_url', 'http://localhost:8008');
                
                // Rediriger vers Element
                window.location.href = '/';
                
            } else {
                const error = await response.json();
                console.error('[Office1789-SSO] Erreur authentification:', error);
                alert('Erreur d\'authentification Matrix: ' + (error.error || 'Erreur inconnue'));
            }
            
        } catch (error) {
            console.error('[Office1789-SSO] Erreur réseau:', error);
            alert('Impossible de se connecter au serveur Matrix');
        }
    }
    
    // Vérifier si un token SSO est présent dans l'URL
    async function checkSSOToken() {
        const urlParams = new URLSearchParams(window.location.search);
        const ssoToken = urlParams.get('sso_token');
        
        if (ssoToken) {
            console.log('[Office1789-SSO] Token SSO détecté');
            
            // Afficher un message de chargement
            document.body.innerHTML = `
                <div style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; font-family: Arial, sans-serif; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white;">
                    <div style="text-align: center; background: rgba(255, 255, 255, 0.1); padding: 48px; border-radius: 24px; backdrop-filter: blur(10px);">
                        <h1 style="font-size: 2rem; margin-bottom: 24px;">🔐 Connexion Office1789</h1>
                        <div class="spinner" style="width: 50px; height: 50px; border: 5px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 24px;"></div>
                        <p style="font-size: 1.1rem;">Authentification automatique en cours...</p>
                    </div>
                </div>
                <style>
                    @keyframes spin {
                        to { transform: rotate(360deg); }
                    }
                </style>
            `;
            
            // Valider et utiliser le token
            const claims = await verifyToken(ssoToken);
            
            if (claims && claims.user_id && claims.password) {
                await authenticateMatrix(claims.user_id, claims.password);
            } else {
                document.body.innerHTML = `
                    <div style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; font-family: Arial, sans-serif; background: #f5f5f7; color: #333;">
                        <div style="text-align: center; background: white; padding: 48px; border-radius: 24px; box-shadow: 0 8px 32px rgba(0,0,0,0.1);">
                            <h1 style="font-size: 2rem; margin-bottom: 16px; color: #e74c3c;">❌ Erreur SSO</h1>
                            <p style="font-size: 1.1rem; margin-bottom: 32px;">Token SSO invalide ou expiré</p>
                            <a href="http://localhost:5173" style="display: inline-block; padding: 12px 32px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; text-decoration: none; border-radius: 24px; font-weight: 600;">
                                Retour à Office1789
                            </a>
                        </div>
                    </div>
                `;
            }
        }
    }
    
    // Exécuter au chargement de la page
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', checkSSOToken);
    } else {
        checkSSOToken();
    }
    
})();
