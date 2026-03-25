import socket
import logging
import signal

from .utils import store_bets
from .protocol import ServerProtocol


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self.__register_signals()

    def __register_signals(self):
        signal.signal(signal.SIGTERM, self.shutdown)
        signal.signal(signal.SIGINT, self.shutdown)

    def shutdown(self, _signum=None, _frame=None):
        if not self._running:
            return

        self._running = False
        logging.info('action: shutdown | result: in_progress | resource: server_socket')
        try:
            self._server_socket.close()
            logging.info('action: shutdown | result: success | resource: server_socket')
        except OSError as e:
            logging.error(f'action: shutdown | result: fail | resource: server_socket | error: {e}')

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._running:
            client_sock = self.__accept_new_connection()
            if client_sock is None:
                continue
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = ServerProtocol(client_sock)
            bet = protocol.recv_bet()

            store_bets([bet])
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
            protocol.send_ack(True)
        except (OSError, ValueError, ConnectionError) as e:
            if self._running:
                logging.error(f'action: apuesta_almacenada | result: fail | error: {e}')
            try:
                ServerProtocol(client_sock).send_ack(False)
            except OSError:
                pass
        finally:
            try:
                client_sock.close()
            except OSError:
                pass

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        if not self._running:
            return None

        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError:
            return None
