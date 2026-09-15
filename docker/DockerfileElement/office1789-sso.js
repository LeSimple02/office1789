// Office1789 SSO pour Element/Matrix
// Ce script intercepte le token SSO dans l'URL et authentifie automatiquement
// SECURITE : bloque l'accès direct à la page de login

(function() {
    'use strict';

    console.log('[Office1789-SSO] Plugin chargé');

    // Secret partagé avec le backend Go
    const SSO_SECRET = 'Office1789-Matrix-SecretKey-ChangeInProduction';

    // Détecter si on est en prod (domaine office1789.com) ou en local
    const DOMAIN = window.location.hostname || '';
    const IS_PROD = DOMAIN.endsWith('office1789.com');

    // URL du frontend Office1789 (pour les liens "Retour")
    const FRONT_BASE = IS_PROD ? 'https://office1789.com' : 'http://localhost:5173';
    // Homeserver Matrix (prod ou local)
    const HS_BASE = IS_PROD ? 'https://matrix.office1789.com' : 'http://localhost:8008';

    // ===== FONCTIONS UTILITAIRES =====

    function base64urlDecode(str) {
        str = str.replace(/-/g, '+').replace(/_/g, '/');
        while (str.length % 4) {
            str += '=';
        }
        return atob(str);
    }

    async function verifyToken(token) {
        const parts = token.split('.');
        if (parts.length !== 2) {
            console.error('[Office1789-SSO] Token invalide - format incorrect');
            return null;
        }

        const [claimsB64, signatureB64] = parts;

        try {
            const claimsJSON = base64urlDecode(claimsB64);
            const claims = JSON.parse(claimsJSON);

            console.log('[Office1789-SSO] Claims décodés:', claims);

            const now = Math.floor(Date.now() / 1000);
            if (claims.exp && claims.exp < now) {
                console.error('[Office1789-SSO] Token expiré');
                return null;
            }

            const encoder = new TextEncoder();
            const key = await crypto.subtle.importKey(
                'raw',
                encoder.encode(SSO_SECRET),
                { name: 'HMAC', hash: 'SHA-256' },
                false,
                ['sign', 'verify']
            );

            const signature = Uint8Array.from(
                atob(signatureB64.replace(/-/g, '+').replace(/_/g, '/')),
                c => c.charCodeAt(0),
            );
            const data = encoder.encode(claimsB64);

            const valid = await crypto.subtle.verify('HMAC', key, signature, data);

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

    // ===== AUTHENTIFICATION MATRIX =====

    async function authenticateMatrix(username, matrixUserId, password) {
        console.log('[Office1789-SSO] Tentative d\'authentification pour:', matrixUserId);

        const existingDeviceId = localStorage.getItem('mx_device_id');

        try {
            const loginPayload = {
                type: 'm.login.password',
                user: matrixUserId,
                password: password,
                initial_device_display_name: 'Office1789 SSO',
            };

            if (existingDeviceId) {
                loginPayload.device_id = existingDeviceId;
                console.log('[Office1789-SSO] Réutilisation du device_id existant:', existingDeviceId);
            }

            const response = await fetch(HS_BASE + '/_matrix/client/r0/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(loginPayload),
            });

            if (response.ok) {
                const data = await response.json();
                console.log('[Office1789-SSO] Authentification réussie');

                localStorage.setItem('mx_access_token', data.access_token);
                localStorage.setItem('mx_user_id', data.user_id);
                localStorage.setItem('mx_device_id', data.device_id);
                localStorage.setItem('mx_hs_url', HS_BASE);
                localStorage.setItem('office1789_matrix_auth', Date.now().toString());

                window.location.href = '/';
            } else if (response.status === 429) {
                console.warn('[Office1789-SSO] Rate limit détecté');
                document.body.innerHTML = `
                    <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh;font-family:Arial,sans-serif;background:linear-gradient(135deg,#f093fb 0%,#f5576c 100%);color:white;">
                        <div style="text-align:center;background:rgba(255,255,255,0.1);padding:48px;border-radius:24px;backdrop-filter:blur(10px);">
                            <h1 style="font-size:2rem;margin-bottom:24px;">Trop de tentatives</h1>
                            <p style="font-size:1.1rem;margin-bottom:16px;">Veuillez patienter quelques instants avant de réessayer.</p>
                        </div>
                    </div>
                `;
            } else {
                console.error('[Office1789-SSO] Erreur HTTP Matrix:', response.status);
                document.body.innerHTML = `
                    <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh;font-family:Arial,sans-serif;background:#f5f5f7;color:#333;">
                        <div style="text-align:center;background:white;padding:48px;border-radius:24px;box-shadow:0 8px 32px rgba(0,0,0,0.1);">
                            <h1 style="font-size:2rem;margin-bottom:16px;color:#e74c3c;">Erreur de connexion Matrix</h1>
                            <p style="font-size:1.1rem;margin-bottom:32px;">Impossible de vous connecter au serveur Matrix.</p>
                            <a href="https://office1789.com" style="display:inline-block;padding:12px 32px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;text-decoration:none;border-radius:24px;font-weight:600;">Retour à Office1789</a>
                        </div>
                    </div>
                `;
            }
        } catch (error) {
            console.error('[Office1789-SSO] Erreur réseau Matrix:', error);
            document.body.innerHTML = `
                <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh;font-family:Arial,sans-serif;background:#f5f5f7;color:#333;">
                    <div style="text-align:center;background:white;padding:48px;border-radius:24px;box-shadow:0 8px 32px rgba(0,0,0,0.1);">
                        <h1 style="font-size:2rem;margin-bottom:16px;color:#e74c3c;">Erreur réseau</h1>
                        <p style="font-size:1.1rem;margin-bottom:32px;">Impossible de contacter le serveur Matrix.</p>
                        <a href="${FRONT_BASE}" style="display:inline-block;padding:12px 32px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;text-decoration:none;border-radius:24px;font-weight:600;">Retour à Office1789</a>
                    </div>
                </div>
            `;
        }
    }

    // ===== FLUX SSO =====

    const urlParams = new URLSearchParams(window.location.search);
    const ssoToken = urlParams.get('sso_token');

    function isLoginPage() {
        const hash = window.location.hash;
        return (
            hash === '' ||
            hash === '#/' ||
            hash === '#/login' ||
            hash === '#/welcome' ||
            hash === '#/register'
        );
    }

    function isAuthenticated() {
        return !!localStorage.getItem('mx_access_token');
    }

    async function processSSO(token) {
        console.log('[Office1789-SSO] Traitement du token SSO');
        const claims = await verifyToken(token);

        if (claims && claims.username && claims.matrixUserId && claims.password) {
            await authenticateMatrix(claims.username, claims.matrixUserId, claims.password);
        } else {
            document.body.innerHTML = `
                <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh;font-family:Arial,sans-serif;background:#f5f5f7;color:#333;">
                    <div style="text-align:center;background:white;padding:48px;border-radius:24px;box-shadow:0 8px 32px rgba(0,0,0,0.1);">
                        <h1 style="font-size:2rem;margin-bottom:16px;color:#e74c3c;">Erreur SSO</h1>
                        <p style="font-size:1.1rem;margin-bottom:32px;">Token SSO invalide ou expiré.</p>
                        <a href="${FRONT_BASE}" style="display:inline-block;padding:12px 32px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;text-decoration:none;border-radius:24px;font-weight:600;">Retour à Office1789</a>
                    </div>
                </div>
            `;
        }
    }

    function init() {
        if (ssoToken) {
            console.log('[Office1789-SSO] Token SSO détecté, traitement immédiat...');
            processSSO(ssoToken);
        } else if (isLoginPage() && !isAuthenticated()) {
            console.warn('[Office1789-SSO] Accès à la page de login bloqué - pas de token SSO');

            const blockPage = document.createElement('div');
            blockPage.innerHTML = `
                <div style="position:fixed;top:0;left:0;width:100%;height:100%;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);display:flex;justify-content:center;align-items:center;z-index:9999999;font-family:Arial,sans-serif">
                    <div style="background:white;padding:48px;border-radius:16px;box-shadow:0 20px 60px rgba(0,0,0,0.3);text-align:center;max-width:500px">
                        <div style="font-size:72px;margin-bottom:24px">🚫</div>
                        <h1 style="color:#e74c3c;margin:0 0 24px 0;font-size:36px;font-weight:700">Accès refusé</h1>
                        <p style="color:#555;line-height:1.8;margin:16px 0;font-size:16px">L'accès direct à la page de connexion Element est désactivé pour des raisons de sécurité.</p>
                        <p style="color:#555;line-height:1.8;margin:16px 0;font-size:16px">Veuillez vous connecter via <strong>Office1789</strong> pour accéder au chat Matrix.</p>
                        <a href="${FRONT_BASE}" style="display:inline-block;margin-top:32px;padding:16px 40px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;text-decoration:none;border-radius:8px;font-weight:600;font-size:16px;transition:transform 0.2s">🏛️ Retour à Office1789</a>
                    </div>
                </div>
            `;

            if (document.body) {
                document.body.appendChild(blockPage);
            } else {
                document.addEventListener('DOMContentLoaded', () => document.body.appendChild(blockPage));
            }
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
