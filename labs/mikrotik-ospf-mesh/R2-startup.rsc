# --- SYSTEM CONFIG ---
/system identity set name="R2-Leaf"

# --- LOOPBACK ---
/interface bridge add name=lo0
/ip address add address=10.255.255.2/32 interface=lo0

# --- POINT-TO-POINT LINKS ---
# R2 to R1 (eth1)
/ip address add address=10.0.12.2/30 interface=ether1
# R2 to R4 (eth2)
/ip address add address=10.0.24.1/30 interface=ether2
# R2 to R3 (eth3)
/ip address add address=10.0.23.1/30 interface=ether3

# --- OSPFv3 (RouterOS v7) ---
/routing ospf instance add name=default router-id=10.255.255.2
/routing ospf area add instance=default name=backbone area-id=0.0.0.0

# OSPF Templates for v7 (PTP network type for fast convergence)
/routing ospf interface-template add area=backbone interfaces=lo0 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether1 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether2 type=ptp
/routing ospf interface-template add area=backbone interfaces=ether3 type=ptp
