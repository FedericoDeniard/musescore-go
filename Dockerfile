# Etapa 1: Construir el frontend con Node.js
FROM node:22 as frontend-builder
WORKDIR /app

COPY src/static/frontend/package*.json ./src/static/frontend/
WORKDIR /app/src/static/frontend
RUN npm install
COPY src/static/frontend/ ./
RUN npm run build

# Etapa 2: Construir y correr el backend con Go
FROM golang:1.23.4

# Instalar dependencias del sistema (para rod/chromium)
RUN apt-get update && apt-get install -y \
    chromium \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libc6 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libexpat1 \
    libfontconfig1 \
    libgbm1 \
    libgcc1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libstdc++6 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    lsb-release \
    wget \
    xdg-utils \
    libvips-dev \
    inkscape \
    --no-install-recommends

# Crear carpeta de trabajo y copiar todo el código del proyecto
WORKDIR /app
COPY . .

# Copiar build del frontend desde la etapa 1
COPY --from=frontend-builder /app/src/static/frontend/dist ./src/static/frontend/dist

# Compilar el binario de Go
RUN go mod download
RUN go build -o dist/main src/main.go

EXPOSE 8000

# Ejecutar el binario
CMD ["./dist/main"]
