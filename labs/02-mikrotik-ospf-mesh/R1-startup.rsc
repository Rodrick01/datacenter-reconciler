# --- SYSTEM CONFIG ---
/system identity set name="R1-Spine"

# --- LOOPBACK ---
/interface bridge add name=lo0
/ip address add address=10.255.255.1/32 interface=lo0

# --- POINT-TO-POINT LINKS ---
# R1 to R2 (eth1)
/ip address add address=10.0.12.1/30 interface=ether1
# R1 to R3 (eth2)
/ip address add address=10.0.13.1/30 interface=ether2

# --- OSPFv3 (RouterOS v7) ---
/routing ospf instance add name=default router-id=10.255.255.1
/routing ospf area add instance=default name=backbone area-id=0.0.0.0

# OSPF Templates for v7 (PTP network type for fast convergence)
/routing ospf interface-template add area=backbone interfaces=lo0 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether1 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether2 type=ptp
