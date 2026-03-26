import logging


def log_server_shutdown_in_progress():
    logging.info('action: shutdown | result: in_progress | resource: server_socket')


def log_server_shutdown_success():
    logging.info('action: shutdown | result: success | resource: server_socket')


def log_server_shutdown_failure(err):
    logging.error(f'action: shutdown | result: fail | resource: server_socket | error: {err}')


def log_accept_connections_in_progress():
    logging.info('action: accept_connections | result: in_progress')


def log_accept_connections_success(ip):
    logging.info(f'action: accept_connections | result: success | ip: {ip}')


def log_bet_received_success(batch_size):
    logging.info(f'action: apuesta_recibida | result: success | cantidad: {batch_size}')


def log_bet_received_failure(batch_size, err):
    logging.error(f'action: apuesta_recibida | result: fail | cantidad: {batch_size} | error: {err}')


def log_draw_success():
    logging.info('action: sorteo | result: success')


def log_server_config_success(port, listen_backlog, logging_level, expected_agencies):
    logging.debug(
        f'action: config | result: success | port: {port} | '
        f'listen_backlog: {listen_backlog} | logging_level: {logging_level} | '
        f'expected_agencies: {expected_agencies}'
    )
