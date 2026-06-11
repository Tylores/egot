package routes

import "fmt"

const (
	// API Gateway
	Gateway = "localhost:8443"

	// services (generated from scaffold)
	BRS             = "localhost:8010"
	Bill            = "localhost:8011"
	DCAP            = "localhost:8012"
	DR              = "localhost:8014"
	EDevice         = "localhost:8015"
	File            = "localhost:8016"
	MUP             = "localhost:8017"
	Messaging       = "localhost:8018"
	Notify          = "localhost:8019"
	PPY             = "localhost:8020"
	SDevice         = "localhost:8021"
	TariffProfile   = "localhost:8022"
	TimeOfUse       = "localhost:8023"
	UPT             = "localhost:8024"
	DER             = "localhost:8026"
	FlowReservation = "localhost:8027"
	Operator        = "localhost:8028"
	Rsps            = "localhost:8041"
)

// serviceMap maps each registered URL path template to the host:port of the
// service that handles it. Keys use the {id1}/{id2}/... naming convention
// consistent with Go's HTTP pattern syntax and the WADL samplePath attributes.
// Use LinkFor to build a complete https URL from any registered path.
var serviceMap = map[string]string{
	// Root Discovery & Core
	"/dcap": DCAP,
	"/tm":   TimeOfUse,

	// EDevice
	"/edev":                                         EDevice,
	"/edev/{id1}":                                   EDevice,
	"/edev/{id1}/adev":                              EDevice,
	"/edev/{id1}/adev/{id2}":                        EDevice,
	"/edev/{id1}/aggp":                              EDevice,
	"/edev/{id1}/cfg":                               EDevice,
	"/edev/{id1}/cfg/prcfg":                         EDevice,
	"/edev/{id1}/cfg/prcfg/{id2}":                   EDevice,
	"/edev/{id1}/di":                                EDevice,
	"/edev/{id1}/di/loc":                            EDevice,
	"/edev/{id1}/di/loc/{id2}":                      EDevice,
	"/edev/{id1}/dstat":                             EDevice,
	"/edev/{id1}/fs":                                EDevice,
	"/edev/{id1}/fsa":                               EDevice,
	"/edev/{id1}/fsa/{id2}":                         EDevice,
	"/edev/{id1}/lel":                               EDevice,
	"/edev/{id1}/lel/{id2}":                         EDevice,
	"/edev/{id1}/lsl":                               EDevice,
	"/edev/{id1}/lsl/{id2}":                         EDevice,
	"/edev/{id1}/ns":                                EDevice,
	"/edev/{id1}/ns/{id2}":                          EDevice,
	"/edev/{id1}/ns/{id2}/addr":                     EDevice,
	"/edev/{id1}/ns/{id2}/addr/{id3}":               EDevice,
	"/edev/{id1}/ns/{id2}/addr/{id3}/rpl":           EDevice,
	"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}":     EDevice,
	"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt": EDevice,
	"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}": EDevice,
	"/edev/{id1}/ns/{id2}/ll":                             EDevice,
	"/edev/{id1}/ns/{id2}/ll/{id3}":                       EDevice,
	"/edev/{id1}/ns/{id2}/ll/{id3}/nbh":                   EDevice,
	"/edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}":             EDevice,
	"/edev/{id1}/prxy":                                    EDevice,
	"/edev/{id1}/prxy/{id2}":                              EDevice,
	"/edev/{id1}/ps":                                      EDevice,
	"/edev/{id1}/rg":                                      EDevice,
	"/edev/{id1}/sub":                                     EDevice,
	"/edev/{id1}/sub/{id2}":                               EDevice,

	// DER (Extracted from EDevice)
	"/der":                    DER,
	"/der/{id1}":              DER,
	"/der/{id1}/cdc":          DER,
	"/der/{id1}/cdp":          DER,
	"/der/{id1}/dera":         DER,
	"/der/{id1}/dercap":       DER,
	"/der/{id1}/dercom":       DER,
	"/der/{id1}/dercom/{id2}": DER,
	"/der/{id1}/derg":         DER,
	"/der/{id1}/derp":         DER,
	"/der/{id1}/ders":         DER,
	"/der/{id1}/upt":          DER,

	// FlowReservation (Extracted from EDevice)
	"/frq":       FlowReservation,
	"/frq/{id1}": FlowReservation,
	"/frp":       FlowReservation,
	"/frp/{id1}": FlowReservation,

	// BRS
	"/brs":                BRS,
	"/brs/{id1}":          BRS,
	"/brs/{id1}/br":       BRS,
	"/brs/{id1}/br/{id2}": BRS,

	// Bill
	"/bill":                          Bill,
	"/bill/{id1}":                    Bill,
	"/bill/{id1}/ca":                 Bill,
	"/bill/{id1}/ca/{id2}":           Bill,
	"/bill/{id1}/ca/{id2}/actbp":     Bill,
	"/bill/{id1}/ca/{id2}/bp":        Bill,
	"/bill/{id1}/ca/{id2}/bp/{id3}":  Bill,
	"/bill/{id1}/ca/{id2}/pro":       Bill,
	"/bill/{id1}/ca/{id2}/pro/{id3}": Bill,
	"/bill/{id1}/ca/{id2}/tar":       Bill,
	"/bill/{id1}/ca/{id2}/tar/{id3}": Bill,
	"/bill/{id1}/ca/{id2}/ver":       Bill,
	"/bill/{id1}/ca/{id2}/ver/{id3}": Bill,
	"/bill/{id1}/ss":                 Bill,

	// DERP
	"/derp":                  DER,
	"/derp/{id1}":            DER,
	"/derp/{id1}/actderc":    DER,
	"/derp/{id1}/dc":         DER,
	"/derp/{id1}/dc/{id2}":   DER,
	"/derp/{id1}/dderc":      DER,
	"/derp/{id1}/derc":       DER,
	"/derp/{id1}/derc/{id2}": DER,

	// DR
	"/dr":                 DR,
	"/dr/{id1}":           DR,
	"/dr/{id1}/actedc":    DR,
	"/dr/{id1}/edc":       DR,
	"/dr/{id1}/edc/{id2}": DR,

	// File
	"/file":       File,
	"/file/{id1}": File,

	// Messaging
	"/msg":                 Messaging,
	"/msg/{id1}":           Messaging,
	"/msg/{id1}/acttxt":    Messaging,
	"/msg/{id1}/txt":       Messaging,
	"/msg/{id1}/txt/{id2}": Messaging,

	// MUP
	"/mup":       MUP,
	"/mup/{id1}": MUP,

	// Notify
	"/ntfy":       Notify,
	"/ntfy/{id1}": Notify,

	// PPY
	"/ppy":                PPY,
	"/ppy/{id1}":          PPY,
	"/ppy/{id1}/ab":       PPY,
	"/ppy/{id1}/actsi":    PPY,
	"/ppy/{id1}/cr":       PPY,
	"/ppy/{id1}/cr/{id2}": PPY,
	"/ppy/{id1}/os":       PPY,
	"/ppy/{id1}/si":       PPY,
	"/ppy/{id1}/si/{id2}": PPY,

	// SDevice
	"/sdev": SDevice,

	// TariffProfile
	"/tp":                                    TariffProfile,
	"/tp/{id1}":                              TariffProfile,
	"/tp/{id1}/rc":                           TariffProfile,
	"/tp/{id1}/rc/{id2}":                     TariffProfile,
	"/tp/{id1}/rc/{id2}/acttti":              TariffProfile,
	"/tp/{id1}/rc/{id2}/tti":                 TariffProfile,
	"/tp/{id1}/rc/{id2}/tti/{id3}":           TariffProfile,
	"/tp/{id1}/rc/{id2}/tti/{id3}/cti":       TariffProfile,
	"/tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}": TariffProfile,

	// UPT
	"/upt":                                 UPT,
	"/upt/{id1}":                           UPT,
	"/upt/{id1}/mr":                        UPT,
	"/upt/{id1}/mr/{id2}":                  UPT,
	"/upt/{id1}/mr/{id2}/rs":               UPT,
	"/upt/{id1}/mr/{id2}/rs/{id3}":         UPT,
	"/upt/{id1}/mr/{id2}/rs/{id3}/r":       UPT,
	"/upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}": UPT,
	"/upt/{id1}/mr/{id2}/rt":               UPT,

	// Rsps
	"/rsps":                 Rsps,
	"/rsps/{id1}":           Rsps,
	"/rsps/{id1}/rsp":       Rsps,
	"/rsps/{id1}/rsp/{id2}": Rsps,
	"/operator":             Operator,
}

// LinkFor returns the full https URL template for a registered path.
func LinkFor(path string) string {
	host, ok := serviceMap[path]
	if !ok {
		panic(fmt.Sprintf("routes: no service mapping for path %q", path))
	}
	return "https://" + host + path
}
