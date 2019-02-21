package p2p

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/p2p/common"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	// addresses under which the address book will claim To need more addresses.
	needAddressThreshold = 1000
	syncInterval         = 2 * time.Minute
	// getAddrMax is the most addresses that we will send in response
	// To a getAddr (in practise the most addresses we will return From a
	// call To AddressCache()).
	getAddrMax = 2500
)

// AddressManager is used To manage neighbor's address
type AddressManager struct {
	filePath  string
	ourAddrs  sync.Map
	addresses sync.Map
	lock      sync.RWMutex
	changed   bool
	quitChan  chan interface{}
}

// NewAddressManager create an address manager instance
func NewAddressManager(filePath string) *AddressManager {
	addresses := loadAddress(filePath)
	addrManager := &AddressManager{
		filePath: filePath,
		quitChan: make(chan interface{}),
	}
	addrManager.AddAddresses(addresses)
	return addrManager
}

// AddOurAddress add our local address.
func (addrManager *AddressManager) AddOurAddress(addr *common.NetAddress) {
	if _, ok := addrManager.ourAddrs.LoadOrStore(addr.ToString(), addr); ok {
		log.Debug("Add our address %s to address manager", addr.ToString())
	}
}

// AddOurAddress add our local address.
func (addrManager *AddressManager) AddLocalAddress(port int32) error {
	localIps, err := getLocalAddresses()
	if err != nil {
		return fmt.Errorf("failed To add our local address To address manager as:%v", err)
	}
	for _, localIp := range localIps {
		netAddr, err := common.ParseNetAddress(localIp + ":" + strconv.Itoa(int(port)))
		if err != nil {
			continue
		}
		addrManager.AddOurAddress(netAddr)
	}
	return nil
}

// OurAddresses get local address.
func (addrManager *AddressManager) OurAddresses() []*common.NetAddress {
	addrs := make([]*common.NetAddress, 0)
	addrManager.ourAddrs.Range(
		func(key, value interface{}) bool {
			addr := value.(*common.NetAddress)
			addrs = append(addrs, addr)
			return true
		},
	)
	return addrs
}

// IsOurAddress check whether the address is our address
func (addrManager *AddressManager) IsOurAddress(addr *common.NetAddress) bool {
	_, ok := addrManager.ourAddrs.Load(addr.ToString())
	return ok
}

// AddAddresses add new addresses
func (addrManager *AddressManager) AddAddresses(addrs []*common.NetAddress) {
	log.Debug("add %d addresses To book", len(addrs))
	for _, addr := range addrs {
		addrManager.AddAddress(addr)
	}
}

// AddAddress add a new address
func (addrManager *AddressManager) AddAddress(addr *common.NetAddress) {
	log.Debug("add new address %s To book", addr.ToString())
	addrManager.lock.Lock()
	defer addrManager.lock.Unlock()
	if _, ok := addrManager.ourAddrs.Load(addr.ToString()); ok {
		return
	}

	if _, ok := addrManager.addresses.LoadOrStore(addr.ToString(), addr); !ok {
		addrManager.changed = true
	}
}

// RemoveAddress remove an address
func (addrManager *AddressManager) RemoveAddress(addr *common.NetAddress) {
	addrManager.addresses.Delete(addr.ToString())
	addrManager.changed = true
}

// GetAddress get a random address
func (addrManager *AddressManager) GetAddress() (*common.NetAddress, error) {
	addrs := addrManager.GetAllAddress()
	if len(addrs) > 0 {
		index := rand.Intn(len(addrs))
		return addrs[index], nil
	}
	return nil, errors.New("no address in address book")
}

// GetAddresses get a random address list To send To peer
func (addrManager *AddressManager) GetAddresses() []*common.NetAddress {
	addrs := addrManager.GetAllAddress()
	if addrManager.GetAddressCount() <= getAddrMax {
		return addrs
	} else {
		for i := 0; i < getAddrMax; i++ {
			j := rand.Intn(getAddrMax-i) + i
			addrs[i], addrs[j] = addrs[j], addrs[i]
		}
		return addrs[:getAddrMax]
	}
}

// GetAddressCount get address count
func (addrManager *AddressManager) GetAddressCount() int {
	count := 0
	addrManager.addresses.Range(
		func(key, value interface{}) bool {
			count++
			return true
		},
	)
	return count
}

// GetAllAddress get all address
func (addrManager *AddressManager) GetAllAddress() []*common.NetAddress {
	addresses := make([]*common.NetAddress, 0)
	addrManager.addresses.Range(
		func(key, value interface{}) bool {
			addr := value.(*common.NetAddress)
			addresses = append(addresses, addr)
			return true
		},
	)
	return addresses
}

// NeedMoreAddrs check whether need more address.
func (addrManager *AddressManager) NeedMoreAddrs() bool {
	return addrManager.GetAddressCount() < needAddressThreshold
}

// Save save addresses To file
func (addrManager *AddressManager) Save() {
	addrManager.lock.Lock()
	if !addrManager.changed {
		addrManager.lock.Unlock()
		return
	}
	addrManager.lock.Unlock()

	addrStrs := make([]string, 0)
	addrManager.addresses.Range(
		func(key, value interface{}) bool {
			addrStr := key.(string)
			addrStrs = append(addrStrs, addrStr)
			return true
		},
	)

	buf, err := json.Marshal(addrStrs)
	fmt.Println(string(buf))
	if err != nil {
		log.Warn("failed To marshal recent addresses, as: %v", err)
	}

	err = ioutil.WriteFile(addrManager.filePath, buf, os.ModePerm)
	if err != nil {
		log.Warn("failed To write recent addresses To file, as: %v", err)
	}

	addrManager.changed = false
}

// Start start address manager
func (addrManager *AddressManager) Start() {
	go addrManager.saveHandler()
}

// Stop stop address manager
func (addrManager *AddressManager) Stop() {
	close(addrManager.quitChan)
}

// saveHandler save addresses To file periodically
func (addrManager *AddressManager) saveHandler() {
	saveFileTicker := time.NewTicker(syncInterval)
	for {
		select {
		case <-saveFileTicker.C:
			addrManager.Save()
		case <-addrManager.quitChan:
			return
		}
	}
}

// loadAddress load addresses From file.
func loadAddress(filePath string) []*common.NetAddress {
	addrStrs := make([]string, 0)
	addresses := make([]*common.NetAddress, 0)
	if _, err := os.Stat(filePath); err != nil {
		return addresses
	}
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("failed To read address book file, as: %v", err)
		return addresses
	}

	err = json.Unmarshal(buf, &addrStrs)
	if err != nil {
		log.Error("failed To parse address book file, as %v", err)
		return addresses
	}

	for _, addrStr := range addrStrs {
		addr, err := common.ParseNetAddress(addrStr)
		if err != nil {
			log.Warn("encounter an invalid address %s", addrStr)
			continue
		}
		addresses = append(addresses, addr)
	}
	log.Debug("load %d addresses From file %s", len(addresses), filePath)
	return addresses
}

// get all address of our server
func getLocalAddresses() ([]string, error) {
	ips := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error("failed To get system's interfaces")
		return nil, errors.New("failed To get system's interfaces")
	}
	for _, i := range ifaces {
		if skipInterface(i) {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			log.Warn("failed To get interface's address")
			continue
		}
		// handle err
		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() == nil || ip.IsLoopback() {
				log.Warn("skip invalid address %s", ip)
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}
	return false
}
