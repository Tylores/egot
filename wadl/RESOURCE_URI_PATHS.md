# SEP 2 Resource URI Path Structure

Generated from `sep_wadl.xml` - highlights the complete resource path hierarchy for the IEEE 2030.5 Smart Energy Profile (SEP 2.2) API.

**Base URI:** `http://localhost/sep/`

## Resource Path Hierarchy Tree

```
/
├── /bill [CustomerAccountList]
│   └── /{id1} [CustomerAccount]
│       ├── /ca [CustomerAgreementList]
│       │   └── /{id2} [CustomerAgreement]
│       │       ├── /actbp [ActiveBillingPeriodList]
│       │       ├── /bp [BillingPeriodList]
│       │       │   └── /{id3} [BillingPeriod]
│       │       ├── /pro [ProjectionReadingList]
│       │       │   └── /{id3} [ProjectionReading]
│       │       ├── /tar [TargetReadingList]
│       │       │   └── /{id3} [TargetReading]
│       │       └── /ver [HistoricalReadingList]
│       │           └── /{id3} [HistoricalReading]
│       └── /ss [ServiceSupplier]
├── /brs [BillingReadingSetList]
│   └── /{id1} [BillingReadingSet]
│       └── /br [BillingReadingList]
│           └── /{id2} [BillingReading]
├── /dcap [DeviceCapability]
├── /derp [DERProgramList]
│   └── /{id1} [DERProgram]
│       ├── /actderc [ActiveDERControlList]
│       ├── /dc [DERCurveList]
│       │   └── /{id2} [DERCurve]
│       ├── /dderc [DefaultDERControl]
│       └── /derc [DERControlList]
│           └── /{id2} [DERControl]
├── /dr [DemandResponseProgramList]
│   └── /{id1} [DemandResponseProgram]
│       ├── /actedc [ActiveEndDeviceControlList]
│       └── /edc [EndDeviceControlList]
│           └── /{id2} [EndDeviceControl]
├── /edev [EndDeviceList]
│   └── /{id1} [EndDevice]
│       ├── /adev [AggregatedDeviceList]
│       │   └── /{id2} [AggregatedDevice]
│       ├── /aggp [AggregationPriority]
│       ├── /cfg [Configuration]
│       │   └── /prcfg [PriceResponseCfgList]
│       │       └── /{id2} [PriceResponseCfg]
│       ├── /der [DERList]
│       │   └── /{id2} [DER]
│       │       ├── /cdc [CurrentDERControls]
│       │       ├── /cdp [CurrentDERProgram]
│       │       ├── /dera [DERAvailability]
│       │       ├── /dercap [DERCapability]
│       │       ├── /dercom [DERComponentList]
│       │       │   └── /{id3} [DERComponent]
│       │       ├── /derg [DERSettings]
│       │       ├── /derp [AssociatedDERProgramList]
│       │       ├── /ders [DERStatus]
│       │       └── /upt [AssociatedUsagePoint]
│       ├── /di [DeviceInformation]
│       │   └── /loc [SupportedLocaleList]
│       │       └── /{id2} [SupportedLocale]
│       ├── /dstat [DeviceStatus]
│       ├── /frp [FlowReservationResponseList]
│       │   └── /{id2} [FlowReservationResponse]
│       ├── /frq [FlowReservationRequestList]
│       │   └── /{id2} [FlowReservationRequest]
│       ├── /fs [FileStatus]
│       ├── /fsa [FunctionSetAssignmentsList]
│       │   └── /{id2} [FunctionSetAssignments]
│       ├── /lel [LogEventList]
│       │   └── /{id2} [LogEvent]
│       ├── /lsl [LoadShedAvailabilityList]
│       │   └── /{id2} [LoadShedAvailability]
│       ├── /ns [IPInterfaceList]
│       │   └── /{id2} [IPInterface]
│       │       ├── /addr [IPAddrList]
│       │       │   └── /{id3} [IPAddr]
│       │       │       └── /rpl [RPLInstanceList]
│       │       │           └── /{id4} [RPLInstance]
│       │       │               └── /srt [RPLSourceRoutesList]
│       │       │                   └── /{id5} [RPLSourceRoutes]
│       │       └── /ll [LLInterfaceList]
│       │           └── /{id3} [LLInterface]
│       │               └── /nbh [NeighborList]
│       │                   └── /{id4} [Neighbor]
│       ├── /prxy [ProxiedDeviceList]
│       │   └── /{id2} [ProxiedDevice]
│       ├── /ps [PowerStatus]
│       ├── /rg [Registration]
│       └── /sub [SubscriptionList]
│           └── /{id2} [Subscription]
├── /file [FileList]
│   └── /{id1} [File]
├── /msg [MessagingProgramList]
│   └── /{id1} [MessagingProgram]
│       ├── /acttxt [ActiveTextMessageList]
│       └── /txt [TextMessageList]
│           └── /{id2} [TextMessage]
├── /mup [MirrorUsagePointList]
│   └── /{id1} [MirrorUsagePoint]
├── /ntfy [NotificationList]
│   └── /{id1} [Notification]
├── /ppy [PrepaymentList]
│   └── /{id1} [Prepayment]
│       ├── /ab [AccountBalance]
│       ├── /actsi [ActiveSupplyInterruptionOverrideList]
│       ├── /cr [CreditRegisterList]
│       │   └── /{id2} [CreditRegister]
│       ├── /os [PrepayOperationStatus]
│       └── /si [SupplyInterruptionOverrideList]
│           └── /{id2} [SupplyInterruptionOverride]
├── /rsps [ResponseSetList]
│   └── /{id1} [ResponseSet]
│       └── /rsp [ResponseList]
│           └── /{id2} [DrResponse]
├── /sdev [SelfDevice]
├── /tm [Time]
├── /tp [TariffProfileList]
│   └── /{id1} [TariffProfile]
│       └── /rc [RateComponentList]
│           └── /{id2} [RateComponent]
│               ├── /acttti [ActiveTimeTariffIntervalList]
│               └── /tti [TimeTariffIntervalList]
│                   └── /{id3} [TimeTariffInterval]
│                       └── /cti [ConsumptionTariffIntervalList]
│                           └── /{id4} [ConsumptionTariffInterval]
└── /upt [UsagePointList]
    └── /{id1} [UsagePoint]
        └── /mr [MeterReadingList]
            └── /{id2} [MeterReading]
                ├── /rs [ReadingSetList]
                │   └── /{id3} [ReadingSet]
                │       └── /r [ReadingList]
                │           └── /{id4} [Reading]
                └── /rt [ReadingType]
```

## Key Root Resources (Depth 1)

| Path | Resource ID | Purpose |
|------|-------------|---------|
| `/dcap` | DeviceCapability | Device capability discovery |
| `/sdev` | SelfDevice | Server device information |
| `/edev` | EndDeviceList | End device management |
| `/ntfy` | NotificationList | Event notifications |
| `/rsps` | ResponseSetList | Response management |
| `/file` | FileList | File management |
| `/dr` | DemandResponseProgramList | Demand response programs |
| `/msg` | MessagingProgramList | Messaging programs |
| `/upt` | UsagePointList | Usage point metering |
| `/tp` | TariffProfileList | Tariff profiles |
| `/bill` | CustomerAccountList | Customer billing |
| `/brs` | BillingReadingSetList | Billing readings |
| `/mup` | MirrorUsagePointList | Mirror usage points |
| `/ppy` | PrepaymentList | Prepayment programs |
| `/derp` | DERProgramList | DER programs |
| `/tm` | Time | Server time |

## Deep Nesting Examples

### End Device Resource Chain
```
/edev/{id1}                          [EndDevice]
├── /edev/{id1}/ns/{id2}             [IPInterface]
│   ├── /edev/{id1}/ns/{id2}/addr/{id3}                    [IPAddr]
│   │   └── /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}      [RPLInstance]
│   │       └── /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}  [RPLSourceRoutes]
│   └── /edev/{id1}/ns/{id2}/ll/{id3}                      [LLInterface]
│       └── /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}        [Neighbor]
└── /edev/{id1}/di/loc/{id2}         [SupportedLocale]
```

### Usage Point Resource Chain
```
/upt/{id1}                           [UsagePoint]
└── /upt/{id1}/mr/{id2}              [MeterReading]
    ├── /upt/{id1}/mr/{id2}/rs/{id3}                       [ReadingSet]
    │   └── /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}           [Reading]
    └── /upt/{id1}/mr/{id2}/rt                             [ReadingType]
```

### Billing Resource Chain
```
/bill/{id1}                          [CustomerAccount]
└── /bill/{id1}/ca/{id2}             [CustomerAgreement]
    ├── /bill/{id1}/ca/{id2}/bp/{id3}                      [BillingPeriod]
    ├── /bill/{id1}/ca/{id2}/pro/{id3}                     [ProjectionReading]
    ├── /bill/{id1}/ca/{id2}/tar/{id3}                     [TargetReading]
    └── /bill/{id1}/ca/{id2}/ver/{id3}                     [HistoricalReading]
```

### Tariff Profile Resource Chain
```
/tp/{id1}                            [TariffProfile]
└── /tp/{id1}/rc/{id2}               [RateComponent]
    └── /tp/{id1}/rc/{id2}/tti/{id3}                       [TimeTariffInterval]
        └── /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}         [ConsumptionTariffInterval]
```

## Template Parameter Patterns

- `{id1}` - First-level resource ID (e.g., CustomerAccount ID, EndDevice ID)
- `{id2}` - Second-level resource ID (e.g., CustomerAgreement ID, MeterReading ID)
- `{id3}` - Third-level resource ID (e.g., BillingPeriod ID, ReadingSet ID)
- `{id4}` - Fourth-level resource ID (e.g., Reading ID, TimeTariffInterval ID)
- `{id5}` - Fifth-level resource ID (e.g., RPLSourceRoutes ID)

## Path Abbreviations

| Abbreviation | Meaning |
|--------------|---------|
| `dcap` | Device Capability |
| `sdev` | Self Device |
| `edev` | End Device |
| `ntfy` | Notification |
| `rsps` | Response Sets |
| `rsp` | Response |
| `dr` | Demand Response |
| `edc` | End Device Control |
| `msg` | Messaging |
| `txt` | Text Messages |
| `upt` | Usage Point |
| `mr` | Meter Reading |
| `rs` | Reading Set |
| `rt` | Reading Type |
| `tp` | Tariff Profile |
| `rc` | Rate Component |
| `tti` | Time Tariff Interval |
| `cti` | Consumption Tariff Interval |
| `bill` | Billing Account |
| `ca` | Customer Agreement |
| `bp` | Billing Period |
| `brs` | Billing Reading Set |
| `br` | Billing Reading |
| `mup` | Mirror Usage Point |
| `ppy` | Prepayment |
| `derp` | DER Program |
| `der` | DER |
| `fsa` | Function Set Assignments |
| `adev` | Aggregated Device |
| `prxy` | Proxied Device |
| `lsl` | Load Shed Availability |
| `ns` | Network Services (IP) |
| `addr` | IP Address |
| `rpl` | RPL Instance |
| `srt` | RPL Source Routes |
| `ll` | Link Layer |
| `nbh` | Neighbor |
| `lel` | Log Events |
| `di` | Device Information |
| `loc` | Supported Locales |
| `ps` | Power Status |
| `fs` | File Status |
| `tm` | Time |

## Notes

- Template parameters (`{id1}`, `{id2}`, etc.) represent resource identifiers that vary at runtime
- List resources (ending in "List", e.g., `/edev`) support collection operations (GET, POST)
- Individual resources (identified by `/{idN}`, e.g., `/edev/{id1}`) support GET, PUT, DELETE operations
- Nested resources show hierarchical relationships (e.g., a MeterReading belongs to a UsagePoint)
- Maximum nesting depth observed: 5 levels (`/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}`)
- Full URIs: Combine base URL + path segments (e.g., `http://localhost/sep/edev/123/mr/456`)
