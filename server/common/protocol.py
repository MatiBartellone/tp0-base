import struct

from .utils import Bet


class ServerProtocol:
    def __init__(self, peer):
        self._peer = peer

    def recv_bet(self):
        agency = self._recv_string()
        first_name = self._recv_string()
        last_name = self._recv_string()
        document = self._recv_string()
        birthdate = self._recv_string()
        number = self._recv_string()
        return Bet(agency, first_name, last_name, document, birthdate, number)

    def send_ack(self, ok):
        code = b"\x01" if ok else b"\x00"
        self._peer.sendall(code)

    def _recv_exact(self, size):
        data = bytearray()
        while len(data) < size:
            chunk = self._peer.recv(size - len(data))
            if not chunk:
                raise ConnectionError("connection closed while receiving data")
            data.extend(chunk)
        return bytes(data)

    def _recv_string(self):
        raw_size = self._recv_exact(2)
        size = struct.unpack(">H", raw_size)[0]
        raw_value = self._recv_exact(size)
        return raw_value.decode("utf-8")


def recv_bet(peer):
    return ServerProtocol(peer).recv_bet()


def send_ack(peer, ok):
    ServerProtocol(peer).send_ack(ok)
