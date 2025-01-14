from __future__ import print_function
from ctypes import *
import json

lib = cdll.LoadLibrary("./sim.so")

# define class GoString to map:
# C type struct { const char *p; GoInt n; }
class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

# describe and call Log()
lib.Run.argtypes = [GoString]
lib.Run.restype = c_char_p

# config file

cfg = b"""
options mode=average iteration=1000 duration=90 simhp=0;

char+=xiangling ele=pyro lvl=80 hp=10121.775 atk=209.549 def=622.549 em=96.000 cr=.05 cd=0.5 cons=6 talent=6,8,8;
weapon+=xiangling label="favonius lance" atk=564.784 er=0.306 refine=1;
#weapon+=xiangling label="kitain cross spear" atk=564.784 em=110 refine=1;
art+=xiangling label="gladiator's finale" count=2;
art+=xiangling label="noblesse oblige" count=2;
stats+=xiangling label=main hp=4780 atk=311 pyro%=0.466 atk%=0.466 cr=0.311;
stats+=xiangling label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=xingqiu ele=hydro lvl=80 hp=9514.469 atk=187.803 def=705.132 atk%=0.240 cr=.05 cd=0.5 cons=6 talent=6,8,8;
weapon+=xingqiu label="blackcliff longsword" atk=565 cd=0.368 refine=1;
art+=xingqiu label="gladiator's finale" count=2;
art+=xingqiu label="noblesse oblige" count=2;
stats+=xingqiu label=main hp=4780 atk=311 hydro%=0.466 er=0.518 cr=0.311;
stats+=xingqiu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
weapon+=bennett label="favonius sword" atk=509.606 er=0.459 refine=5;
art+=bennett label="noblesse oblige" count=4;
stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=raiden ele=electro lvl=90 hp=12907 atk=337 def=789 electro%=0.288 cr=.05 cd=0.5 cons=0 talent=8,8,8;
weapon+=raiden label="grasscutter's light" atk=608 er=.551 refine=1;
art+=raiden label="seal of insulation" count=4;
stats+=raiden label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
stats+=raiden label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

##ENEMY
target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1 cryo=.1;
active+=raiden;

actions+=skill target=raiden if=.status.raidenskill==0;
actions+=sequence_strict target=xingqiu exec=skill,burst lock=100;
actions+=skill target=xingqiu if=.energy.xingqiu<80 lock=100;
actions+=burst target=xingqiu;
actions+=burst target=bennett;
actions+=sequence_strict target=xiangling exec=skill,burst;
actions+=sequence_strict target=raiden exec=burst,attack,attack,attack,attack,attack,attack,attack,attack,attack;
actions+=skill target=xiangling active=xiangling;
actions+=skill target=bennett if=.energy.xiangling<70 swap=xiangling;
actions+=skill target=raiden;
actions+=attack target=xiangling active=xiangling;
actions+=attack target=raiden active=raiden;
actions+=attack target=xingqiu active=xingqiu;
actions+=attack target=bennett active=bennett;
"""


msg = GoString(cfg, len(cfg))

result = str(lib.Run(msg), 'utf-8')

parsed = json.loads(result)

print(json.dumps(parsed, indent=4, sort_keys=True))