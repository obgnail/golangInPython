from ctypes import *

lib = CDLL("main.so")


lib.sshTerminal.argtypes = [c_char_p]
def sshTerminal(configPath: str) -> None:
  lib.sshTerminal(configPath.encode("utf-8"))

