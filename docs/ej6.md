# Ejercicio 6

## Objetivo

Enviar apuestas en lotes (batch) desde cada agencia al servidor, garantizando que no se pierdan apuestas y que la agrupacion respete el maximo configurado.

## Solucion (resumen)

1. El cliente lee el archivo CSV de la agencia y convierte cada registro en una apuesta valida.
2. Las apuestas se agrupan en batches de tamano `max_amount` (o menor en el ultimo lote).
3. Cada batch se serializa en un mensaje binario y se envia al servidor.
4. El servidor procesa cada lote y registra la cantidad recibida con `action: apuesta_recibida`.
5. El total acumulado de apuestas procesadas coincide con la cantidad original del archivo.

## Aspectos clave de la implementacion

1. Se desacopla la lectura del dataset, la construccion de batches y el envio por protocolo.
2. La logica de particion de lotes garantiza que solo el ultimo batch pueda tener menos elementos que `max_amount`.
3. Se mantiene una semantica simple request/ack para detectar errores de envio o procesamiento por lote.
