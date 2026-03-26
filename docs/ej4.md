# Ejercicio 4

## Objetivo

Implementar apagado graceful ante SIGTERM para cliente y servidor, asegurando liberacion ordenada de recursos y finalizacion controlada de procesos.

## Solucion (resumen)

1. El sistema se ejecuta con cliente y servidor en Docker Compose.
2. Se envia SIGTERM a un servicio puntual o a toda la composicion.
3. Cada proceso entra en flujo de cierre, detiene su loop principal y cierra sockets abiertos.
4. El proceso termina con estado exitoso cuando el cierre ocurre sin errores de recursos.
5. Los logs reflejan la secuencia de shutdown para validar que no haya abortos bruscos.

## Aspectos clave de la implementacion

1. Se registran handlers de señal para capturar SIGTERM y disparar cierre controlado.
2. El shutdown es idempotente para evitar dobles cierres sobre la misma conexion.
3. El criterio de validacion principal es que el contenedor objetivo finalice en success al recibir la señal.
