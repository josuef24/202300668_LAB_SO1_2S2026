# Proyecto 2 - SO1

Este directorio contiene una base funcional para la parte central del proyecto:

- `daemon/main.go`: daemon en Go para análisis, gestión y corrección de contenedores.
- `scripts/cron_create_containers.sh`: cronjob para desplegar contenedores de prueba.
- `kernel/continfo_pr2_so1_202300668.c`: plantilla de módulo de kernel en C para exponer métricas en `/proc`.

## Ejecución del daemon

```bash
cd proyecto2
go run ./daemon
```

## Carga del módulo kernel

El módulo real debe compilarse con headers del kernel y luego cargarse con:

```bash
sudo insmod continfo_pr2_so1_202300668.ko
cat /proc/continfo_pr2_so1_202300668
```

## Nota

La parte del daemon ha sido implementada para compilar en un entorno de desarrollo local; la parte del módulo de kernel debe ejecutarse en un host Linux con cabeceras del kernel disponibles.
