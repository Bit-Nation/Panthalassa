package panthalassa

import (
	"errors"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	mesh "github.com/Bit-Nation/panthalassa/mesh"
	log "github.com/ipfs/go-log"
	"github.com/segmentio/objconv/json"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

var panthalassaInstance *Panthalassa
var logger = log.Logger("panthalassa")

type UpStream interface {
	Send(data string)
}

type StartConfig struct {
	EncryptedKeyManager string   `valid:"required"`
	RendezvousKey       string   `valid:"required"`
	Client              UpStream `valid:"required"`
	SignedProfile       string   `valid:"required"`
}

// create a new panthalassa instance
func start(km *keyManager.KeyManager, config StartConfig) error {

	//Exit if instance was already created and not stopped
	if panthalassaInstance != nil {
		return errors.New("call stop first in order to create a new panthalassa instance")
	}

	//Mesh network
	pk, err := km.MeshPrivateKey()
	if err != nil {
		return err
	}

	m, errReporter, err := mesh.New(pk, config.RendezvousKey)
	if err != nil {
		return err
	}
	//Report error's from mesh network to current logger
	go func() {
		for {
			select {
			case err := <-errReporter:
				logger.Error(err)
			}
		}
	}()

	//Create panthalassa instance
	panthalassaInstance = &Panthalassa{
		km:        km,
		upStream:  config.Client,
		deviceApi: deviceApi.New(config.Client),
		mesh:      m,
	}

	// register all housekeepers
	SearchContacts(panthalassaInstance)

	return nil

}

// start panthalassa
func Start(config string, password string) error {

	// unmarshal config
	var c StartConfig
	if err := json.Unmarshal([]byte(config), &c); err != nil {
		return err
	}

	// validate config
	_, err := valid.ValidateStruct(config)
	if err != nil {
		return err
	}

	// open key manager with password
	km, err := keyManager.OpenWithPassword(c.EncryptedKeyManager, password)
	if err != nil {
		return err
	}

	return start(km, c)
}

// create a new panthalassa instance with the mnemonic
func StartFromMnemonic(config string, mnemonic string) error {

	// unmarshal config
	var c StartConfig
	if err := json.Unmarshal([]byte(config), &c); err != nil {
		return err
	}

	// validate config
	_, err := valid.ValidateStruct(config)
	if err != nil {
		return err
	}

	// create key manager
	km, err := keyManager.OpenWithMnemonic(c.EncryptedKeyManager, mnemonic)
	if err != nil {
		return err
	}

	// create panthalassa instance
	return start(km, c)

}

//Eth Private key
func EthPrivateKey() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.GetEthereumPrivateKey()

}

func EthAddress() (string, error) {
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.GetEthereumAddress()
}

func SendResponse(id uint32, data string) error {

	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.deviceApi.Receive(id, data)
}

//Export the current account store with given password
func ExportAccountStore(pw, pwConfirm string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.Export(pw, pwConfirm)

}

func IdentityPublicKey() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.IdentityPublicKey()
}

//Stop panthalassa
func Stop() error {

	//Exit if not started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa in order to stop it")
	}

	//Stop panthalassa
	err := panthalassaInstance.Stop()
	if err != nil {
		return err
	}

	//Reset singleton
	panthalassaInstance = nil

	return nil
}
