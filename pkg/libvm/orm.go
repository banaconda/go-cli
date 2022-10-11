package libvm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseImage struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Path   string
	Format string
	Size   int64
}

type Volume struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Path     string
	Size     int64
	Format   string
	OriginID uint
	Origin   BaseImage
}

type Key struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Username string
	Rsa      string
	Path     string
}

type Network struct {
	gorm.Model
	Name    string `gorm:"unique"`
	Vlan    int32
	Cidr    string
	Gateway string
	Dns     string
}

type Domain struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Cpu       int64
	Memory    int64
	Mac       string
	Ip        string
	KeyId     uint
	Key       Key
	VolumeID  uint
	Volume    Volume
	NetworkID uint
	Network   Network
}

type VmerDB struct {
	db *gorm.DB
}

func NewVmerDB(path string) (*VmerDB, error) {
	vmerDB := &VmerDB{}
	err := vmerDB.Open(path)
	if err != nil {
		return nil, err
	}

	return vmerDB, nil
}

// vmer open db by path
func (vmerDB *VmerDB) Open(path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}
	vmerDB.db = db

	db.AutoMigrate(&BaseImage{})
	db.AutoMigrate(&Key{})
	db.AutoMigrate(&Network{})
	db.AutoMigrate(&Volume{})
	db.AutoMigrate(&Domain{})

	return nil
}

// insert base image
func (vmerDB *VmerDB) InsertBaseImage(baseImage *BaseImage) error {
	return vmerDB.db.Create(baseImage).Error
}

// delete base image
func (vmerDB *VmerDB) DeleteBaseImage(baseImage *BaseImage) error {
	return vmerDB.db.Unscoped().Delete(baseImage).Error
}

// update base image
func (vmerDB *VmerDB) UpdateBaseImage(baseImage *BaseImage) error {
	return vmerDB.db.Save(baseImage).Error
}

// get base image by name
func (vmerDB *VmerDB) GetBaseImageByName(name string) (*BaseImage, error) {
	var baseImage BaseImage
	err := vmerDB.db.Where("name = ?", name).First(&baseImage).Error
	if err != nil {
		return nil, err
	}
	return &baseImage, nil
}

// get base image by id
func (vmerDB *VmerDB) GetBaseImageById(id uint) (*BaseImage, error) {
	var baseImage BaseImage
	err := vmerDB.db.First(&baseImage, id).Error
	if err != nil {
		return nil, err
	}
	return &baseImage, nil
}

// get all base images
func (vmerDB *VmerDB) GetAllBaseImages() ([]BaseImage, error) {
	var baseImages []BaseImage
	err := vmerDB.db.Find(&baseImages).Error
	if err != nil {
		return nil, err
	}
	return baseImages, nil
}

// insert key
func (vmerDB *VmerDB) InsertKey(key *Key) error {
	return vmerDB.db.Create(key).Error
}

// delete key
func (vmerDB *VmerDB) DeleteKey(key *Key) error {
	return vmerDB.db.Unscoped().Delete(key).Error
}

// update key
func (vmerDB *VmerDB) UpdateKey(key *Key) error {
	return vmerDB.db.Save(key).Error
}

// get key by name
func (vmerDB *VmerDB) GetKeyByName(name string) (*Key, error) {
	var key Key
	err := vmerDB.db.Where("name = ?", name).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// get key by id
func (vmerDB *VmerDB) GetKeyById(id uint) (*Key, error) {
	var key Key
	err := vmerDB.db.First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// get all keys
func (vmerDB *VmerDB) GetAllKeys() ([]Key, error) {
	var keys []Key
	err := vmerDB.db.Find(&keys).Error
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// insert network
func (vmerDB *VmerDB) InsertNetwork(network *Network) error {
	return vmerDB.db.Create(network).Error
}

// delete network
func (vmerDB *VmerDB) DeleteNetwork(network *Network) error {
	return vmerDB.db.Unscoped().Delete(network).Error
}

// update network
func (vmerDB *VmerDB) UpdateNetwork(network *Network) error {
	return vmerDB.db.Save(network).Error
}

// get network by name
func (vmerDB *VmerDB) GetNetworkByName(name string) (*Network, error) {
	var network Network
	err := vmerDB.db.Where("name = ?", name).First(&network).Error
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// get network by id
func (vmerDB *VmerDB) GetNetworkById(id uint) (*Network, error) {
	var network Network
	err := vmerDB.db.First(&network, id).Error
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// get all networks
func (vmerDB *VmerDB) GetAllNetworks() ([]Network, error) {
	var networks []Network
	err := vmerDB.db.Find(&networks).Error
	if err != nil {
		return nil, err
	}
	return networks, nil
}

// insert volume
func (vmerDB *VmerDB) InsertVolume(volume *Volume) error {
	return vmerDB.db.Create(volume).Error
}

// delete volume
func (vmerDB *VmerDB) DeleteVolume(volume *Volume) error {
	return vmerDB.db.Unscoped().Delete(volume).Error
}

// update volume
func (vmerDB *VmerDB) UpdateVolume(volume *Volume) error {
	return vmerDB.db.Save(volume).Error
}

// get volume by name
func (vmerDB *VmerDB) GetVolumeByName(name string) (*Volume, error) {
	var volume Volume
	err := vmerDB.db.Where("name = ?", name).Preload(clause.Associations).First(&volume).Error
	if err != nil {
		return nil, err
	}
	return &volume, nil
}

// get volume by id
func (vmerDB *VmerDB) GetVolumeById(id uint) (*Volume, error) {
	var volume Volume
	err := vmerDB.db.Preload(clause.Associations).First(&volume, id).Error
	if err != nil {
		return nil, err
	}
	return &volume, nil
}

// get all volumes
func (vmerDB *VmerDB) GetAllVolumes() ([]Volume, error) {
	var volumes []Volume
	err := vmerDB.db.Preload(clause.Associations).Find(&volumes).Error
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

// insert domain
func (vmerDB *VmerDB) InsertDomain(domain *Domain) error {
	return vmerDB.db.Create(domain).Error
}

// delete domain
func (vmerDB *VmerDB) DeleteDomain(domain *Domain) error {
	return vmerDB.db.Unscoped().Delete(domain).Error
}

// update domain
func (vmerDB *VmerDB) UpdateDomain(domain *Domain) error {
	return vmerDB.db.Save(domain).Error
}

// get domain by name
func (vmerDB *VmerDB) GetDomainByName(name string) (*Domain, error) {
	var domain Domain
	err := vmerDB.db.Where("name = ?", name).Preload(clause.Associations).First(&domain).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// get domain by id
func (vmerDB *VmerDB) GetDomainById(id uint) (*Domain, error) {
	var domain Domain
	err := vmerDB.db.Preload(clause.Associations).First(&domain, id).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// get all domains
func (vmerDB *VmerDB) GetAllDomains() ([]Domain, error) {
	var domains []Domain
	err := vmerDB.db.Preload(clause.Associations).Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}
