# --- SYSTEM CONFIG ---
/system identity set name="R4-Edge"

# --- LOOPBACK ---
/interface bridge add name=lo0
/ip address add address=10.255.255.4/32 interface=lo0

# --- POINT-TO-POINT LINKS ---
# R4 to R2 (eth1)
/ip address add address=10.0.24.2/30 interface=ether1
# R4 to R3 (eth2)
/ip address add address=10.0.34.2/30 interface=ether2

# --- OSPFv3 (RouterOS v7) ---
/routing ospf instance add name=default router-id=10.255.255.4
/routing ospf area add instance=default name=backbone area-id=0.0.0.0

# OSPF Templates for v7 (PTP network type for fast convergence)
/routing ospf interface-template add area=backbone interfaces=lo0 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether1 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether2 type=ptp
