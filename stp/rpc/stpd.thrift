namespace go stpd
typedef i32 int
typedef i16 uint16
struct Dot1dStpPortEntryConfig{
	1 : i32 	Dot1dStpPortKey
	2 : i32 	Dot1dStpPortPriority
	3 : i32 	Dot1dStpPortEnable
	4 : i32 	Dot1dStpPortPathCost
	5 : i32 	Dot1dStpPortPathCost32
	6 : i32 	Dot1dStpPortProtocolMigration
	7 : i32 	Dot1dStpPortAdminPointToPoint
	8 : i32 	Dot1dStpPortAdminEdgePort
	9 : i32 	Dot1dStpPortAdminPathCost
}
struct Dot1dStpPortEntryState{
	1 : i32 	Dot1dStpPortPathCost32
	2 : i32 	Dot1dStpPortPriority
	3 : i32 	Dot1dStpPortAdminPointToPoint
	4 : i32 	Dot1dStpPortProtocolMigration
	5 : i32 	Dot1dStpPortKey
	6 : i32 	Dot1dStpPortAdminEdgePort
	7 : i32 	Dot1dStpPortPathCost
	8 : i32 	Dot1dStpPortAdminPathCost
	9 : i32 	Dot1dStpPortEnable
	10 : i32 	Dot1dStpPortState
	11 : string 	Dot1dStpPortDesignatedRoot
	12 : i32 	Dot1dStpPortDesignatedCost
	13 : string 	Dot1dStpPortDesignatedBridge
	14 : string 	Dot1dStpPortDesignatedPort
	15 : i32 	Dot1dStpPortForwardTransitions
	16 : i32 	Dot1dStpPortOperEdgePort
	17 : i32 	Dot1dStpPortOperPointToPoint
}
struct Dot1dStpPortEntryStateGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Dot1dStpPortEntryState> Dot1dStpPortEntryStateList
}
struct Dot1dStpBridgeConfig{
	1 : string 	Dot1dBridgeAddressKey
	2 : i32 	Dot1dStpPriorityKey
	3 : i32 	Dot1dStpBridgeMaxAge
	4 : i32 	Dot1dStpBridgeHelloTime
	5 : i32 	Dot1dStpBridgeForwardDelay
	6 : i32 	Dot1dStpBridgeForceVersion
	7 : i32 	Dot1dStpBridgeTxHoldCount
}
struct Dot1dStpBridgeState{
	1 : i32 	Dot1dStpBridgeForceVersion
	2 : string 	Dot1dBridgeAddressKey
	3 : i32 	Dot1dStpBridgeHelloTime
	4 : i32 	Dot1dStpBridgeTxHoldCount
	5 : i32 	Dot1dStpBridgeForwardDelay
	6 : i32 	Dot1dStpBridgeMaxAge
	7 : i32 	Dot1dStpPriorityKey
	8 : i32 	Dot1dStpProtocolSpecification
	9 : i32 	Dot1dStpTimeSinceTopologyChange
	10 : i32 	Dot1dStpTopChanges
	11 : string 	Dot1dStpDesignatedRoot
	12 : i32 	Dot1dStpRootCost
	13 : i32 	Dot1dStpRootPort
	14 : i32 	Dot1dStpMaxAge
	15 : i32 	Dot1dStpHelloTime
	16 : i32 	Dot1dStpHoldTime
	17 : i32 	Dot1dStpForwardDelay
}
struct Dot1dStpBridgeStateGetInfo {
	1: int StartIdx
	2: int EndIdx
	3: int Count
	4: bool More
	5: list<Dot1dStpBridgeState> Dot1dStpBridgeStateList
}
service STPDServices {
	bool CreateDot1dStpPortEntryConfig(1: Dot1dStpPortEntryConfig config);
	bool UpdateDot1dStpPortEntryConfig(1: Dot1dStpPortEntryConfig origconfig, 2: Dot1dStpPortEntryConfig newconfig, 3: list<bool> attrset);
	bool DeleteDot1dStpPortEntryConfig(1: Dot1dStpPortEntryConfig config);

	Dot1dStpPortEntryStateGetInfo GetBulkDot1dStpPortEntryState(1: int fromIndex, 2: int count);
	bool CreateDot1dStpBridgeConfig(1: Dot1dStpBridgeConfig config);
	bool UpdateDot1dStpBridgeConfig(1: Dot1dStpBridgeConfig origconfig, 2: Dot1dStpBridgeConfig newconfig, 3: list<bool> attrset);
	bool DeleteDot1dStpBridgeConfig(1: Dot1dStpBridgeConfig config);

	Dot1dStpBridgeStateGetInfo GetBulkDot1dStpBridgeState(1: int fromIndex, 2: int count);
}