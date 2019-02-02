package vm

func initSpecClass(vm *VM) {
	vm.mainThread.execGoobyLib("spec.gb")
}
