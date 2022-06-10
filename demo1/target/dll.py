from ctypes import *

lib = CDLL("main.so")


lib.noArgNoRes.argtypes = []
def noArgNoRes() -> None:
  lib.noArgNoRes()


lib.oneArgNoRes.argtypes = [c_int]
def oneArgNoRes(b: int) -> None:
  lib.oneArgNoRes(b)


lib.twoArgsNoRes.argtypes = [c_int, c_int]
def twoArgsNoRes(a: int, b: int) -> None:
  lib.twoArgsNoRes(a, b)

lib.noArgOneRes.restype = c_char_p
lib.noArgOneRes.argtypes = []
def noArgOneRes() -> str:
  return lib.noArgOneRes().decode("utf-8")

lib.oneArgOneRes.restype = c_char_p
lib.oneArgOneRes.argtypes = [c_int]
def oneArgOneRes(a: int) -> str:
  return lib.oneArgOneRes(a).decode("utf-8")

lib.twoArgsOneRes.restype = c_int
lib.twoArgsOneRes.argtypes = [c_int, c_int]
def twoArgsOneRes(a: int, b: int) -> int:
  return lib.twoArgsOneRes(a, b)

lib.GetName.restype = c_char_p
lib.GetName.argtypes = [c_char_p]
def GetName(pyStr: str) -> str:
  return lib.GetName(pyStr.encode("utf-8")).decode("utf-8")

lib.number_add.restype = c_int
lib.number_add.argtypes = [c_int, c_int]
def number_add(a: int, b: int) -> int:
  return lib.number_add(a, b)

lib.addFloat.restype = c_float
lib.addFloat.argtypes = [c_float, c_float]
def addFloat(a: float, b: float) -> float:
  return lib.addFloat(a, b)

lib.addDouble.restype = c_double
lib.addDouble.argtypes = [c_double, c_double]
def addDouble(a: float, b: float) -> float:
  return lib.addDouble(a, b)

