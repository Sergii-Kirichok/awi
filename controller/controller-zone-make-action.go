package controller

func (c *Controller) MakeAction(zoneId string) error {
	//fmt.Printf("ASK FOR MAKE ACTION FOR ZONE: %s\n", zoneId)
	_, err := c.auth.MakeBookmark(zoneId)
	return err
}
