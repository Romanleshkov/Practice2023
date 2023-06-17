package main

type AccessTicket struct {
	Username string	`json:"username"`
	Password string	`json:"password"`
}

type AccessTicketResponse struct {
	Data struct{
		CSRFPreventionToken string	`json:"CSRFPreventionToken"`
		Ticket string				`json:"ticket"`
	}								`json:"data"`
}

type GetNodes struct {
	Nodes []struct{
		Node string				`json:"node"`
		Status string			`json:"status"`
	}							`json:"data"`
}

type GetLxcs struct {
	Lxcs []struct{
		Status string			`json:"status"`
		VmId string				`json:"vmid"`
	}							`json:"data"`
}

type GetLxcStatus struct {
	Data struct{
		Status string		`json:"status"`
	}						`json:"data"`
}

type GetNodeStatus struct {
	Data struct{
		Uptime int		`json:"uptime"`
	}						`json:"data"`
}


