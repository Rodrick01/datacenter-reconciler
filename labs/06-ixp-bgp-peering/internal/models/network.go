package models

// Node representa un equipo de red en la simulación (Router, Switch o Host).
// Usar punteros para esta estructura en la lógica previene copias de memoria en
// topologías masivas.
type Node struct {
	Name     string
	Kind     string
	Image    string
	MgmtIPv4 string
	Ports    []string
	Binds    []string
	Exec     []string
}

// Link representa una conexión L2 (veth pair) entre dos endpoints dentro del orquestador.
type Link struct {
	Endpoints []string
}

// Topology representa la estructura de datos que Containerlab necesita para renderizar el cluster.
// Encapsula metadata global y slices pre-calculados de nodos y enlaces.
type Topology struct {
	Name       string
	MgmtSubnet string
	Nodes      []Node
	Links      []Link
}
