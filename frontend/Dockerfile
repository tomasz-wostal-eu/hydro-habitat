# --- Build Stage ---
FROM node:24-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install --legacy-peer-deps
COPY . .
RUN npm run build --legacy-peer-deps

# --- Final Stage ---
FROM caddy:2-alpine
RUN rm /etc/caddy/Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /app/dist /usr/share/caddy
EXPOSE 80
