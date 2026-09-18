# Guía de Instalación

## 1. Objetivo

Esta guía explica cómo instalar y ejecutar el proyecto para validar el funcionamiento de las APIs desplegadas en los entornos de Docker, Podman y Containerd.

## 2. Requisitos previos

Antes de comenzar, asegúrate de tener instalado lo siguiente:

- Git
- Go 1.22 o superior
- Docker
- Podman
- Containerd o `nerdctl`
- curl
- acceso administrativo (`sudo` o permisos de root)

## 3. Clonar el repositorio

```bash
git clone <URL_DEL_REPOSITORIO>
cd 202300668_LAB_SO1_2S2026
```

## 4. Verificación de ambiente

### 4.1 Verificar Go

```bash
go version
```

Debe mostrar una versión compatible con Go 1.22 o posterior.

### 4.2 Verificar Docker

```bash
docker --version
```

### 4.3 Verificar Podman

```bash
podman --version
```

### 4.4 Verificar Containerd y nerdctl

```bash
containerd --version
nerdctl --version
```

## 5. Estructura del proyecto

```text
.
├── VM1_Containerd/
│   ├── API1/
│   └── API2/
├── VM2_Podman/
│   └── API3/
├── VM3_Docker/
├── README.md
├── MANUAL_TECNICO.md
└── GUIA_INSTALACION.md
```

## 6. Instalación de dependencias

### 6.1 Instalar Go en Ubuntu/Debian

```bash
sudo apt update
sudo apt install -y golang
```

### 6.2 Instalar Docker

```bash
sudo apt update
sudo apt install -y docker.io
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
newgrp docker
```

### 6.3 Instalar Podman

```bash
sudo apt update
sudo apt install -y podman
```

### 6.4 Instalar Containerd y nerdctl

```bash
sudo apt update
sudo apt install -y containerd
sudo apt install -y nerdctl
sudo systemctl enable --now containerd
```

## 7. Ejecución local sin contenedor

Puedes probar las APIs ejecutándolas directamente con Go.

### 7.1 API1

```bash
cd VM1_Containerd/API1
go run .
```

### 7.2 API2

```bash
cd VM1_Containerd/API2
go run .
```

### 7.3 API3

```bash
cd VM2_Podman/API3
go run .
```

### 7.4 VM3_Docker

```bash
cd VM3_Docker
go run .
```

Las aplicaciones quedarán disponibles en los siguientes puertos:

- `API1`: `http://localhost:8081/health`
- `API2`: `http://localhost:8082/health`
- `API3`: `http://localhost:8080/health`
- `VM3_Docker`: `http://localhost:8080/health`

## 8. Ejecución con Docker

### 8.1 Construir la imagen

```bash
cd VM3_Docker
docker build -t api-docker .
```

### 8.2 Levantar el contenedor

```bash
docker run --rm -d -p 8080:8080 --name api-docker api-docker
```

### 8.3 Verificar funcionamiento

```bash
curl -i http://localhost:8080/health
```

Salida esperada:

```json
{"message":"¡API funcionando correctamente!","status":"success","runtime":"Go Standard Library"}
```

## 9. Ejecución con Podman

### 9.1 Construir la imagen

```bash
cd VM2_Podman/API3
podman build -t api-podman .
```

### 9.2 Levantar el contenedor

```bash
podman run --rm -d -p 8080:8080 --name api-podman api-podman
```

### 9.3 Verificar funcionamiento

```bash
curl -i http://localhost:8080/health
```

Salida esperada:

```json
{"status":"UP","message":"API3 is Ready","timestamp":"2026-08-21T00:00:00Z","VM":"VM2","carnet":"202300668"}
```

## 10. Ejecución con Containerd

### 10.1 Construir imagen API1

```bash
cd VM1_Containerd/API1
nerdctl build -t api1:latest .
```

### 10.2 Levantar API1

```bash
nerdctl run --rm -d -p 8081:8081 --name api1 api1:latest
```

### 10.3 Construir imagen API2

```bash
cd ../API2
nerdctl build -t api2:latest .
```

### 10.4 Levantar API2

```bash
nerdctl run --rm -d -p 8082:8082 --name api2 api2:latest
```

### 10.5 Verificar funcionamiento

```bash
curl -i http://localhost:8081/health
curl -i http://localhost:8082/health
```

Salida esperada:

```json
{"status":"UP","message":"API1 is Ready","timestamp":"2026-08-21T00:00:00Z","VM":"VM1","carnet":"202300668"}
```

```json
{"status":"UP","message":"API2 is Ready","timestamp":"2026-08-21T00:00:00Z","VM":"VM1","carnet":"202300668"}
```

## 11. Capturas de pantalla para la entrega

Se recomienda guardar las capturas en esta ruta:

```text
./assets/capturas/
```

### 11.1 Captura de la API1

Archivo sugerido:

```text
./assets/capturas/curl-api1-health.png
```

Qué hacer:

1. Ejecuta: `curl -i http://localhost:8081/health`
2. Muestra la terminal completa.
3. Guarda la imagen con el comando y la respuesta JSON.

### 11.2 Captura de la API2

Archivo sugerido:

```text
./assets/capturas/curl-api2-health.png
```

Qué hacer:

1. Ejecuta: `curl -i http://localhost:8082/health`
2. Captura la terminal con la salida del JSON.
3. Asegúrate de que se observe el puerto `8082`.

### 11.3 Captura de la API3

Archivo sugerido:

```text
./assets/capturas/curl-api3-health.png
```

Qué hacer:

1. Ejecuta: `curl -i http://localhost:8080/health`
2. Muestra la respuesta del servicio de Podman.
3. Guarda la captura en la ruta indicada.

### 11.4 Captura final de Docker

Archivo sugerido:

```text
./assets/capturas/curl-docker-health.png
```

Qué hacer:

1. Ejecuta: `curl -i http://localhost:8080/health`
2. Comprueba que la salida muestre `status: "success"`.
3. Guarda la imagen para incluirla en el documento final.

## 12. Detección y limpieza

Para detener contenedores:

```bash
docker stop api-docker
podman stop api-podman
nerdctl stop api1 api2
```

## 13. Solución de problemas comunes

- `go: command not found`: instalar Go y verificar el `PATH`
- `docker: command not found`: revisar la instalación de Docker
- `podman: command not found`: instalar Podman
- `nerdctl: command not found`: instalar `nerdctl` o Containerd
- puerto ocupado: cambiar el mapeo `-p` o detener el servicio que lo está usando
- permisos: usar `sudo` para iniciar contenedores

## 14. Verificación final

La instalación será exitosa si se cumplen estas condiciones:

- cada servicio responde en la ruta `/health`
- cada puerto expuesto responde en `localhost`
- la salida devuelve JSON válido
- el contenedor queda levantado sin errores en consola

Si la respuesta incluye `status: "UP"` o `status: "success"`, la validación es satisfactoria.

## 15. Nota final

Las capturas de pantalla deben documentar la ejecución real de los comandos y la salida en consola para evidenciar el correcto funcionamiento del proyecto.
