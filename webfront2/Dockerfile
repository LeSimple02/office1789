# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments pour les URLs de production
ARG VITE_APP_API
ARG VITE_APP_MAIL_URL
ARG VITE_APP_CHAT_URL
ARG VITE_APP_DOCS_URL

# Les exporter comme variables d'environnement pour le build
ENV VITE_APP_API=$VITE_APP_API
ENV VITE_APP_MAIL_URL=$VITE_APP_MAIL_URL
ENV VITE_APP_CHAT_URL=$VITE_APP_CHAT_URL
ENV VITE_APP_DOCS_URL=$VITE_APP_DOCS_URL
ENV NODE_ENV=production

# Build the application
RUN npm run build

# Runtime stage - nginx to serve static files
FROM nginx:alpine

# Copy built files from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration for SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
