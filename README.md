# Guía de instalación - Proyecto 1 - Sistemas Operativos 1

Este repositorio contiene una pequeña infraestructura de APIs desarrolladas en Go para trabajar con tres motores de contenedores distintos: Docker, Podman y Containerd. La finalidad es levantar cada servicio de forma independiente y verificar que responden correctamente en sus respectivos puertos.

## 1. Requisitos previos

Antes de comenzar, asegúrate de tener instalado lo siguiente en tu sistema:

- Git
- Go 1.22 o superior
- Docker o Podman o Containerd, según el entorno que quieras probar
- curl (opcional, para verificar la salud de las APIs)
- Acceso root o sudo en Linux para instalar y configurar los motores de contenedores

## 2. Clonar el repositorio

```bash
git clone <URL_DEL_REPOSITORIO>
cd 202300668_LAB_SO1_2S2026
```

## 3. Estructura del proyecto

```text
.
├── VM1_Containerd/
│   ├── API1/
│   └── API2/
├── VM2_Podman/
│   └── API3/
└── VM3_Docker/
```

- VM1_Containerd: contiene dos servicios Go para ejecutar sobre Containerd.
- VM2_Podman: contiene la API 3 para ejecutar con Podman.
- VM3_Docker: contiene la API principal para ejecutar con Docker.

## 4. Instalación de dependencias

### 4.1 Instalar Go

En Ubuntu/Debian:

```bash
sudo apt update
sudo apt install -y golang
```

Verifica la versión:

```bash
go version
```

### 4.2 Instalar Docker

```bash
sudo apt update
sudo apt install -y docker.io
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
newgrp docker
```

### 4.3 Instalar Podman

```bash
sudo apt update
sudo apt install -y podman
```

### 4.4 Instalar Containerd

```bash
sudo apt update
sudo apt install -y containerd
sudo systemctl enable --now containerd
```

Si prefieres usar nerdctl para construir y ejecutar contenedores con Containerd:

```bash
sudo apt install -y nerdctl
```

## 5. Ejecución local sin contenedor (opción alternativa)

Cada API también puede ejecutarse directamente con Go:

```bash
cd VM1_Containerd/API1
go run .
```

En otra terminal:

```bash
cd VM1_Containerd/API2
go run .
```

```bash
cd VM2_Podman/API3
go run .
```

```bash
cd VM3_Docker
go run .
```

Las APIs responderán en los siguientes puertos:

- API1: 8081
- API2: 8082
- API3: 8080
- VM3_Docker: 8080

## 6. Ejecutar con Docker

### 6.1 Construir la imagen

```bash
cd VM3_Docker
docker build -t api-docker .
```

### 6.2 Levantar el contenedor

```bash
docker run --rm -d -p 8080:8080 --name api-docker api-docker
```

### 6.3 Verificar salud

```bash
curl http://localhost:8080/health
```

## 7. Ejecutar con Podman

### 7.1 Construir la imagen

```bash
cd VM2_Podman/API3
podman build -t api-podman .
```

### 7.2 Levantar el contenedor

```bash
podman run --rm -d -p 8080:8080 --name api-podman api-podman
```

### 7.3 Verificar salud

```bash
curl http://localhost:8080/health
```

## 8. Ejecutar con Containerd

### 8.1 Construir la imagen usando nerdctl

```bash
cd VM1_Containerd/API1
nerdctl build -t api1:latest .
```

```bash
cd ../API2
nerdctl build -t api2:latest .
```

### 8.2 Levantar los contenedores

```bash
nerdctl run --rm -d -p 8081:8081 --name api1 api1:latest
nerdctl run --rm -d -p 8082:8082 --name api2 api2:latest
```

### 8.3 Verificar salud

```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## 9. Deteener y limpiar contenedores

```bash
docker stop api-docker
podman stop api-podman
nerdctl stop api1 api2
```

## 10. Solución de problemas comunes

- Error de permisos con Docker/Podman: ejecuta `sudo usermod -aG docker $USER` y vuelve a iniciar sesión.
- Puerto ocupado: cambia el mapeo `-p` o finaliza el proceso que esté usando ese puerto.
- `go: command not found`: instala Go y verifica la variable de entorno `$PATH`.
- `containerd` no inicia: revisa el estado del servicio con `sudo systemctl status containerd`.

## 11. Verificación final

Una instalación correcta debe permitir:

- que cada servicio responda en la ruta `/health`
- que los puertos expuestos queden accesibles desde localhost
- que cada entorno pueda compilar y ejecutar la aplicación sin errores

Si la respuesta devuelve JSON con `status: "UP"` o similar, la instalación ha sido exitosa.

## 12. Créditos

Proyecto académico de Sistemas Operativos 1 - Carnet: 202300668
