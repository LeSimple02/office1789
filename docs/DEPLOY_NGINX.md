# Configuration Nginx pour Office1789

Ce document décrit la configuration nginx nécessaire pour exposer les services Office1789 en production.

## Architecture

- **Backend API**: `backend.office1789.com` → `http://127.0.0.1:8080`
- **Frontend**: `office1789.com` → `http://127.0.0.1:5173`
- **Mail (Roundcube)**: `mail.office1789.com` → `http://127.0.0.1:8081`
- **Chat (Element)**: `chat.office1789.com` → `http://127.0.0.1:8083`
- **Matrix API**: `chat.office1789.com/_matrix/` → `http://127.0.0.1:8008`

## Configuration nginx

Créer ou éditer `/etc/nginx/sites-available/office1789`:

```nginx
# Backend API
server {
    listen 80;
    server_name backend.office1789.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers (si nécessaire)
        add_header Access-Control-Allow-Origin * always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Origin, Content-Type, Accept, Authorization" always;
    }
}

# Frontend
server {
    listen 80;
    server_name office1789.com www.office1789.com;

    location / {
        proxy_pass http://127.0.0.1:5173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Mail (Roundcube)
server {
    listen 80;
    server_name mail.office1789.com;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Chat (Element + Matrix)
server {
    listen 80;
    server_name chat.office1789.com;

    # Element Web
    location / {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Matrix API
    location /_matrix/ {
        proxy_pass http://127.0.0.1:8008/_matrix/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support pour Matrix sync
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Activation et test

```bash
# Lien symbolique vers sites-enabled
sudo ln -sf /etc/nginx/sites-available/office1789 /etc/nginx/sites-enabled/office1789

# Test de la configuration
sudo nginx -t

# Rechargement nginx
sudo systemctl reload nginx

# Si nginx n'est pas démarré
sudo systemctl start nginx
sudo systemctl enable nginx
```

## Passage en HTTPS avec Let's Encrypt

Une fois les services accessibles en HTTP, activer HTTPS:

```bash
cd /home/debian/office1789
./scripts/prod_letsencrypt.sh
```

Puis adapter la configuration nginx pour écouter sur le port 443 avec les certificats générés:

```nginx
server {
    listen 443 ssl http2;
    server_name backend.office1789.com;
    
    ssl_certificate /etc/letsencrypt/live/office1789.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/office1789.com/privkey.pem;
    
    # ... reste de la config proxy_pass
}
```

## Vérification

Test depuis votre machine locale:

```bash
# Backend
curl https://backend.office1789.com/api/welcome

# Frontend
curl -I https://office1789.com

# Mail
curl -I https://mail.office1789.com

# Chat
curl -I https://chat.office1789.com
```

## Dépannage

### nginx ne démarre pas

```bash
sudo journalctl -u nginx -n 50 --no-pager
sudo nginx -t
```

### Service non accessible

Vérifier que le conteneur Docker correspondant tourne:

```bash
cd /home/debian/office1789/docker
docker compose ps
```

### Port déjà utilisé

Vérifier les processus écoutant sur les ports:

```bash
sudo ss -tlnp | grep ':80'
sudo ss -tlnp | grep ':8080'
```
