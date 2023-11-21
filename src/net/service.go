package net

import (
	"log"
)

//
// --- BroadCaster start ---
//

// BroadCaster broadcasts the message across the network,
// each node has its own BroadCaster service
type BroadCaster struct {
	Net      *Network // Network
	NodeName string   // node that uses broadcaster
}

func MakeBroadCaster(n *Network) *BroadCaster {
	return &BroadCaster{
		Net:      n,
		NodeName: "",
	}
}

func (bc *BroadCaster) GatherQuorum() map[string]*Node {
	return bc.Net.GetNodes()
}

//
// --- BroadCaster end ---
//

//
// --- Method Invoker Start ---
//

type anyFunc func(interface{}) interface{}

type MethodInvoker struct {
	reflectionMap map[string]anyFunc // func name <-> func
	NodeName      string
}

func MakeMethodInvoker(methods []string) *MethodInvoker {
	mi := &MethodInvoker{
		reflectionMap: make(map[string]anyFunc),
		NodeName:      "",
	}
	for _, methodName := range methods {
		mi.reflectionMap[methodName+"Map"] = getFuncByName(methodName + "Map")
		mi.reflectionMap[methodName+"Reduce"] = getFuncByName(methodName + "Reduce")
	}
	return mi
}

func (mi *MethodInvoker) InvokeMethod(methodName string, args interface{}) interface{} {
	if method, ok := mi.reflectionMap[methodName]; ok {
		a := method(args)
		return a
	}
	log.Fatalf("Service.InvokeMethod(): Method '%s' not found", methodName)
	return ""
}

func getFuncByName(name string) anyFunc {
	switch name {
	case "CountNodesMap":
		return CountNodesMap
	case "CountNodesReduce":
		return CountNodesReduce
	default:
		log.Fatalf("Service:getFuncByName(): No func named %s", name)
		return nil
	}
}

//
// --- Method Invoker End ---
//

type Storage struct {
	// for each node different Storage instance
}
