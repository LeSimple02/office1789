# 🏛 Office1789

Office1789 est une application web qui se voudra à terme être une alternative à office365 basée sur **Vue.js** (frontend), **Gin (Go)** (backend) et **PostgreSQL** (base de données).  
L’ensemble est orchestré avec **Docker** et **Docker Compose** pour simplifier le développement et le déploiement.

---

## 🚀 Stack technique

- **Frontend** : Vue.js  
- **Backend** : Gin (Go)  
- **Base de données** : PostgreSQL  
- **Conteneurisation** : Docker

---

## ⚙️ Prérequis

python3 config.py

---

## ▶️ Lancer le projet en local

1. **Cloner le dépôt**
   ```bash
   git clone https://github.com/LeSimple02/office1789.git
   cd office1789

   cd webfront2
   npm install
   npm run serve

   cd backend
   go run .