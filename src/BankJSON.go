package main

/* sample

	{
    "LiveBanks": [
        {
            "bank_folder": "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
            "bank_name": "Pop",
            "display_order": 0,
            "is_factory": true
        },
        {
            "bank_folder": "6683CD2A-39D9-EEC0-6510-369203A24FFE",
            "bank_name": "Bank 2",
            "display_order": 9,
            "is_factory": false
        }
    ],
    "Favorites": []
}
*/

// LiveBanksJSON : json struct
type LiveBanksJSON struct {
	BankFolder   string `json:"bank_folder"`
	BankName     string `json:"bank_name"`
	DisplayOrder int    `json:"display_order"`
	IsFactory    bool   `json:"is_factory"`
}

// BankListJSON : json struct
type BankListJSON struct {
	LiveBanks []LiveBanksJSON `json:"LiveBanks"`
	Favorites []interface{}   `json:"Favorites"`
}
