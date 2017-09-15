package main

import (
  "fmt"
  "log"
  "net/http"
  "io/ioutil"
  emulator "./GoBoyAdvance/Emulator"
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
  emulator.ROM = buf
}

func cpuStateHandler(w http.ResponseWriter, r *http.Request) {

  emulator.Clock_Tick()

  modeString := modes[(emulator.Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((emulator.Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }


  var CpuState cpuState
  regs := make([]uint32, 19, 19)
  for i := 0; i < 16; i++ {
		regs[i] = *(emulator.Cpu.R[i])
	}
  regs[16]=emulator.Cpu.CPSR
  regs[17]=*(emulator.Cpu.SPSR)
  regs[18]=emulator.Cpu.Fetch_reg
  CpuState.Mode = modeString
  CpuState.Thumb = thumbString
  CpuState.R = regs

  writer := json.NewEncoder(w)
  writer.Encode(CpuState)
}

func regIncreaseHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("the cpu RegIncrease handler was triggered! \n")
  emulator.Reset()

  modeString := modes[(emulator.Cpu.CPSR&0x0000000F)]

  var thumbString string
  if((emulator.Cpu.CPSR&0x00000020)!=0){
    thumbString = "Thumb"
  } else {
    thumbString = "ARM"
  }

  var CpuState cpuState
  regs := make([]uint32, 19, 19)
  for i := 0; i < 16; i++ {
		regs[i] = *(emulator.Cpu.R[i])
	}
  regs[16]=emulator.Cpu.CPSR
  regs[17]=*(emulator.Cpu.SPSR)
  regs[18]=emulator.Cpu.Fetch_reg
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
