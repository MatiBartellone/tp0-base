# Ejercicio 2

## Objetivo

Permitir configurar cliente y servidor desde archivos montados por volumen, validando que los parametros de ejecucion se apliquen correctamente sin modificar la imagen.

## Solucion (resumen)

1. Se generan servicios con Docker Compose y se montan archivos de configuracion para cliente y servidor.
2. Antes de levantar, se escriben distintos valores de config (cliente y servidor) en sus archivos correspondientes.
3. Al iniciar contenedores, cada proceso carga configuracion desde volumen.
4. Los servicios loguean `action: config | result: success | ...` con los valores efectivos cargados.
5. Las pruebas comparan los campos de logs contra los valores esperados para confirmar que la inyeccion de config funciona.

## Aspectos clave de la implementacion

1. Configuracion desacoplada de la imagen: los cambios de parametros no requieren rebuild.
2. Validacion por observabilidad: los logs estructurados son la fuente de verdad de la config aplicada.
3. El servidor debe mantenerse activo aunque clientes finalicen su loop correctamente.
