package models

type RequestGetServiceCost struct {
	Pickup       City   `json:"pickup"`                      // откуда
	Delivery     City   `json:"delivery"`                    // куда
	SelfPickup   bool   `json:"selfPickup" example:"false"`  // самовывоз отправителя
	SelfDelivery bool   `json:"selfDelivery" example:"true"` // самовывоз получателя
	Parcel       Parcel `json:"parcel"`
}

//	@Description	Посылка
type Parcel struct {
	Weight   int `json:"weight" example:"5"`
	Length   int `json:"length" example:"20"`
	Width    int `json:"width" example:"20"`
	Height   int `json:"height" example:"20"`
	Quantity int `json:"quantity" example:"1"`
}
