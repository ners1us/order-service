package enum

type ProductType string

const (
	ProductElectronics ProductType = "электроника"
	ProductClothes     ProductType = "одежда"
	ProductShoes       ProductType = "обувь"
)

func (pr ProductType) String() string {
	return string(pr)
}
