# Dockerfile Element avec Office1789 SSO
FROM vectorim/element-web:latest

USER root

# Copier le script SSO
COPY office1789-sso.js /app/

# Copier et modifier index.html
RUN cp /app/index.html /tmp/index.html && \
    sed 's|</head>|<script src="/office1789-sso.js"></script></head>|' /tmp/index.html > /app/index.html.new && \
    mv /app/index.html.new /app/index.html

# Configuration Element
COPY config.json /app/config.json

USER nginx
