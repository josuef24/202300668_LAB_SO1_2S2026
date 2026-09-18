# Guía de instalación y Manual Técnico - Proyecto 1 - Sistemas Operativos 1

Este repositorio contiene la infraestructura y código fuente de un sistema distribuido compuesto por tres APIs desarrolladas en Go. El proyecto implementa tres motores de contenedores distintos (Containerd, Podman y Docker) distribuidos en tres Máquinas Virtuales bajo el hipervisor KVM, centralizando las imágenes mediante un registro local Zot.

## 1. Arquitectura del Sistema

El proyecto se distribuye en la siguiente red (KVM/QEMU):
- **VM 1 (Containerd):** Aloja la `API 1` (Puerto 8081) y `API 2` (Puerto 8082).
- **VM 2 (Podman):** Aloja la `API 3` (Puerto 8080).
- **VM 3 (Docker):** Aloja el `Zot Registry` (Puerto 5000) utilizado para la distribución de imágenes.

## 2. Requisitos previos

Para desplegar este entorno, el sistema host debe contar con:
- Distribución de Linux (ej. Ubuntu 22.04/24.04).
- KVM/QEMU y Virt-Manager instalados.
- 3 Máquinas Virtuales configuradas en la misma red NAT (`192.168.122.x`).

## 3. Estructura del repositorio

```text
.
├── VM1_Containerd/
│   ├── API1/          # Código fuente (main.go) y Dockerfile
│   └── API2/          # Código fuente (main.go) y Dockerfile
└── VM2_Podman/
    └── API3/          # Código fuente (main.go) y Dockerfile