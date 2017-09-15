package main

import (
  "fmt"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "os"
)

type cpuState struct{
  R []uint32
  Mode string
  Thumb string
}

var modes = [16]string{
  "User", "FIQ", "IRQ", "Supervisor", "Invalid",
  "Invalid", "Invalid", "Abort", "Invalid", "Invalid", "Invalid",
  "Undefined", "Invalid", "Invalid", "System",
 }

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("the send rom handler was triggered! \n");
  buf, err := ioutil.ReadAll(r.Body)
  if err!=nil {log.Fatal("request",err)}
  ROM = buf
}

func cpuStateHandler(w http.ResponseWriter, r *http.Request) {

  Clock_Tick()

  modeString := modes[(Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }


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

  writer := json.NewEncoder(w)
  writer.Encode(CpuState)
}

func regIncreaseHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("the cpu RegIncrease handler was triggered! \n")
  Reset()

  modeString := modes[(Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }

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
  http.HandleFunc("/cpuState", cpuStateHandler)
  http.HandleFunc("/regIncrease", regIncreaseHandler)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
