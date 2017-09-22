package accessors

import "log"

func GetEnv(hostname string) (Device, error) {

	log.Printf("[dbo] getting environment variables of hostname: %s", hostname)

	return Device{}, nil
}

func GetUi(hostname string) (string, error) {

	log.Printf("[dbo] getting ui configuration of %s", hostname)

	return "", nil
}

func AddNewRoom(room Room) (Room, error) {

	log.Printf("[dbo] adding new %s room: %s", room.Desig.Name, room.Name)

	return room, nil
}

//func ValidateDesignation(desig Designation) error {
//
//	log.Printf("[accessors] validating designation... %s", desig)
//
//	return nil
//}
//
//func ValidateRoom(room Room) error {
//
//	log.Printf("[accessors] validating room...")
//
//	if (desig.Name == nil) || (desig.ID == nil) {
//		return errors.New("invalid designation")
//	}
//
//}
