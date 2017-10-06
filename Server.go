package main

import (
  "fmt"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "os"
)

//to send state of cpu to client
type cpuState struct{
  R []uint32
  Mode string
  Thumb string
}

//cpu operating modes, data for mode found in lower 5 bits of CPSR
var modes = [16]string{
  "User", "FIQ", "IRQ", "Supervisor", "Invalid",
  "Invalid", "Invalid", "Abort", "Invalid", "Invalid", "Invalid",
  "Undefined", "Invalid", "Invalid", "System",
 }

//sends page
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("the send rom handler was triggered! \n");
  buf, err := ioutil.ReadAll(r.Body)
  if err!=nil {log.Fatal("request",err)}
  ROM = buf
}

func cpuStateHandler(w http.ResponseWriter, r *http.Request) {

  //run a clock cycle
  Clock_Tick()

  //get cpu operating mode and thumb vs arm mode
  modeString := modes[(Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }

  //create a data structure with state of emulated cpu to send to client
  var CpuState cpuState
  regs := make([]uint32, 19, 19)
  for i := 0; i < 16; i++ {
		regs[i] = *(Cpu.R[i])
	}
  regs[16]=Cpu.CPSR
  regs[17]=*(Cpu.SPSR)
  regs[18]=Cpu.Fetch_reg
  CpuState.Mode = modeString
  CpuState.Thumb = thumbString
  CpuState.R = regs

	//send data
  writer := json.NewEncoder(w)
  writer.Encode(CpuState)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {

	//Reset the Cpu's state
	Reset()

  //get cpu operating mode and thumb vs arm mode
  modeString := modes[(Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }

  //create a data structure with state of emulated cpu to send to client
  var CpuState cpuState
  regs := make([]uint32, 19, 19)
  for i := 0; i < 16; i++ {
		regs[i] = *(Cpu.R[i])
	}
  regs[16]=Cpu.CPSR
  regs[17]=*(Cpu.SPSR)
  regs[18]=Cpu.Fetch_reg
  CpuState.Mode = modeString
  CpuState.Thumb = thumbString
  CpuState.R = regs

  //send data
  writer := json.NewEncoder(w)
  writer.Encode(CpuState)
}


func main() {
  port, exist := os.LookupEnv("PORT")
	if !(exist) {
		port = "8080"
	}
  http.Handle("/", http.FileServer(http.Dir("./Client")))
  http.HandleFunc("/sendRom", handler)
  http.HandleFunc("/cpuState", cpuStateHandler) // for running an instruction
  http.HandleFunc("/reset", resetHandler)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
