import socket, pprint
from time import time

from pynuoadmin import nuodb_cli

def latency_point(host, port = 48004, timeout = 5):
    """
    Calculate a latency point using sockets. If error return large number
    """

    # New Socket and Time out
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)

    # Start a timer
    s_start = time()

    # Try to Connect
    try:
        s.connect((host, int(port)))
        s.shutdown(socket.SHUT_RD)
    except:
        raise ValueError("Connection Timeout")

    # Stop Timer
    s_runtime = (time() - s_start) * 1000
    return float(s_runtime)

def measure(host,port=48005,runs=9):
    total = 0
    try:
        for n in range(0,runs):
            total += latency_point(host)
        return total / float(runs)
    except Value:
        return 99999

class LatencyCommands(nuodb_cli.AdminCommands):

    @nuodb_cli.subcommand
    def find_closest_admin(self):
        """
        Find the closest (based upon latency) admin process.  This should be in same zone
        """

        closest = None
        closest_measurement = 99999
        
        for s in self.conn._get_all('peers'):
            addr = s['address']
            host,port = addr.split(":")
            measurement = measure(host,int(port))
            if measurement < closest_measurement:
                closest = s['id']
                closest_measurement = measurement
        print(closest)
        return closest
