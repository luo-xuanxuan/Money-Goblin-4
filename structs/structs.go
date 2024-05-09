package structs

type World_Status struct {
	Name              string                 `json:"name"`
	Free_Company_List []*Free_Company_Status `json:"free_company_list"`
}

type Free_Company_Status struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Tanks            int                   `json:"tanks"`
	Repairs          int                   `json:"repairs"`
	Submersible_List []*Submersible_Status `json:"submersible_list"`
}

type Submersible_Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Return_Time int64  `json:"return_time"`
}

type Loot_Model struct {
	Timestamp int64  `json:"time"`
	Fc_id     string `json:"fcid"`
	Sub_id    string `json:"sub_id"`
	Player    string `json:"player"`
	World     string `json:"world"`
	Sector_id int    `json:"sector_id"`
	Item_id   int    `json:"item_id"`
	Quantity  int    `json:"quantity"`
}

type Timer_Model struct {
	Return_time int64  `json:"return_time"`
	Fc_id       string `json:"fcid"`
	Name        string `json:"name"`
	Sub_id      int32  `json:"sub_id"`
}
