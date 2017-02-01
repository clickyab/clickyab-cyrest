package ad

import (
	"modules/misc/base"
	_ "modules/telegram/ad/adControllers"   // controller
	_ "modules/telegram/ad/ads"             // models
	_ "modules/telegram/ad/chanControllers" // controller
	_ "modules/telegram/ad/planControllers" // controller
)

func init() {
	base.RegisterPermission("confirm_ad", "confirm_ad")
}
