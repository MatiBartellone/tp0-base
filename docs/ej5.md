# Ejercicio 5

## Objetivo

Enviar una apuesta desde el cliente al servidor y verificar que la informacion llega integra y queda persistida correctamente.

## Solucion (resumen)

1. El cliente construye una apuesta a partir de la configuracion de entrada.
2. La apuesta se serializa y se envia por TCP al servidor.
3. El servidor recibe, deserializa y valida el mensaje.
4. Si la operacion es correcta, el servidor persiste la apuesta y responde confirmacion.
5. En logs se puede verificar trazabilidad de punta a punta (`apuesta_enviada` y `apuesta_almacenada`).

## Aspectos clave de la implementacion

1. Se mantiene un protocolo binario simple para asegurar interoperabilidad cliente-servidor.
2. La validacion principal del ejercicio es la consistencia del DNI enviado por cliente y almacenado por servidor.
3. La observabilidad se apoya en logs de exito/error por accion para facilitar diagnostico.
