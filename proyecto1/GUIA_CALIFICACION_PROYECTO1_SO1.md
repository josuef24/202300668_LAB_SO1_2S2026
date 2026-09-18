# Guía de Calificación - Sistemas Operativos 1 (Proyecto 1)

Esta guía contiene los comandos de validación directa para agilizar la calificación de los entregables del Proyecto 1.

- Llevar preparadas las entregas de UEDI y Classroom.
- Ingresar al Meet del laboratorio durante la calificación.
- Firmar el Excel al estar de acuerdo con la nota.

## 1. Preparación del Entorno (20 pts)

### 1.1 Validación de Máquinas Virtuales y KVM (10 pts)

Para comprobar que las 3 VMs están en ejecución bajo KVM desde la terminal del host:

```bash
# Información del host
neofetch

# Listar las máquinas virtuales activas
virsh list --all
```

### 1.2 Validación de Runtimes (10 pts)

Ingresar por SSH o directamente a cada VM y comprobar los motores de contenedores correspondientes:

```bash
# VM1: Verificar Containerd
sudo ctr namespace ls
sudo systemctl status containerd
docker info  # no debería estar instalado

# VM2: Verificar Podman
podman info
docker info  # no debería estar instalado

# VM3: Verificar Docker
docker info
```

## 2. Desarrollo y Contenerización (30 pts)

### 2.1 Desarrollo de APIs en Go (10 pts)

Comprobar la existencia de los contenedores en ejecución y comprobar que use el formato solicitado:

- `/api1/#CARNET/call-api2`
- `/api2/#CARNET/call-api1`
- `/api3/#CARNET/call-api1`
- y demás endpoints según corresponda.

Comprobar que el código de las API es de Go:

```bash
# VM1 (Containerd): Listar contenedores
sudo ctr -n default containers ls

# VM2 (Podman):
podman ps

# VM3 (Docker):
docker ps
```

### 2.2 Endpoint /health (10 pts)

Lanzar peticiones a los endpoints de salud de cada API para asegurar que funciona correctamente en Go:

```bash
# Realizar curl a la API correspondiente
curl -s http://<IP_VM>:<PUERTO>/health | jq .
```

Ejemplo:

```bash
curl -s http://192.168.122.10:8080/health | jq .
```

### 2.3 Creación de Dockerfiles y Construcción de Imágenes (10 pts)

Comprobar Dockerfiles.

Imágenes respetan el formato:

```text
[API#-#CARNET]
```

Utilizar la interfaz de Zot.

## 3. Registros y Comunicación Cruzada (25 pts)

### 3.1 / 3.2 Zot Registry y Descarga de Imágenes (15 pts)

Comprobar que el registro Zot está activo en la VM3 y listar las imágenes:

```bash
# Utilizar la interfaz de Zot
```

Descargar imagen en alguna VM.

### 3.3 / 4.1 Comunicación entre APIs (10 pts)

Comprobar la lógica de comunicación cruzada de las APIs.

Usar Thunder Client o similares.

#### Endpoints de llamada

- API1
  - GET `/api1/#CARNET/call-api2`
  - GET `/api1/#CARNET/call-api3`

- API2
  - GET `/api2/#CARNET/call-api1`
  - GET `/api2/#CARNET/call-api3`

- API3
  - GET `/api3/#CARNET/call-api1`
  - GET `/api3/#CARNET/call-api2`

Se realizará la comprobación con todas las VM funcionales, y deteniendo una VM para verificar el correcto funcionamiento de "status".

## 4. Documentación y Entregables

### 4.2 Manual técnico y Guía de instalación

Debe presentarse la documentación técnica y la guía de instalación del proyecto.

### 4.3 Repositorio en GitHub con estructura y guía

Verificar que el repositorio tenga la estructura correcta y que incluya la guía correspondiente.

## 5. Evaluación Final

### 5.1 Respuesta a preguntas relacionadas al proyecto

Se evaluará la comprensión del proyecto y la capacidad de explicar la implementación.

### 5.2 Capacidad de modificar y justificar código en revisión

Se revisará la capacidad del estudiante para modificar y justificar sus decisiones de código.

## Entregas requeridas

- Entrega en UEDI
- Entrega en Classroom

## Importante

Únicamente se calificará lo que se presente en el horario del alumno.
